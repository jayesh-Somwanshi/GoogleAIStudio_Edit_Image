package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"google.golang.org/genai"
)

func HandleCodeExecution(w http.ResponseWriter, r *http.Request) {
	prompt := r.FormValue("prompt")
	if prompt == "" {
		http.Error(w, "Prompt is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		http.Error(w, "Gemini client error", http.StatusInternalServerError)
		return
	}

	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{CodeExecution: &genai.ToolCodeExecution{}},
		},
	}

	chat, err := client.Chats.Create(ctx, "gemini-2.5-flash", config, nil)
	if err != nil {
		http.Error(w, "Chat create error", http.StatusInternalServerError)
		return
	}

	result, err := chat.SendMessage(ctx, genai.Part{Text: prompt})
	if err != nil {
		http.Error(w, "Chat send error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"text":      result.Text(),
		"code":      result.ExecutableCode(),
		"result":    result.CodeExecutionResult(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
