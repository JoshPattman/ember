package ember

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"errors"
)

func NewSQLCache(db *sql.DB) (EmbedderResponseCache, error) {
	c := &sqlCache{
		db: db,
	}
	err := c.setupDB()
	if err != nil {
		return nil, err
	}
	return c, nil
}

type sqlCache struct {
	db *sql.DB
}

func (cache *sqlCache) setupDB() error {
	// New table for embeddings
	query := `
	CREATE TABLE IF NOT EXISTS embed_cache (
		text TEXT PRIMARY KEY,
		embedding BLOB NOT NULL
	);`
	_, err := cache.db.Exec(query)
	if err != nil {
		return wrap(err, "failed to create embedding cache table")
	}
	return nil
}

// ===== Implementation of EmbedderResponseCache =====

func (cache *sqlCache) GetCachedEmbedding(input string) (bool, []float64, error) {
	row := cache.db.QueryRow(`SELECT embedding FROM embed_cache WHERE text=?;`, input)
	blob := []byte{}
	err := row.Scan(&blob)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, err
	}

	var embedding []float64
	err = gob.NewDecoder(bytes.NewBuffer(blob)).Decode(&embedding)
	if err != nil {
		return false, nil, err
	}

	return true, embedding, nil
}

func (cache *sqlCache) SetCachedEmbedding(input string, embedding []float64) error {
	blob := bytes.NewBuffer(nil)
	err := gob.NewEncoder(blob).Encode(embedding)
	if err != nil {
		return err
	}

	_, err = cache.db.Exec(`
		INSERT INTO embed_cache (text, embedding)
		VALUES (?, ?)
		ON CONFLICT(text) DO UPDATE SET embedding = excluded.embedding;
	`, input, blob.Bytes())
	return err
}
