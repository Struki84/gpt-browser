package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type UserPrompt struct {
	Content string `json:"content"`
}

type Prompts struct {
	Persona string
	Search  string
	Analyse string
}

type GPTService struct {
	// Conversation is a list of all the responses from GPT API.
	// This way we store context of the conversation internally
	// since GPT API doesn't retain context between requests
	Conversation []string
	Client       *openai.Client
	Prompts      Prompts
}

func NewGPTService() *GPTService {
	return &GPTService{
		Conversation: []string{},
		Client:       openai.NewClient("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"),
		Prompts:      Prompts{},
	}
}

func (gpt *GPTService) LoadPrompts() *GPTService {
	// Open and read the JSON file
	jsonFile, err := os.Open("prompts.json")

	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Decode the JSON file into the data structure
	err = json.Unmarshal(byteValue, &gpt.Prompts)

	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}

	return gpt
}

func (gpt *GPTService) Query(prompt string) *GPTService {
	gpt.Conversation = []string{prompt}

	return gpt
}

func (gpt *GPTService) BuildSearchQuery(prompt *UserPrompt) *GPTService {
	p := fmt.Sprintf(gpt.Prompts.Search, prompt.Content)
	gpt.Conversation = append(gpt.Conversation, gpt.Prompts.Persona)
	gpt.Conversation = append(gpt.Conversation, p)

	return gpt
}

func (gpt *GPTService) AnalyseResults(data string) *GPTService {
	p := fmt.Sprintf(gpt.Prompts.Analyse, data)
	gpt.Conversation = append(gpt.Conversation, p)

	return gpt
}

func (gpt *GPTService) PromptGPT() (openai.ChatCompletionResponse, error) {
	GPTContext := strings.Join(gpt.Conversation, "\n")

	log.Println("GPT Context:", GPTContext)

	msg := []openai.ChatCompletionMessage{{
		Role:    openai.ChatMessageRoleUser,
		Content: GPTContext,
	}}

	request := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: msg,
	}

	return gpt.Client.CreateChatCompletion(
		context.Background(),
		request,
	)
}
