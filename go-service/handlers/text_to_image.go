package handlers

import (
	"context"
	"net/http"
	"os"

	"google.golang.org/genai"
)

func HandleTextToImage(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	if description == "" {
		http.Error(w, "Description is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		http.Error(w, "Failed to create Gemini client", http.StatusInternalServerError)
		return
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-preview-image-generation",
		[]*genai.Content{{
			Role: "user",
			Parts: []*genai.Part{
				genai.NewPartFromText(description),
			},
		}},
		&genai.GenerateContentConfig{ResponseModalities: []string{"TEXT", "IMAGE"}},
	)
	if err != nil {
		http.Error(w, "Gemini content generation failed", http.StatusInternalServerError)
		return
	}

	for _, part := range result.Candidates[0].Content.Parts {
		if part.InlineData != nil {
			w.Header().Set("Content-Type", "image/png")
			w.Write(part.InlineData.Data)
			return
		}
	}

	http.Error(w, "No image returned", http.StatusInternalServerError)
}
