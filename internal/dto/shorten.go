package dto

type ShortenRequestDTO struct {
	URL string `json:"url"`
}

type ShortenResponseDTO struct {
	Result string `json:"result"`
}
