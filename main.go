package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Pokémon স্ট্রাকচার
type Pokemon struct {
	ID            int         `json:"id"`
	Num           string      `json:"num"`
	Name          string      `json:"name"`
	Img           string      `json:"img"`
	Type          []string    `json:"type"`
	Height        string      `json:"height"`
	Weight        string      `json:"weight"`
	Candy         string      `json:"candy"`
	CandyCount    *int        `json:"candy_count,omitempty"`
	Egg           string      `json:"egg"`
	SpawnChance   float64     `json:"spawn_chance"`
	AvgSpawns     float64     `json:"avg_spawns"`
	SpawnTime     string      `json:"spawn_time"`
	Multipliers   []float64   `json:"multipliers"`
	Weaknesses    []string    `json:"weaknesses"`
	NextEvolution []Evolution `json:"next_evolution,omitempty"`
	PrevEvolution []Evolution `json:"prev_evolution,omitempty"`
	CreatedAt     time.Time   `json:"created_at,omitempty"`
	UpdatedAt     time.Time   `json:"updated_at,omitempty"`
}

type Evolution struct {
	Num  string `json:"num"`
	Name string `json:"name"`
}

// ইন-মেমরি ডাটাবেস
type PokemonDB struct {
	sync.RWMutex
	pokemons  []Pokemon
	idCounter int
}

var db *PokemonDB

func init() {
	db = &PokemonDB{
		pokemons:  make([]Pokemon, 0),
		idCounter: 1,
	}
	loadInitialData()
}

// প্রারম্ভিক ডেটা লোড
func loadInitialData() {
	// JSON ফাইল থেকে ডেটা লোড করার চেষ্টা করুন
	if err := loadFromJSON("pokemon.json"); err == nil {
		log.Println("Data loaded from pokemon.json")
		return
	}

	// স্যাম্পল ডেটা লোড
	log.Println("Using sample data")
	sampleData := getSampleData()

	db.Lock()
	defer db.Unlock()

	for i := range sampleData {
		sampleData[i].ID = db.idCounter
		db.idCounter++
		sampleData[i].CreatedAt = time.Now()
		sampleData[i].UpdatedAt = time.Now()
		db.pokemons = append(db.pokemons, sampleData[i])
	}

	log.Printf("Loaded %d Pokémon\n", len(db.pokemons))
}

// JSON ফাইল থেকে লোড
func loadFromJSON(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var pokemons []Pokemon
	if err := json.Unmarshal(data, &pokemons); err != nil {
		return err
	}

	db.Lock()
	defer db.Unlock()

	for i := range pokemons {
		pokemons[i].ID = db.idCounter
		db.idCounter++
		pokemons[i].CreatedAt = time.Now()
		pokemons[i].UpdatedAt = time.Now()
		db.pokemons = append(db.pokemons, pokemons[i])
	}

	return nil
}

// JSON ফাইলে সেভ
func saveToJSON(filename string) error {
	db.RLock()
	defer db.RUnlock()

	data, err := json.MarshalIndent(db.pokemons, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// স্যাম্পল ডেটা
func getSampleData() []Pokemon {
	return []Pokemon{
		{
			Num:           "001",
			Name:          "Bulbasaur",
			Img:           "http://www.serebii.net/pokemongo/pokemon/001.png",
			Type:          []string{"Grass", "Poison"},
			Height:        "0.71 m",
			Weight:        "6.9 kg",
			Candy:         "Bulbasaur Candy",
			CandyCount:    intPtr(25),
			Egg:           "2 km",
			SpawnChance:   0.69,
			AvgSpawns:     69,
			SpawnTime:     "20:00",
			Multipliers:   []float64{1.58},
			Weaknesses:    []string{"Fire", "Ice", "Flying", "Psychic"},
			NextEvolution: []Evolution{{Num: "002", Name: "Ivysaur"}, {Num: "003", Name: "Venusaur"}},
		},
		{
			Num:           "002",
			Name:          "Ivysaur",
			Img:           "http://www.serebii.net/pokemongo/pokemon/002.png",
			Type:          []string{"Grass", "Poison"},
			Height:        "0.99 m",
			Weight:        "13.0 kg",
			Candy:         "Bulbasaur Candy",
			CandyCount:    intPtr(100),
			Egg:           "Not in Eggs",
			SpawnChance:   0.042,
			AvgSpawns:     4.2,
			SpawnTime:     "07:00",
			Multipliers:   []float64{1.2, 1.6},
			Weaknesses:    []string{"Fire", "Ice", "Flying", "Psychic"},
			PrevEvolution: []Evolution{{Num: "001", Name: "Bulbasaur"}},
			NextEvolution: []Evolution{{Num: "003", Name: "Venusaur"}},
		},
		{
			Num:           "003",
			Name:          "Venusaur",
			Img:           "http://www.serebii.net/pokemongo/pokemon/003.png",
			Type:          []string{"Grass", "Poison"},
			Height:        "2.01 m",
			Weight:        "100.0 kg",
			Candy:         "Bulbasaur Candy",
			Egg:           "Not in Eggs",
			SpawnChance:   0.017,
			AvgSpawns:     1.7,
			SpawnTime:     "11:30",
			Weaknesses:    []string{"Fire", "Ice", "Flying", "Psychic"},
			PrevEvolution: []Evolution{{Num: "001", Name: "Bulbasaur"}, {Num: "002", Name: "Ivysaur"}},
		},
		{
			Num:           "004",
			Name:          "Charmander",
			Img:           "http://www.serebii.net/pokemongo/pokemon/004.png",
			Type:          []string{"Fire"},
			Height:        "0.61 m",
			Weight:        "8.5 kg",
			Candy:         "Charmander Candy",
			CandyCount:    intPtr(25),
			Egg:           "2 km",
			SpawnChance:   0.253,
			AvgSpawns:     25.3,
			SpawnTime:     "08:45",
			Multipliers:   []float64{1.65},
			Weaknesses:    []string{"Water", "Ground", "Rock"},
			NextEvolution: []Evolution{{Num: "005", Name: "Charmeleon"}, {Num: "006", Name: "Charizard"}},
		},
		{
			Num:           "005",
			Name:          "Charmeleon",
			Img:           "http://www.serebii.net/pokemongo/pokemon/005.png",
			Type:          []string{"Fire"},
			Height:        "1.09 m",
			Weight:        "19.0 kg",
			Candy:         "Charmander Candy",
			CandyCount:    intPtr(100),
			Egg:           "Not in Eggs",
			SpawnChance:   0.012,
			AvgSpawns:     1.2,
			SpawnTime:     "19:00",
			Multipliers:   []float64{1.79},
			Weaknesses:    []string{"Water", "Ground", "Rock"},
			PrevEvolution: []Evolution{{Num: "004", Name: "Charmander"}},
			NextEvolution: []Evolution{{Num: "006", Name: "Charizard"}},
		},
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

// CORS মিডলওয়্যার
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// JSON রেসপন্স হেল্পার
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// ==================== CRUD OPERATIONS ====================

// 1. CREATE - POST /api/pokemons
func createPokemon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var pokemon Pokemon
	if err := json.NewDecoder(r.Body).Decode(&pokemon); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validation
	if pokemon.Name == "" {
		respondError(w, http.StatusBadRequest, "Name is required")
		return
	}

	db.Lock()
	defer db.Unlock()

	// Check if Pokémon number already exists
	for _, p := range db.pokemons {
		if p.Num == pokemon.Num {
			respondError(w, http.StatusConflict, "Pokemon with this number already exists")
			return
		}
	}

	// Set ID and timestamps
	pokemon.ID = db.idCounter
	db.idCounter++
	pokemon.CreatedAt = time.Now()
	pokemon.UpdatedAt = time.Now()

	db.pokemons = append(db.pokemons, pokemon)

	// Save to file
	if err := saveToJSON("pokemon_backup.json"); err != nil {
		log.Printf("Warning: Could not save backup: %v", err)
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Pokemon created successfully",
		"pokemon": pokemon,
	})
}

// 2. READ ALL - GET /api/pokemons
func getAllPokemons(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Query parameters
	query := r.URL.Query()
	typeFilter := query.Get("type")
	search := query.Get("search")
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	db.RLock()
	defer db.RUnlock()

	// Filtering
	filteredPokemons := make([]Pokemon, len(db.pokemons))
	copy(filteredPokemons, db.pokemons)

	if typeFilter != "" {
		var result []Pokemon
		for _, p := range filteredPokemons {
			for _, t := range p.Type {
				if strings.EqualFold(t, typeFilter) {
					result = append(result, p)
					break
				}
			}
		}
		filteredPokemons = result
	}

	if search != "" {
		var result []Pokemon
		for _, p := range filteredPokemons {
			if strings.Contains(strings.ToLower(p.Name), strings.ToLower(search)) ||
				strings.Contains(p.Num, search) {
				result = append(result, p)
			}
		}
		filteredPokemons = result
	}

	// Pagination
	page := 1
	limit := 20

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	start := (page - 1) * limit
	end := start + limit

	if start > len(filteredPokemons) {
		start = len(filteredPokemons)
	}
	if end > len(filteredPokemons) {
		end = len(filteredPokemons)
	}

	response := map[string]interface{}{
		"total":       len(filteredPokemons),
		"page":        page,
		"limit":       limit,
		"total_pages": (len(filteredPokemons) + limit - 1) / limit,
		"data":        filteredPokemons[start:end],
	}

	respondJSON(w, http.StatusOK, response)
}

// 3. READ ONE - GET /api/pokemons/{id}
func getPokemonByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/pokemons/")
	id, err := strconv.Atoi(path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	db.RLock()
	defer db.RUnlock()

	for _, pokemon := range db.pokemons {
		if pokemon.ID == id {
			respondJSON(w, http.StatusOK, pokemon)
			return
		}
	}

	respondError(w, http.StatusNotFound, "Pokemon not found")
}

// 4. UPDATE - PUT /api/pokemons/{id}
func updatePokemon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/pokemons/")
	id, err := strconv.Atoi(path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var updatedPokemon Pokemon
	if err := json.NewDecoder(r.Body).Decode(&updatedPokemon); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	db.Lock()
	defer db.Unlock()

	for i, pokemon := range db.pokemons {
		if pokemon.ID == id {
			// Keep original ID and creation timestamp
			updatedPokemon.ID = pokemon.ID
			updatedPokemon.CreatedAt = pokemon.CreatedAt
			updatedPokemon.UpdatedAt = time.Now()

			db.pokemons[i] = updatedPokemon

			// Save backup
			if err := saveToJSON("pokemon_backup.json"); err != nil {
				log.Printf("Warning: Could not save backup: %v", err)
			}

			respondJSON(w, http.StatusOK, map[string]interface{}{
				"message": "Pokemon updated successfully",
				"pokemon": updatedPokemon,
			})
			return
		}
	}

	respondError(w, http.StatusNotFound, "Pokemon not found")
}

// 5. PARTIAL UPDATE - PATCH /api/pokemons/{id}
func patchPokemon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/pokemons/")
	id, err := strconv.Atoi(path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	db.Lock()
	defer db.Unlock()

	for i, pokemon := range db.pokemons {
		if pokemon.ID == id {
			// Convert to JSON and back to apply updates
			pokemonJSON, _ := json.Marshal(pokemon)
			var pokemonMap map[string]interface{}
			json.Unmarshal(pokemonJSON, &pokemonMap)

			// Apply updates
			for key, value := range updates {
				if key != "id" && key != "created_at" && key != "updated_at" {
					pokemonMap[key] = value
				}
			}

			// Convert back to Pokemon struct
			updatedJSON, _ := json.Marshal(pokemonMap)
			var updatedPokemon Pokemon
			json.Unmarshal(updatedJSON, &updatedPokemon)

			// Restore original ID and timestamps
			updatedPokemon.ID = pokemon.ID
			updatedPokemon.CreatedAt = pokemon.CreatedAt
			updatedPokemon.UpdatedAt = time.Now()

			db.pokemons[i] = updatedPokemon

			// Save backup
			if err := saveToJSON("pokemon_backup.json"); err != nil {
				log.Printf("Warning: Could not save backup: %v", err)
			}

			respondJSON(w, http.StatusOK, map[string]interface{}{
				"message": "Pokemon updated successfully",
				"pokemon": updatedPokemon,
			})
			return
		}
	}

	respondError(w, http.StatusNotFound, "Pokemon not found")
}

// 6. DELETE - DELETE /api/pokemons/{id}
func deletePokemon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/pokemons/")
	id, err := strconv.Atoi(path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	db.Lock()
	defer db.Unlock()

	for i, pokemon := range db.pokemons {
		if pokemon.ID == id {
			// Remove the element
			db.pokemons = append(db.pokemons[:i], db.pokemons[i+1:]...)

			// Save backup
			if err := saveToJSON("pokemon_backup.json"); err != nil {
				log.Printf("Warning: Could not save backup: %v", err)
			}

			respondJSON(w, http.StatusOK, map[string]string{
				"message": "Pokemon deleted successfully",
				"id":      fmt.Sprintf("%d", id),
			})
			return
		}
	}

	respondError(w, http.StatusNotFound, "Pokemon not found")
}

// 7. BULK CREATE - POST /api/pokemons/bulk
func bulkCreatePokemons(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var newPokemons []Pokemon
	if err := json.NewDecoder(r.Body).Decode(&newPokemons); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	db.Lock()
	defer db.Unlock()

	var created []Pokemon
	var errors []string

	for _, pokemon := range newPokemons {
		if pokemon.Name == "" {
			errors = append(errors, "Pokemon name cannot be empty")
			continue
		}

		// Check if Pokémon number already exists
		exists := false
		for _, p := range db.pokemons {
			if p.Num == pokemon.Num {
				errors = append(errors, fmt.Sprintf("Pokemon with number %s already exists", pokemon.Num))
				exists = true
				break
			}
		}
		if exists {
			continue
		}

		// Set ID and timestamps
		pokemon.ID = db.idCounter
		db.idCounter++
		pokemon.CreatedAt = time.Now()
		pokemon.UpdatedAt = time.Now()

		db.pokemons = append(db.pokemons, pokemon)
		created = append(created, pokemon)
	}

	// Save backup
	if err := saveToJSON("pokemon_backup.json"); err != nil {
		log.Printf("Warning: Could not save backup: %v", err)
	}

	response := map[string]interface{}{
		"message":          "Bulk create completed",
		"created_count":    len(created),
		"failed_count":     len(errors),
		"created_pokemons": created,
		"errors":           errors,
	}

	status := http.StatusCreated
	if len(errors) > 0 {
		status = http.StatusPartialContent
	}

	respondJSON(w, status, response)
}

// 8. DELETE ALL - DELETE /api/pokemons
func deleteAllPokemons(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query()
	confirm := query.Get("confirm")

	if confirm != "true" {
		respondError(w, http.StatusBadRequest, "Add ?confirm=true to confirm deletion")
		return
	}

	db.Lock()
	defer db.Unlock()

	count := len(db.pokemons)
	db.pokemons = make([]Pokemon, 0)
	db.idCounter = 1

	// Save empty backup
	if err := saveToJSON("pokemon_backup.json"); err != nil {
		log.Printf("Warning: Could not save backup: %v", err)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":       "All Pokemon deleted successfully",
		"deleted_count": count,
	})
}

// 9. SPECIAL ENDPOINTS
func getPokemonsByType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/pokemons/type/")
	if path == "" {
		respondError(w, http.StatusBadRequest, "Type parameter required")
		return
	}

	db.RLock()
	defer db.RUnlock()

	var result []Pokemon
	for _, pokemon := range db.pokemons {
		for _, t := range pokemon.Type {
			if strings.EqualFold(t, path) {
				result = append(result, pokemon)
				break
			}
		}
	}

	respondJSON(w, http.StatusOK, result)
}

func getPokemonsWeakAgainst(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/pokemons/weakness/")
	if path == "" {
		respondError(w, http.StatusBadRequest, "Weakness parameter required")
		return
	}

	db.RLock()
	defer db.RUnlock()

	var result []Pokemon
	for _, pokemon := range db.pokemons {
		for _, w := range pokemon.Weaknesses {
			if strings.EqualFold(w, path) {
				result = append(result, pokemon)
				break
			}
		}
	}

	respondJSON(w, http.StatusOK, result)
}

// 10. STATISTICS - GET /api/stats
func getStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	db.RLock()
	defer db.RUnlock()

	// Calculate statistics
	typeStats := make(map[string]int)
	totalSpawnChance := 0.0
	highestSpawn := db.pokemons[0]
	lowestSpawn := db.pokemons[0]

	for _, pokemon := range db.pokemons {
		// Type statistics
		for _, t := range pokemon.Type {
			typeStats[t]++
		}

		// Spawn statistics
		totalSpawnChance += pokemon.SpawnChance
		if pokemon.SpawnChance > highestSpawn.SpawnChance {
			highestSpawn = pokemon
		}
		if pokemon.SpawnChance < lowestSpawn.SpawnChance {
			lowestSpawn = pokemon
		}
	}

	avgSpawnChance := totalSpawnChance / float64(len(db.pokemons))

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"total_pokemons":   len(db.pokemons),
		"by_type":          typeStats,
		"avg_spawn_chance": avgSpawnChance,
		"highest_spawn":    highestSpawn,
		"lowest_spawn":     lowestSpawn,
		"last_updated":     time.Now().Format(time.RFC3339),
	})
}

// 11. SEARCH - GET /api/pokemons/search/{query}
func searchPokemons(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/pokemons/search/")
	if path == "" {
		respondError(w, http.StatusBadRequest, "Search query required")
		return
	}

	db.RLock()
	defer db.RUnlock()

	var result []Pokemon
	query := strings.ToLower(path)

	for _, pokemon := range db.pokemons {
		if strings.Contains(strings.ToLower(pokemon.Name), query) ||
			strings.Contains(pokemon.Num, query) ||
			strings.Contains(strings.ToLower(pokemon.Candy), query) {
			result = append(result, pokemon)
		}
	}

	respondJSON(w, http.StatusOK, result)
}

// 12. HOME HANDLER
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	response := map[string]interface{}{
		"message": "Pokémon REST API with Full CRUD Operations",
		"version": "2.0.0",
		"endpoints": map[string]string{
			"GET    /":          "API Documentation",
			"GET    /api/stats": "Get Pokémon statistics",

			// CRUD Operations
			"GET    /api/pokemons":      "Get all Pokémon (with pagination)",
			"POST   /api/pokemons":      "Create new Pokémon",
			"GET    /api/pokemons/{id}": "Get Pokémon by ID",
			"PUT    /api/pokemons/{id}": "Update Pokémon (full)",
			"PATCH  /api/pokemons/{id}": "Update Pokémon (partial)",
			"DELETE /api/pokemons/{id}": "Delete Pokémon by ID",

			// Bulk Operations
			"POST   /api/pokemons/bulk":         "Bulk create Pokémon",
			"DELETE /api/pokemons?confirm=true": "Delete all Pokémon",

			// Special Queries
			"GET    /api/pokemons/type/{type}":     "Get Pokémon by type",
			"GET    /api/pokemons/weakness/{type}": "Get Pokémon weak against",
			"GET    /api/pokemons/search/{query}":  "Search Pokémon",

			// Query Parameters
			"?type=Fire":       "Filter by type",
			"?search=pika":     "Search by name/ID",
			"?page=2&limit=10": "Pagination",
		},
		"example_payload": map[string]interface{}{
			"name":         "Pikachu",
			"num":          "025",
			"type":         []string{"Electric"},
			"height":       "0.41 m",
			"weight":       "6.0 kg",
			"candy":        "Pikachu Candy",
			"candy_count":  50,
			"egg":          "2 km",
			"spawn_chance": 0.21,
			"avg_spawns":   21,
			"spawn_time":   "04:00",
			"weaknesses":   []string{"Ground"},
		},
	}

	respondJSON(w, http.StatusOK, response)
}

// ==================== MAIN FUNCTION ====================
func main() {
	// রাউটিং সেটআপ
	http.HandleFunc("/", enableCORS(homeHandler))

	// Stats
	http.HandleFunc("/api/stats", enableCORS(getStats))

	// Main Pokémon endpoints
	http.HandleFunc("/api/pokemons", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAllPokemons(w, r)
		case http.MethodPost:
			createPokemon(w, r)
		case http.MethodDelete:
			deleteAllPokemons(w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	// Bulk operations
	http.HandleFunc("/api/pokemons/bulk", enableCORS(bulkCreatePokemons))

	// Individual Pokémon operations
	http.HandleFunc("/api/pokemons/", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/pokemons/")

		// Check for special endpoints first
		if strings.HasPrefix(path, "type/") {
			getPokemonsByType(w, r)
			return
		}
		if strings.HasPrefix(path, "weakness/") {
			getPokemonsWeakAgainst(w, r)
			return
		}
		if strings.HasPrefix(path, "search/") {
			searchPokemons(w, r)
			return
		}

		// Check if it's just "/api/pokemons/"
		if path == "" {
			switch r.Method {
			case http.MethodGet:
				getAllPokemons(w, r)
			case http.MethodPost:
				createPokemon(w, r)
			default:
				respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
			return
		}

		// Regular ID-based operations
		switch r.Method {
		case http.MethodGet:
			getPokemonByID(w, r)
		case http.MethodPut:
			updatePokemon(w, r)
		case http.MethodPatch:
			patchPokemon(w, r)
		case http.MethodDelete:
			deletePokemon(w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	// Start server
	port := ":8080"

	fmt.Println("=====================================")
	fmt.Printf("📡 Server URL: http://localhost%s\n", port)
	fmt.Println("📚 API Documentation: http://localhost" + port)
	fmt.Println("\n📋 Available Endpoints:")
	fmt.Println("  GET    /                         - API Documentation")
	fmt.Println("  GET    /api/stats                - Statistics")
	fmt.Println("  GET    /api/pokemons             - Get all Pokémon")
	fmt.Println("  POST   /api/pokemons             - Create new Pokémon")
	fmt.Println("  GET    /api/pokemons/{id}        - Get Pokémon by ID")
	fmt.Println("  PUT    /api/pokemons/{id}        - Update Pokémon (full)")
	fmt.Println("  PATCH  /api/pokemons/{id}        - Update Pokémon (partial)")
	fmt.Println("  DELETE /api/pokemons/{id}        - Delete Pokémon by ID")
	fmt.Println("  POST   /api/pokemons/bulk        - Bulk create Pokémon")
	fmt.Println("  DELETE /api/pokemons?confirm=true - Delete all Pokémon")
	fmt.Println("  GET    /api/pokemons/type/{type} - Filter by type")
	fmt.Println("  GET    /api/pokemons/search/{query} - Search Pokémon")
	fmt.Println("\n⚡ Example curl commands saved to curl_examples.txt")

	// Start server
	log.Fatal(http.ListenAndServe(port, nil))
}
