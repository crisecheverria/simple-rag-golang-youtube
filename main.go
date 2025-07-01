package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
)

func main() {
	fmt.Println("Simple RAG AI Agent")
	fmt.Println("====================")

	// Initialize RAg System
	rag, err := NewRAGSystem()
	if err != nil {
		log.Fatal("Failed to Initialize RAG System:", err)
	}

	fmt.Println("RAG System Initialized successfully!")

	if len(os.Args) > 1 {
		pdfPath := os.Args[1]
		fmt.Printf("Loading PDF: %s\n", pdfPath)

		err = rag.LoadPDF(pdfPath)
		if err != nil {
			log.Fatal("Failed to load PDF:", err)
		}

		fmt.Printf("PDF loaded successfully! Found %d document(s)\n", len(rag.documents))
		if len(rag.documents) > 0 {
			fmt.Printf("Document has %d chunks\n", len(rag.documents[0].Chunks))
		}

	}

	_ = rag
}

type RAGSystem struct {
	documents []Document
}

type Document struct {
	Name    string
	Content string
	Chunks  []string
}

func NewRAGSystem() (*RAGSystem, error) {
	return &RAGSystem{
		documents: make([]Document, 0),
	}, nil
}

// LoadPDF loads and process a PDF file
func (r *RAGSystem) LoadPDF(filePath string) error {
	file, reader, err := pdf.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open PDF: %w", err)
	}
	defer file.Close()

	var content strings.Builder
	totalPages := reader.NumPage()

	for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
		page := reader.Page(pageIndex)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			log.Printf("Warning: failed to extract text from page %d: %v", pageIndex, err)
			continue
		}

		content.WriteString(text)
		content.WriteString("\n")
	}

	chunks := r.chunkText(content.String())
	doc := Document{
		Name:    filePath,
		Content: content.String(),
		Chunks:  chunks,
	}

	r.documents = append(r.documents, doc)

	return nil
}

// chunkText splits text into manageable chunks
func (r *RAGSystem) chunkText(text string) []string {
	const chunkSize = 500
	words := strings.Fields(text)
	var chunks []string

	for i := 0; i < len(words); i += chunkSize {
		end := i + chunkSize
		if end > len(words) {
			end = len(words)
		}
		chunk := strings.Join(words[i:end], " ")
		if len(strings.TrimSpace(chunk)) > 0 {
			chunks = append(chunks, chunk)
		}
	}

	return chunks
}
