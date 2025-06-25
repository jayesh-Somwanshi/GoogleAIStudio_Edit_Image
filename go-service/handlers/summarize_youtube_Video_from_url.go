package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"google.golang.org/genai"
)

func HandleVideoYoutubeURL(w http.ResponseWriter, r *http.Request) {
	videoURL := r.FormValue("videoUrl")
	prompt := r.FormValue("prompt")

	if videoURL == "" || prompt == "" {
		http.Error(w, "videoUrl and prompt are required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GOOGLE_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		http.Error(w, "Gemini client error", http.StatusInternalServerError)
		return
	}

	parts := []*genai.Part{
		genai.NewPartFromText(prompt),
		genai.NewPartFromURI(videoURL, "video/mp4"),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash", contents, nil)
	if err != nil {
		http.Error(w, "Gemini summarization failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"summary": result.Text(),
	})
}
