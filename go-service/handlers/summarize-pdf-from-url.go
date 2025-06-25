package handlers

import (
    "context"
    "encoding/json"
    "io"
    "net/http"
    "os"

    "google.golang.org/genai"
)

func HandlerSummarizePdfFromUrl(w http.ResponseWriter, r *http.Request) {
    prompt := r.FormValue("prompt")
    if prompt == "" {
        prompt = "Summarize this document"
    }

    file, _, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "PDF file is required", http.StatusBadRequest)
        return
    }
    defer file.Close()

    pdfData, err := io.ReadAll(file)
    if err != nil {
        http.Error(w, "Failed to read PDF", http.StatusInternalServerError)
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

    parts := []*genai.Part{
        {
            InlineData: &genai.Blob{
                MIMEType: "application/pdf",
                Data:     pdfData,
            },
        },
        genai.NewPartFromText(prompt),
    }

    contents := []*genai.Content{
        genai.NewContentFromParts(parts, genai.RoleUser),
    }

    result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", contents, nil)
    if err != nil {
        http.Error(w, "Gemini generation failed", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "summary": result.Text(),
    })
}
