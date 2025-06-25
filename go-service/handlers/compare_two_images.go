package handlers

import (
    "context"
    "encoding/json"
    "io"
    "net/http"
    "os"

    "google.golang.org/genai"
)

type CompareResponse struct {
    Result string `json:"result"`
}

func HandlerCompareTwoImages(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()

    prompt := r.FormValue("prompt")
    if prompt == "" {
        http.Error(w, "Missing prompt", http.StatusBadRequest)
        return
    }

    // Get image1
    file1, _, err := r.FormFile("image1")
    if err != nil {
        http.Error(w, "Missing image1", http.StatusBadRequest)
        return
    }
    defer file1.Close()

    // Get image2
    file2, _, err := r.FormFile("image2")
    if err != nil {
        http.Error(w, "Missing image2", http.StatusBadRequest)
        return
    }
    defer file2.Close()

    // Read image2 as bytes (inline)
    image2Bytes, err := io.ReadAll(file2)
    if err != nil {
        http.Error(w, "Reading image2 failed", http.StatusInternalServerError)
        return
    }

    // Init Gemini client
    client, err := genai.NewClient(ctx, &genai.ClientConfig{
        APIKey: os.Getenv("GEMINI_API_KEY"),
    })
    if err != nil {
        http.Error(w, "Failed to create Gemini client", http.StatusInternalServerError)
        return
    }

    // Save image1 to temp file and upload via path
    tempFile, err := os.CreateTemp("", "image1-*.jpg")
    if err != nil {
        http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
        return
    }
    defer os.Remove(tempFile.Name())

    _, err = io.Copy(tempFile, file1)
    if err != nil {
        http.Error(w, "Failed to write image1 to temp file", http.StatusInternalServerError)
        return
    }

    uploadedFile, err := client.Files.UploadFromPath(ctx, tempFile.Name(), nil)
    if err != nil {
        http.Error(w, "Uploading image1 failed", http.StatusInternalServerError)
        return
    }

    // Prepare parts
    parts := []*genai.Part{
        genai.NewPartFromText(prompt),
        genai.NewPartFromBytes(image2Bytes, "image/jpeg"), // Inline image2
        genai.NewPartFromURI(uploadedFile.URI, uploadedFile.MIMEType), // URI image1
    }

    contents := []*genai.Content{
        genai.NewContentFromParts(parts, genai.RoleUser),
    }

    result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", contents, nil)
    if err != nil {
        http.Error(w, "Gemini comparison failed: "+err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(CompareResponse{Result: result.Text()})
}
