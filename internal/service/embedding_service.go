package service

import (
	"music-store/utils"
)

// EmbeddingService abstracts embedding generation.
type EmbeddingService interface {
	Generate(text string) []float32
	Dim() int
}

type embeddingService struct {
	client *utils.HugotClient
}

func NewEmbeddingService(client *utils.HugotClient) EmbeddingService {
	return &embeddingService{client: client}
}

func (e *embeddingService) Generate(text string) []float32 {
	return e.client.Generate(text)
}

func (e *embeddingService) Dim() int {
	if e.client == nil || e.clientCfgDim() <= 0 {
		return 512
	}
	return e.clientCfgDim()
}

func (e *embeddingService) clientCfgDim() int {
	// expose cfg.Dim indirectly
	// hugot client keeps cfg private, but we can infer by generating and reading length
	v := e.client.Generate("dim-probe")
	return len(v)
}
