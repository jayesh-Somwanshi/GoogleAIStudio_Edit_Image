package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"google.golang.org/genai"
)

func HandleImageUnderstandingFromUrl(w http.ResponseWriter, r *http.Request) {
	imageURL := r.FormValue("imageUrl")
	prompt := r.FormValue("prompt")

	if imageURL == "" || prompt == "" {
		http.Error(w, "imageUrl and prompt are required", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(imageURL)
	if err != nil {
		http.Error(w, "Failed to fetch image", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read image data", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GOOGLE_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		http.Error(w, "Gemini client creation failed", http.StatusInternalServerError)
		return
	}

	parts := []*genai.Part{
		genai.NewPartFromBytes(imageBytes, "image/jpeg"),
		genai.NewPartFromText(prompt),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", contents, nil)
	if err != nil {
		http.Error(w, "Gemini failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"summary": result.Text(),
	})
}

