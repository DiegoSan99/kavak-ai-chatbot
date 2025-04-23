package load

type EmbeddingConfig struct {
	ChunkSize    int // Size of each text chunk for embedding
	ChunkOverlap int // Number of overlapping characters between chunks
}

type LLMContainer struct {
	Embedder        EmbeddingClient
	EmbeddingConfig EmbeddingConfig
}
