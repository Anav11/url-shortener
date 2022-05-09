package handlers

type ShortenerResponseJSON struct {
	Result string `json:"result"`
}

type ShortenerRequestJSON struct {
	URL string `json:"url"`
}

type UserURLsJSON struct {
	ShortURL	string `json:"short_url"`
	OriginalURL	string `json:"original_url"`
}

type BatchURLs []struct {
	CorrelationID 	string `json:"correlation_id"`
	OriginalURL	string `json:"original_url"`
}

type BatchShortURL struct {
	CorrelationID 	string `json:"correlation_id"`
	ShortURL	string `json:"short_url"`
}
