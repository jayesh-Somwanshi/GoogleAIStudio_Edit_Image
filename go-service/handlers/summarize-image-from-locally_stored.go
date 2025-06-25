package handlers

import (
    "context"
    "encoding/json"
    "io"
    "net/http"
    "os"

    "google.golang.org/genai"
)

type CaptionResponse struct {
    Caption string `json:"caption"`
}

func HandleImageLocallySummary(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()

    prompt := r.FormValue("prompt")
    if prompt == "" {
        http.Error(w, "Missing prompt", http.StatusBadRequest)
        return
    }

    file, _, err := r.FormFile("image")
    if err != nil {
        http.Error(w, "Image upload failed", http.StatusBadRequest)
        return
    }
    defer file.Close()

    imageBytes, err := io.ReadAll(file)
    if err != nil {
        http.Error(w, "Reading image failed", http.StatusInternalServerError)
        return
    }

    client, err := genai.NewClient(ctx, &genai.ClientConfig{
        APIKey: os.Getenv("GEMINI_API_KEY"),
    })
    if err != nil {
        http.Error(w, "Failed to create Gemini client", http.StatusInternalServerError)
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
        http.Error(w, "Gemini caption generation failed", http.StatusInternalServerError)
        return
    }

    resp := CaptionResponse{Caption: result.Text()}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
