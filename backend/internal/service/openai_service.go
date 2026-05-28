package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/deckforge/backend/internal/config"
	"github.com/sashabaranov/go-openai"
)

// GeneratedSlide is the shape GPT returns for each slide.
type GeneratedSlide struct {
	SlideType string   `json:"slide_type"`
	Title     string   `json:"title"`
	Subtitle  string   `json:"subtitle"`
	Bullets   []string `json:"bullets"`
}

// GeneratedDeck is the full structured response from OpenAI.
type GeneratedDeck struct {
	Title  string           `json:"title"`
	Slides []GeneratedSlide `json:"slides"`
}

// OpenAIService calls GPT to turn document text into slide content.
type OpenAIService struct {
	client *openai.Client
	model  string
}

func NewOpenAIService(cfg *config.Config) *OpenAIService {
	return &OpenAIService{
		client: openai.NewClient(cfg.OpenAIKey),
		model:  cfg.OpenAIModel,
	}
}

// GenerateSlides sends extracted document text to GPT and returns structured slides.
func (s *OpenAIService) GenerateSlides(ctx context.Context, documentText string) (*GeneratedDeck, error) {
	if s.client == nil {
		return nil, fmt.Errorf("OpenAI client not configured — set OPENAI_API_KEY")
	}

	const maxChars = 12000
	if len(documentText) > maxChars {
		documentText = documentText[:maxChars] + "\n...[truncated]"
	}

	systemPrompt := `You are an expert pitch deck consultant. Given source document text, create a professional startup pitch deck.

Return JSON with this exact structure:
{"title":"Deck Title","slides":[{"slide_type":"title","title":"...","subtitle":"...","bullets":[]}, ...]}

You MUST include exactly 7 slides in this order with slide_type:
1. title
2. problem
3. solution
4. market
5. features
6. roadmap
7. conclusion

Each slide needs title, subtitle (optional for non-title slides), and bullets (array of strings, can be empty for title).`

	userPrompt := fmt.Sprintf("Create a pitch deck from this document:\n\n%s", documentText)

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: userPrompt},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("OpenAI request failed: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("empty response from OpenAI")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	var deck GeneratedDeck
	if err := json.Unmarshal([]byte(content), &deck); err != nil {
		return nil, fmt.Errorf("parse OpenAI response: %w", err)
	}

	s.normalizeDeck(&deck)
	return &deck, nil
}

func (s *OpenAIService) normalizeDeck(deck *GeneratedDeck) {
	required := []string{"title", "problem", "solution", "market", "features", "roadmap", "conclusion"}
	if len(deck.Slides) >= 7 {
		return
	}
	existing := make(map[string]bool)
	for _, sl := range deck.Slides {
		existing[strings.ToLower(sl.SlideType)] = true
	}
	for _, t := range required {
		if existing[t] {
			continue
		}
		deck.Slides = append(deck.Slides, GeneratedSlide{
			SlideType: t,
			Title:     capitalize(t),
			Bullets:   []string{"Content to be refined"},
		})
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
