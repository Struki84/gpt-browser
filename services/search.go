package services

import (
	"fmt"

	"gitlab.com/psheets/ddgquery"
	"gitlab.strukan.me/sandbox/gpt/browser/models"
)

type SearchService struct {
	Results []models.SearchResult
}

func NewSearchService() *SearchService {
	return &SearchService{}
}

func (search *SearchService) SearchDuckDuckGo(query string) error {
	results, _ := ddgquery.Query(query, 10)

	// log.Println("=======> Results:")
	// log.Println(results)
	// log.Println(url)
	// log.Println(query)

	if len(results) == 0 {
		return fmt.Errorf("no results found for query: %s", query)
	}

	for _, result := range results {
		search.Results = append(search.Results, models.SearchResult{
			Source: "duckduckgo.com",
			Title:  result.Title,
			Info:   result.Info,
			Ref:    result.Ref,
		})
	}

	return nil
}

func (search *SearchService) FormatResults() string {
	formattedResults := ""

	for _, result := range search.Results {
		formattedResults += fmt.Sprintf("Title: %s\nSource: %s\nDescription: %s\nURL: %s\n\n", result.Title, result.Source, result.Info, result.Ref)
	}

	return formattedResults
}
