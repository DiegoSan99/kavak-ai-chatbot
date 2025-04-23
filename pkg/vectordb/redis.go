package vectordb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/load"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/redisvector"
)

// RedisVectorDB is a struct that provides vector database functionality using Redis.
type RedisVectorDB struct {
	embedder      load.EmbeddingClient
	redisHost     string
	redisPassword string
}

// NewRedisVectorDB creates a new RedisVectorDB instance.
func NewRedisVectorDB(embedder load.EmbeddingClient, redisHost, redisPassword string) *RedisVectorDB {
	return &RedisVectorDB{
		embedder:      embedder,
		redisHost:     redisHost,
		redisPassword: redisPassword,
	}
}

// getRedisHost returns the Redis connection URL.
func (r *RedisVectorDB) getRedisHost() (string, error) {
	var err error
	host := ""

	// Check if the Redis host is set in the configuration
	if r.redisHost == "" {
		err = errors.New("RedisHost is not set")
	} else {
		// Construct Redis connection string without authentication
		host = "redis://" + r.redisHost

		// If password is provided, include it in the connection string
		if r.redisPassword != "" {
			host = "redis://:" + r.redisPassword + "@" + r.redisHost
		}
	}

	return host, err
}

// AddDocuments adds documents to the Redis vector store.
//
// Parameters:
//   - prefix: A string prefix used to identify relevant vector entries.
//   - docs: The documents to add to the vector store.
//
// Returns:
//   - []string: The IDs of the added documents.
//   - error: An error if the operation fails.
func (r *RedisVectorDB) AddDocuments(prefix string, docs []schema.Document) ([]string, error) {
	if r.embedder == nil {
		return nil, errors.New("missing embedding model")
	} else {
		if !r.embedder.Initialized() {
			return nil, errors.New("embedding model not initialized")
		}
	}

	// Get the embedder from the client
	embedder, err := r.embedder.NewEmbedder()
	if err != nil {
		return nil, err
	}

	// Setup Redis vector store
	redisVector := redisvector.WithIndexName(prefix+"aillm_vector_idx", true)
	embedderVector := redisvector.WithEmbedder(embedder)

	redisHostURL, redisConnectionErr := r.getRedisHost()
	if redisConnectionErr != nil {
		return nil, redisConnectionErr
	}

	store, err := redisvector.New(context.TODO(), redisvector.WithConnectionURL(redisHostURL), redisVector, embedderVector)
	if err != nil {
		return nil, err
	}

	// Add documents to the store
	docIDs, err := store.AddDocuments(context.Background(), docs)
	if err != nil {
		return nil, fmt.Errorf("failed to add documents: %v", err)
	}

	return docIDs, nil
}

func (r *RedisVectorDB) CosineSimilarity(prefix, query string, rowCount int, scoreThreshold float32) ([]schema.Document, error) {
	var result []schema.Document
	if r.embedder == nil {
		return nil, errors.New("missing embedding model")
	} else {
		if !r.embedder.Initialized() {
			// Initialize the embedder if needed
			// This might need to be handled differently depending on your implementation
			return nil, errors.New("embedding model not initialized")
		}
	}

	// Get the embedder from the client
	embedder, err := r.embedder.NewEmbedder()
	if err != nil {
		return result, err
	}

	// Setup Redis vector store
	redisVector := redisvector.WithIndexName(prefix+"aillm_vector_idx", true)
	embedderVector := redisvector.WithEmbedder(embedder)

	redisHostURL, redisConnectionErr := r.getRedisHost()
	if redisConnectionErr != nil {
		return result, redisConnectionErr
	}
	store, err := redisvector.New(context.TODO(), redisvector.WithConnectionURL(redisHostURL), redisVector, embedderVector)
	if err != nil {
		return result, err
	}
	ctx := context.Background()
	optionsVector := []vectorstores.Option{
		vectorstores.WithScoreThreshold(scoreThreshold),
		vectorstores.WithEmbedder(embedder),
	}
	results, err := store.SimilaritySearch(ctx, query, rowCount, optionsVector...)
	if err != nil && !strings.Contains(err.Error(), "no such index") {
		return result, fmt.Errorf("search error: %v", err)
	}
	return results, nil
}

func (r *RedisVectorDB) FindKNN(prefix, searchQuery string, rowCount int, scoreThreshold float32) ([]schema.Document, error) {
	var result []schema.Document

	embedder, err := r.embedder.NewEmbedder()
	if err != nil {
		return result, err
	}

	redisVector := redisvector.WithIndexName(prefix+"aillm_vector_idx", true)
	embedderVector := redisvector.WithEmbedder(embedder)
	redisHostURL, redisConnectionErr := r.getRedisHost()
	if redisConnectionErr != nil {
		return result, redisConnectionErr
	}

	store, err := redisvector.New(context.TODO(), redisvector.WithConnectionURL(redisHostURL), redisVector, embedderVector)
	if err != nil {
		return result, err
	}

	optionsVector := []vectorstores.Option{
		vectorstores.WithScoreThreshold(scoreThreshold),
		vectorstores.WithEmbedder(embedder),
	}

	retriever := vectorstores.ToRetriever(store, rowCount, optionsVector...)

	resDocs, err := retriever.GetRelevantDocuments(context.Background(), searchQuery)
	if err != nil {
		return result, err
	}
	return resDocs, nil
}
