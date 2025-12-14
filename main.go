package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
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
	CandyCount    int         `json:"candy_count,omitempty"`
	Egg           string      `json:"egg"`
	SpawnChance   float64     `json:"spawn_chance"`
	AvgSpawns     float64     `json:"avg_spawns"`
	SpawnTime     string      `json:"spawn_time"`
	Multipliers   []float64   `json:"multipliers"`
	Weaknesses    []string    `json:"weaknesses"`
	NextEvolution []Evolution `json:"next_evolution,omitempty"`
	PrevEvolution []Evolution `json:"prev_evolution,omitempty"`
}

type Evolution struct {
	Num  string `json:"num"`
	Name string `json:"name"`
}

// ইন-মেমরি ডাটা স্টোর
var pokemons []Pokemon

func init() {
	// আপনার JSON ডেটা এখানে লোড করুন
	// আমি উদাহরণ হিসেবে কিছু ডেটা দিচ্ছি
	pokemons = []Pokemon{
		{
			ID:          1,
			Num:         "001",
			Name:        "Bulbasaur",
			Img:         "http://www.serebii.net/pokemongo/pokemon/001.png",
			Type:        []string{"Grass", "Poison"},
			Height:      "0.71 m",
			Weight:      "6.9 kg",
			Candy:       "Bulbasaur Candy",
			CandyCount:  25,
			Egg:         "2 km",
			SpawnChance: 0.69,
			AvgSpawns:   69,
			SpawnTime:   "20:00",
			Multipliers: []float64{1.58},
			Weaknesses:  []string{"Fire", "Ice", "Flying", "Psychic"},
			NextEvolution: []Evolution{
				{Num: "002", Name: "Ivysaur"},
				{Num: "003", Name: "Venusaur"},
			},
		},
		{
			ID:          2,
			Num:         "002",
			Name:        "Ivysaur",
			Img:         "http://www.serebii.net/pokemongo/pokemon/002.png",
			Type:        []string{"Grass", "Poison"},
			Height:      "0.99 m",
			Weight:      "13.0 kg",
			Candy:       "Bulbasaur Candy",
			CandyCount:  100,
			Egg:         "Not in Eggs",
			SpawnChance: 0.042,
			AvgSpawns:   4.2,
			SpawnTime:   "07:00",
			Multipliers: []float64{1.2, 1.6},
			Weaknesses:  []string{"Fire", "Ice", "Flying", "Psychic"},
			PrevEvolution: []Evolution{
				{Num: "001", Name: "Bulbasaur"},
			},
			NextEvolution: []Evolution{
				{Num: "003", Name: "Venusaur"},
			},
		},
		// আপনার সমস্ত ডেটা এখানে যোগ করুন
	}

	// ফাইলে থেকে ডেটা লোড করার জন্য (যদি চান)
	// loadFromJSON("pokemon.json")
}

// JSON থেকে ডেটা লোড করার ফাংশন
func loadFromJSON(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &pokemons)
	if err != nil {
		return err
	}

	return nil
}

// CORS মিডলওয়্যার
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// 1. GET all Pokémon
func getAllPokemons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Query parameters
	query := r.URL.Query()
	typeFilter := query.Get("type")
	search := query.Get("search")
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	// Filtering
	filteredPokemons := pokemons

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

	json.NewEncoder(w).Encode(response)
}

// 2. GET Pokémon by ID
func getPokemonByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/pokemons/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	for _, pokemon := range pokemons {
		if pokemon.ID == id {
			json.NewEncoder(w).Encode(pokemon)
			return
		}
	}

	http.Error(w, `{"error": "Pokemon not found"}`, http.StatusNotFound)
}

// 3. GET Pokémon by type
func getPokemonsByType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, `{"error": "Type parameter required"}`, http.StatusBadRequest)
		return
	}

	pokemonType := pathParts[3]
	var result []Pokemon

	for _, pokemon := range pokemons {
		for _, t := range pokemon.Type {
			if strings.EqualFold(t, pokemonType) {
				result = append(result, pokemon)
				break
			}
		}
	}

	json.NewEncoder(w).Encode(result)
}

// 4. GET weak against
func getPokemonsWeakAgainst(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, `{"error": "Weakness parameter required"}`, http.StatusBadRequest)
		return
	}

	weakness := pathParts[3]
	var result []Pokemon

	for _, pokemon := range pokemons {
		for _, w := range pokemon.Weaknesses {
			if strings.EqualFold(w, weakness) {
				result = append(result, pokemon)
				break
			}
		}
	}

	json.NewEncoder(w).Encode(result)
}

// 5. GET by candy type
func getPokemonsByCandy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, `{"error": "Candy parameter required"}`, http.StatusBadRequest)
		return
	}

	candyType := pathParts[3]
	var result []Pokemon

	for _, pokemon := range pokemons {
		if strings.Contains(strings.ToLower(pokemon.Candy), strings.ToLower(candyType)) {
			result = append(result, pokemon)
		}
	}

	json.NewEncoder(w).Encode(result)
}

// 6. GET by egg distance
func getPokemonsByEgg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, `{"error": "Egg parameter required"}`, http.StatusBadRequest)
		return
	}

	eggDistance := pathParts[3]
	var result []Pokemon

	for _, pokemon := range pokemons {
		if pokemon.Egg == eggDistance {
			result = append(result, pokemon)
		}
	}

	json.NewEncoder(w).Encode(result)
}

// 7. GET spawn chance filter
func getPokemonsBySpawnChance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	minStr := query.Get("min")
	maxStr := query.Get("max")

	var min, max float64 = 0, 100

	if minStr != "" {
		if m, err := strconv.ParseFloat(minStr, 64); err == nil {
			min = m
		}
	}

	if maxStr != "" {
		if m, err := strconv.ParseFloat(maxStr, 64); err == nil {
			max = m
		}
	}

	var result []Pokemon
	for _, pokemon := range pokemons {
		if pokemon.SpawnChance >= min && pokemon.SpawnChance <= max {
			result = append(result, pokemon)
		}
	}

	// Sort by spawn chance (descending)
	sort.Slice(result, func(i, j int) bool {
		return result[i].SpawnChance > result[j].SpawnChance
	})

	json.NewEncoder(w).Encode(result)
}

// 8. GET by height/weight range
func getPokemonsBySize(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	heightMin := query.Get("height_min")
	heightMax := query.Get("height_max")
	weightMin := query.Get("weight_min")
	weightMax := query.Get("weight_max")

	var result []Pokemon

	for _, pokemon := range pokemons {
		height, _ := strconv.ParseFloat(strings.Split(pokemon.Height, " ")[0], 64)
		weight, _ := strconv.ParseFloat(strings.Split(pokemon.Weight, " ")[0], 64)

		heightMatch := true
		weightMatch := true

		if heightMin != "" {
			if hmin, err := strconv.ParseFloat(heightMin, 64); err == nil {
				heightMatch = heightMatch && (height >= hmin)
			}
		}

		if heightMax != "" {
			if hmax, err := strconv.ParseFloat(heightMax, 64); err == nil {
				heightMatch = heightMatch && (height <= hmax)
			}
		}

		if weightMin != "" {
			if wmin, err := strconv.ParseFloat(weightMin, 64); err == nil {
				weightMatch = weightMatch && (weight >= wmin)
			}
		}

		if weightMax != "" {
			if wmax, err := strconv.ParseFloat(weightMax, 64); err == nil {
				weightMatch = weightMatch && (weight <= wmax)
			}
		}

		if heightMatch && weightMatch {
			result = append(result, pokemon)
		}
	}

	json.NewEncoder(w).Encode(result)
}

// রাউট সেটআপ
func setupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/pokemons", getAllPokemons)
	mux.HandleFunc("/api/pokemons/type/", getPokemonsByType)
	mux.HandleFunc("/api/pokemons/weakness/", getPokemonsWeakAgainst)
	mux.HandleFunc("/api/pokemons/candy/", getPokemonsByCandy)
	mux.HandleFunc("/api/pokemons/egg/", getPokemonsByEgg)
	mux.HandleFunc("/api/pokemons/spawn", getPokemonsBySpawnChance)
	mux.HandleFunc("/api/pokemons/size", getPokemonsBySize)
	mux.HandleFunc("/api/pokemons/", getPokemonByID) // Last because it catches /api/pokemons/{id}
}

func main() {
	mux := http.NewServeMux()

	// স্ট্যাটিক ফাইল সার্ভ (যদি ফ্রন্টএন্ড থাকে)
	// mux.Handle("/", http.FileServer(http.Dir("./static")))

	// API রাউটস
	setupRoutes(mux)

	// CORS সহ সার্ভার
	handler := enableCORS(mux)

	// স্বাগত বার্তা রুট
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message":   "Pokémon API is running!",
				"version":   "1.0.0",
				"endpoints": "Available at /api/pokemons",
			})
		}
	})

	log.Println("Server starting on port 8080...")
	log.Println("Available endpoints:")
	log.Println("  GET /api/pokemons - Get all Pokémon with pagination")
	log.Println("  GET /api/pokemons/{id} - Get Pokémon by ID")
	log.Println("  GET /api/pokemons/type/{type} - Get Pokémon by type")
	log.Println("  GET /api/pokemons/weakness/{type} - Get Pokémon weak against")
	log.Println("  GET /api/pokemons/candy/{name} - Get Pokémon by candy")
	log.Println("  GET /api/pokemons/egg/{distance} - Get Pokémon by egg distance")
	log.Println("  GET /api/pokemons/spawn?min=0.1&max=1.0 - Filter by spawn chance")
	log.Println("  GET /api/pokemons/size?height_min=0.5&height_max=2.0 - Filter by size")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
