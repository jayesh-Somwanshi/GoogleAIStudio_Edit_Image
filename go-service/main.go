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
	http.HandleFunc("/code-execute",handlers.HandleCodeExecution)
	http.HandleFunc("/image-understanding-by-url",handlers.HandleImageUnderstandingFromUrl)
    http.HandleFunc("/summarize-youtubeVideo-from-url",handlers.HandleVideoYoutubeURL)
	http.HandleFunc("/summarize-pdf-from-locally-stored", handlers.HandlePDFLocallySummary)
	http.HandleFunc("/summarize-pdf-from-url",handlers.HandlerSummarizePdfFromUrl)
	http.HandleFunc("/summarize-image-from-locally",handlers.HandleImageLocallySummary)
	http.HandleFunc("/compare-two-images",handlers.HandlerCompareTwoImages)

	log.Println("Go server started on :9095")
	log.Fatal(http.ListenAndServe(":9095", nil))
}
