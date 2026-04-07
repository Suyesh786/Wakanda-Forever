package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"google.golang.org/api/option"
)

// AIAnalysisRequest represents the incoming JSON for the artifact analysis
type AIAnalysisRequest struct {
	ArtifactName string `json:"artifactName"`
	Era          string `json:"era"`
}

// CustomJumpRequest represents the incoming user query for a new era
type CustomJumpRequest struct {
	Query string `json:"query"`
}

// init loads the .env file locally.
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Relying on system environment variables.")
	}
}

// getErasHandler serves the database list (Calls GetAllEras from db.go)
func getErasHandler(w http.ResponseWriter, r *http.Request) {
	eras := GetAllEras()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eras)
}

// generateAnalysisHandler calls the Google Gemini API for artifact analysis
func generateAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodOptions {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AIAnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		http.Error(w, "API Key not configured", http.StatusInternalServerError)
		return
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		http.Error(w, "Failed to initialize AI client", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Updated to Gemini 2.5 Flash
	model := client.GenerativeModel("gemini-2.5-flash")
	prompt := fmt.Sprintf("Act as a brilliant archaeologist. Give me a 2 sentence, highly detailed structural analysis of the artifact '%s' from the era '%s'. Make it sound technical and dramatic.", req.ArtifactName, req.Era)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil || len(resp.Candidates) == 0 {
		http.Error(w, "AI Generation failed", http.StatusInternalServerError)
		return
	}

	var aiResponse string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			aiResponse += string(text)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"analysis": aiResponse})
}

// customJumpHandler uses Gemini to generate a brand new Era object on the fly
func customJumpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodOptions {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CustomJumpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("ERROR: GEMINI_API_KEY is empty in .env!")
		http.Error(w, "API Key not configured", http.StatusInternalServerError)
		return
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Println("GEMINI CLIENT ERROR:", err)
		http.Error(w, "Failed to initialize AI client", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Updated to Gemini 2.5 Flash
	model := client.GenerativeModel("gemini-2.5-flash")

	prompt := fmt.Sprintf(`
		The user wants to travel in a VR time machine to: "%s". 
		You must generate the historical data for this destination.
		Return ONLY a valid JSON object. Do not include markdown formatting, backticks, or conversational text.
		The JSON MUST exactly match this structure:
		{
			"id": "A single lowercase word (e.g., 'chernobyl', 'moon')",
			"title": "The exact location name",
			"searchTerm": "A valid Wikipedia search term for this place",
			"desc": "A dramatic, atmospheric 2-sentence description of being there at that exact time.",
			"year": "The year with era (e.g. 1986 CE, 2560 BCE)",
			"coord": "LOC: LAT° N/S, LONG° E/W",
			"artName": "Name of a famous artifact, object, or structure present there",
			"artData": "A 2-sentence technical analysis of that artifact",
			"color": "A vibrant hex code (#RRGGBB) that fits the mood of the location",
			"filter": "custom"
		}
	`, req.Query)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Println("GEMINI API ERROR:", err)
		http.Error(w, "AI Generation failed", http.StatusInternalServerError)
		return
	}

	if len(resp.Candidates) == 0 {
		log.Println("GEMINI RETURNED NO CANDIDATES")
		http.Error(w, "AI Generation empty", http.StatusInternalServerError)
		return
	}

	var aiJSON string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			aiJSON += string(text)
		}
	}

	log.Println("--- RAW AI RESPONSE START ---")
	log.Println(aiJSON)
	log.Println("--- RAW AI RESPONSE END ---")

	// ROBUST JSON EXTRACTION: Find the first '{' and the last '}'
	start := strings.Index(aiJSON, "{")
	end := strings.LastIndex(aiJSON, "}")

	if start == -1 || end == -1 {
		log.Println("ERROR: AI Output contained no valid JSON brackets.")
		http.Error(w, "AI output invalid structure", http.StatusInternalServerError)
		return
	}

	// Slice out exactly the JSON payload
	finalJSON := aiJSON[start : end+1]

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, finalJSON)
}

// getWikimediaImageHandler fetches data safely from Wikipedia
func getWikimediaImageHandler(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("term")
	if searchTerm == "" {
		http.Error(w, "Missing search term", http.StatusBadRequest)
		return
	}

	// 1. Wikipedia prefers underscores instead of spaces for titles
	wikiTitle := strings.ReplaceAll(searchTerm, " ", "_")

	// 2. Safely encode any weird characters
	escapedTerm := url.PathEscape(wikiTitle)

	wikiURL := fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&format=json&prop=pageimages&titles=%s&pithumbsize=1000", escapedTerm)

	req, err := http.NewRequest("GET", wikiURL, nil)
	if err != nil {
		log.Println("Failed to create Wikipedia request:", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// 3. THE VIP PASS: Wikipedia strictly requires contact info in the User-Agent!
	req.Header.Set("User-Agent", "ChronosVR/1.0 (chronos.developer@example.com)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Network error contacting Wikipedia:", err)
		http.Error(w, "Failed to contact Wikimedia", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 4. Check if Wikipedia blocked us
	if resp.StatusCode != http.StatusOK {
		log.Printf("Wikipedia rejected request with status: %d for URL: %s", resp.StatusCode, wikiURL)
		http.Error(w, "Wikipedia API error", http.StatusInternalServerError)
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode Wikipedia JSON. URL was: %s", wikiURL)
		http.Error(w, "Failed to parse Wikimedia response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	mux := http.NewServeMux()

	// Register API Routes
	mux.HandleFunc("/api/eras", getErasHandler)
	mux.HandleFunc("/api/analyze", generateAnalysisHandler)
	mux.HandleFunc("/api/wiki-image", getWikimediaImageHandler)
	mux.HandleFunc("/api/custom-jump", customJumpHandler)

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:5500",
			"http://localhost:5500",
			"https://*.vercel.app", // Wildcard for Vercel
		},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true, // Prints exactly what CORS is doing in your terminal
	})

	handler := c.Handler(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("CHRONOS CORE ONLINE // Port: %s\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}