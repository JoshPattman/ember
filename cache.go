package ember

type EmbedderResponseCache interface {
	GetCachedEmbedding(string) (bool, []float64, error)
	SetCachedEmbedding(string, []float64) error
}
