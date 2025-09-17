package ember

// NewInMemoryCache creates an in-memory implementation of ModelResponseCache.
// It stores model responses in memory using a hash of the input messages as a key.
func NewInMemoryCache() EmbedderResponseCache {
	return &inMemoryCache{
		embs: make(map[string][]float64),
	}
}

type inMemoryCache struct {
	embs map[string][]float64
}

func (cache *inMemoryCache) GetCachedEmbedding(s string) (bool, []float64, error) {
	if emb, ok := cache.embs[s]; ok {
		return true, emb, nil
	}
	return false, nil, nil
}

func (cache *inMemoryCache) SetCachedEmbedding(s string, e []float64) error {
	cache.embs[s] = e
	return nil
}
