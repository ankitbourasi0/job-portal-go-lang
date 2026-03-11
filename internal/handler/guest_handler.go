package handler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ankitbourasi0/job-portal/internal/repository"
)

type GuestHandler struct {
	Repo *repository.GuestRepository
}

func (h GuestHandler) HandleAnalyzeResumeForAtsScore(w http.ResponseWriter, r *http.Request) {
	//1. Max Memory Limit
	r.ParseMultipartForm(5 << 20) //limit, server will reject if file is greater than 5MB

	//2.Get file from FORM(key should be 'resume' or something else)
	file, header, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, "File too large, it should be less than 5MB", http.StatusRequestEntityTooLarge)
		return
	}
	defer file.Close()

	//3. Validation - Check its actually Pdf
	if header.Header.Get("Content-Type") != "application/pdf" {
		http.Error(w, "Only pdf files are allowed", http.StatusUnsupportedMediaType)
		return
	}
	//4. Read File content
	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "File read failed", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Received file: %s, size: %d bytes\n", header.Filename, len(content))

	//5. Send 'content' (byte slice) to PDF parse function
	resumeText, err := h.Repo.PdfParser(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		log.Printf("PDF Parsing Error: %v", err)
		http.Error(w, "Failed to parse PDF", http.StatusInternalServerError)
		return
	}

	//Gemini API and Scoring Logic
	// Now send this resumeText to Gemini!
	fmt.Println("Extracted Text Length:", len(resumeText))
}
