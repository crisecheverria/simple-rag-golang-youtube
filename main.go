package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/philippgille/chromem-go"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
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

			// Interactive query mode
			fmt.Println("\n=== Query Mode ===")
			fmt.Println("You can now ask questions about the document.")
			fmt.Println("Type 'quit' to exit.")

			scanner := bufio.NewScanner(os.Stdin)
			for {
				fmt.Print("\nQuestion: ")
				if !scanner.Scan() {
					break
				}

				question := strings.TrimSpace(scanner.Text())
				if question == "quit" {
					break
				}

				if question == "" {
					continue
				}

				fmt.Println("Generating answer...")
				ctx := context.Background()
				answer, err := rag.Query(ctx, question)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					continue
				}

				fmt.Printf("\nAnswer: %s\n", answer)

			}
		}

	} else {
		fmt.Println("Usage: go tun main.go <path-to-file>")
		fmt.Println("Example: go run main.go ./sample.pdf")
		fmt.Println("\nPrerquisites:")
		fmt.Println("1. Install Ollama: https:ollama.ai")
		fmt.Println("2. Pull a model: ollama pull llama2")
		fmt.Println("3. Ensure Ollama is running: ollama serve")
	}

	_ = rag
}

type RAGSystem struct {
	documents  []Document
	llm        llms.Model
	vectorDB   *chromem.DB
	collection *chromem.Collection
}

type Document struct {
	Name    string
	Content string
	Chunks  []string
}

func NewRAGSystem() (*RAGSystem, error) {
	// Initialize puer Go vector database (chromen-go)
	db := chromem.NewDB()

	// Create collection for documents
	collection, err := db.CreateCollection("rag-documents", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create vector collection: %w", err)
	}

	// Initialize Ollama LLMs
	llm, err := ollama.New(ollama.WithModel("llama2"))
	if err != nil {
		return nil, fmt.Errorf("failed to Initialize Ollama: %w", err)
	}

	return &RAGSystem{
		documents:  make([]Document, 0),
		llm:        llm,
		vectorDB:   db,
		collection: collection,
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

	// Store chunks in vector database
	ctx := context.Background()
	var ids []string
	var contents []string
	var metadatas []map[string]string

	for i, chunk := range chunks {
		ids = append(ids, fmt.Sprintf("%s_chunk_%d", filePath, i))
		contents = append(contents, chunk)
		metadatas = append(metadatas, map[string]string{
			"source":   filePath,
			"chunk_id": fmt.Sprintf("%d", i),
		})
	}

	err = r.collection.Add(ctx, ids, nil, metadatas, contents)
	if err != nil {
		return fmt.Errorf("failed to store chunks in vector database: %w", err)
	}

	fmt.Printf("Stored %d chunks in pure Go vector database\n", len(chunks))

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

// Query performs a semantic query using vector search and Ollama
func (r *RAGSystem) Query(ctx context.Context, question string) (string, error) {
	if len(r.documents) == 0 {
		return "", fmt.Errorf("no documents loaded")
	}

	// Perform semantic search using pure Go vector database
	results, err := r.collection.Query(ctx, question, 3, nil, nil)
	if err != nil {
		// fallback to keyword matching if vector search fails
		fmt.Println("Vector search failedm failinng back to keyword matching...")
		relevantChunks := r.findRelevantChunks(question, 3)
		return r.generateAnswer(ctx, question, relevantChunks)
	}

	var relevantChunks []string
	for i, result := range results {
		fmt.Printf("Chunk %d preview: %.100s...\n", i+1, result.Content)
		relevantChunks = append(relevantChunks, result.Content)
	}

	return r.generateAnswer(ctx, question, relevantChunks)
}

// generateAnswer creates the final answer using llms
func (r *RAGSystem) generateAnswer(ctx context.Context, question string, chunks []string) (string, error) {
	// Build context from relevant chunks
	context := "Context from documents:\n"
	for i, chunk := range chunks {
		context += fmt.Sprintf("%d. %s\n\n", i+1, chunk)
	}

	// Create prompt
	prompt := fmt.Sprintf(`Based on the following context, answer the question: %s Question: %s Answer:`, context, question)
	// Query the LLM
	response, err := r.llm.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	return response.Choices[0].Content, nil
}

// findRelevantChunks performs simple keyword-based retrieval
func (r *RAGSystem) findRelevantChunks(query string, maxChunks int) []string {
	queryWords := strings.Fields(strings.ToLower(query))

	type chunkScore struct {
		chunk string
		score int
	}

	var scored []chunkScore

	// Score chunks based in keyword overlap
	for _, doc := range r.documents {
		for _, chunk := range doc.Chunks {
			chunkLower := strings.ToLower(chunk)
			score := 0

			for _, word := range queryWords {
				if strings.Contains(chunkLower, word) {
					score++
				}
			}

			if score > 0 {
				scored = append(scored, chunkScore{chunk: chunk, score: score})
			}
		}
	}

	// Sort by score (simple bubble sort for simplicity)
	for i := range scored {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// Return top chunks
	var result []string
	limit := min(len(scored), maxChunks)

	for i := range limit {
		result = append(result, scored[i].chunk)
	}

	return result
}
