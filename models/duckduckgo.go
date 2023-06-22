package models

type DuckDuckGoResult struct {
	Results []struct {
		Title    string `json:"Title"`
		FirstURL string `json:"FirstURL"`
		Text     string `json:"Text"`
	} `json:"Results"`
	RelatedTopics []struct {
		Title    string `json:"Title"`
		FirstURL string `json:"FirstURL"`
		Text     string `json:"Text"`
	} `json:"RelatedTopics"`
}
