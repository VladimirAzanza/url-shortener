package dto

type ShortenRequestDTO struct {
	URL string `json:"url"`
}

type ShortenResponseDTO struct {
	Result string `json:"result"`
}

type BatchRequestDTO struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponseDTO struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
