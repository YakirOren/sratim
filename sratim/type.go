package sratim

type SearchResponse struct {
	Success bool           `json:"success"`
	Results []SearchResult `json:"results"`
}

type SearchResult struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	Errors      []string `json:"errors,omitempty"`
	ErrorsCount int      `json:"errors_count,omitempty"`
	Success     bool     `json:"success"`
	Watch       struct {
		URL string `json:"480"`
	} `json:"watch,omitempty"`
}
