package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"google.golang.org/genai"
)

func HandlePDFLocallySummary(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	if description == "" {
		description = "Give me a summary of this PDF file."
	}

	file, _, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "PDF is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tmpFile, err := os.CreateTemp("", "*.pdf")
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, file)
	if err != nil {
		http.Error(w, "Failed to save temp PDF", http.StatusInternalServerError)
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

	uploadConfig := &genai.UploadFileConfig{MIMEType: "application/pdf"}
	uploadedFile, err := client.Files.UploadFromPath(ctx, tmpFile.Name(), uploadConfig)
	if err != nil {
		http.Error(w, "Upload failed", http.StatusInternalServerError)
		return
	}

	parts := []*genai.Part{
		genai.NewPartFromURI(uploadedFile.URI, uploadedFile.MIMEType),
		genai.NewPartFromText(description),
	}
	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", contents, nil)
	if err != nil {
		http.Error(w, "Failed to generate content", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, result.Text())
}
