package main

import (
	"strings"

	"gitlab.sintezis.co/sintezis/sdk/web/rest"
	"gitlab.strukan.me/sandbox/gpt/browser/services"

	"gitlab.com/psheets/ddgquery"
)

func main() {
	var userPrompt services.UserPrompt

	gptService := services.NewGPTService()
	gptService.LoadPrompts()

	searchService := services.NewSearchService()

	router := rest.NewRouter()
	ctrl := rest.NewController(router)

	ctrl.Post("/gpt/test", func(ctx *rest.Context) {
		ctx.SetContentType("application/json")

		err := ctx.JsonDecode(&userPrompt)

		if err != nil {
			ctx.JsonResponse(404, err)
			return
		}

		resp, err := gptService.Query(userPrompt.Content).PromptGPT()

		if err != nil {
			ctx.JsonResponse(500, err)
			return
		}

		ctx.JsonResponse(200, resp)
	})

	ctrl.Post("/gpt/search", func(ctx *rest.Context) {
		ctx.SetContentType("application/json")

		err := ctx.JsonDecode(&userPrompt)

		if err != nil {
			ctx.JsonResponse(400, err)
			return
		}

		// Create a prompt asking GPT for the best search query
		searchQueryResp, err := gptService.BuildSearchQuery(&userPrompt).PromptGPT()

		if err != nil {
			ctx.JsonResponse(500, err)
			return
		}

		searchQuery := strings.TrimSpace(searchQueryResp.Choices[0].Message.Content)

		// Fetch search results from DuckDuckGo
		err = searchService.SearchDuckDuckGo(searchQuery)

		if err != nil {
			ctx.JsonResponse(404, err.Error())
			return
		}

		formattedResults := searchService.FormatResults()

		// Prompt GPT to analyze the search results
		resp, err := gptService.AnalyseResults(formattedResults).PromptGPT()

		if err != nil {
			ctx.JsonResponse(500, err)
			return
		}

		ctx.JsonResponse(200, resp)
	})

	ctrl.Get("/search", func(ctx *rest.Context) {
		ctx.SetContentType("application/json")

		results, _ := ddgquery.Query("Who is Current World FIFA Champion?", 5)

		if len(results) == 0 {
			ctx.JsonResponse(404, "No search results")
			return
		}

		ctx.JsonResponse(200, results)

	})

	router.Mux.StrictSlash(true)
	router.Listen(":8082")
}
