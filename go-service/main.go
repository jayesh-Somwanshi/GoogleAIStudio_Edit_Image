package main

import (
	"log"
	"net/http"
	"github.com/joho/godotenv"
	"go-service/handlers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	http.HandleFunc("/edit-image", handlers.HandleEditImage)
	http.HandleFunc("/text-to-image", handlers.HandleTextToImage)
	http.HandleFunc("/text-to-text", handlers.HandleTextOnly)
	http.HandleFunc("/generate-multi-image",handlers.HandleGenerateMulti)
	http.HandleFunc("/summarize-pdf", handlers.HandlePDFSummary)
	http.HandleFunc("/code-execute",handlers.HandleCodeExecution)
	http.HandleFunc("/image-understanding-by-url",handlers.HandleImageUnderstandingFromUrl)
	http.HandleFunc("/youtube-video-understanding-By-URL",handlers.HandleVideoYoutubeURL)

	log.Println("Go server started on :9095")
	log.Fatal(http.ListenAndServe(":9095", nil))
}
