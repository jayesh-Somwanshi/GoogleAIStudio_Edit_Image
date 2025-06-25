package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"google.golang.org/genai"
)

var generatedImageDir = "GeneratedImage"

func ServeGeneratedImages() {
	fs := http.FileServer(http.Dir(generatedImageDir))
	http.Handle("/GeneratedImage/", http.StripPrefix("/GeneratedImage/", fs))
}

func clearGeneratedImages() {
	entries, err := os.ReadDir(generatedImageDir)
	if err != nil {
		log.Printf("Failed to read directory: %v", err)
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			os.Remove(filepath.Join(generatedImageDir, entry.Name()))
		}
	}
}

func HandleGenerateMulti(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	if description == "" {
		http.Error(w, "Description is required", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imgData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read image", http.StatusInternalServerError)
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

	if _, err := os.Stat(generatedImageDir); os.IsNotExist(err) {
		os.MkdirAll(generatedImageDir, 0755)
	} else {
		clearGeneratedImages()
	}

	var urls []string
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

	for i := 1; i <= 3; i++ {    // generate by default 3 images
		prompt := fmt.Sprintf("%s (variation %c)", description, 'A'+i-1)

		parts := []*genai.Part{
			genai.NewPartFromText(prompt),
			{
				InlineData: &genai.Blob{
					MIMEType: "image/png",
					Data:     imgData,
				},
			},
		}

		result, err := client.Models.GenerateContent(
			ctx,
			"gemini-2.0-flash-preview-image-generation",
			[]*genai.Content{{Role: "user", Parts: parts}},
			&genai.GenerateContentConfig{ResponseModalities: []string{"TEXT", "IMAGE"}},
		)
		if err != nil {
			log.Printf("Generation failed on attempt %d: %v", i, err)
			continue
		}

		if len(result.Candidates) == 0 {
			log.Printf("No candidates returned on attempt %d", i)
			continue
		}

		found := false
		for _, part := range result.Candidates[0].Content.Parts {
			if part.InlineData != nil {
				filename := fmt.Sprintf("img_%s_%d.png", timestamp, i)
				filepath := filepath.Join(generatedImageDir, filename)
				if err := os.WriteFile(filepath, part.InlineData.Data, 0644); err != nil {
					log.Printf("Failed to write file: %v", err)
					continue
				}
				urls = append(urls, fmt.Sprintf("/GeneratedImage/%s", filename))
				found = true
			}
		}

		if !found {
			log.Printf("No image returned in result on attempt %d", i)
		}
	}

	if len(urls) == 0 {
		http.Error(w, "No images generated", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{
		"images": urls,
	})
}