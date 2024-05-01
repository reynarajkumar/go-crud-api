package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var profiles []Profile = []Profile{}

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}
type Profile struct {
	Department  string `json:"department"`
	Designation string `json:"designation"`
	Employee    User   `json:"employee"`
}

func addItem(q http.ResponseWriter, r *http.Request) {
	var newProfile Profile
	json.NewDecoder(r.Body).Decode(&newProfile)

	q.Header().Set("Content-Type", "application/json")

	profiles = append(profiles, newProfile)

	json.NewEncoder(q).Encode(profiles)
}

func getAllProfiles(q http.ResponseWriter, r *http.Request) {
	q.Header().Set("Content-Type", "application/json")
	json.NewEncoder(q).Encode(profiles)
}

func getProfile(q http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		q.WriteHeader(400)
		q.Write([]byte("ID could not be converted to integer"))
		return
	}
	if id >= len(profiles) {
		q.WriteHeader(404)
		q.Write([]byte("no profile found with specified ID"))
	}
	profile := profiles[id]
	q.Header().Set("Content-Type", "application/json")
	json.NewEncoder(q).Encode(profile)
}

func updateProfile(q http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		q.WriteHeader(400)
		q.Write([]byte("ID could not be converted to integer"))
		return
	}
	if id >= len(profiles) {
		q.WriteHeader(404)
		q.Write([]byte("no profile found with specified ID"))
	}
	var updateProfile Profile
	json.NewDecoder(r.Body).Decode(&updateProfile)
	profiles[id] = updateProfile
	q.Header().Set("Content-Type", "application/json")
	json.NewEncoder(q).Encode(updateProfile)
}

func deleteProfile(q http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		q.WriteHeader(400)
		q.Write([]byte("ID could not be converted to integer"))
		return
	}
	if id >= len(profiles) {
		q.WriteHeader(404)
		q.Write([]byte("no profile found with specified ID"))
	}
	profiles = append(profiles[:id], profiles[id+1:]...)
	q.WriteHeader(200)

}

// makes a POST request with unique name
func makeAPICallParallelized(path string, body []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Post("http://localhost:8080"+path, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error making API call: %v\n", err)
		return
	}
	defer resp.Body.Close()
}

// makes a POST request with unique name
func makeAPICallNotParallelized(path string, body []byte) {
	resp, err := http.Post("http://localhost:8080"+path, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error making API call: %v\n", err)
		return
	}
	defer resp.Body.Close()
}

func generateRandJSONProfile(firstName string) []byte {
	profile := Profile{
		Department:  "Engineering",
		Designation: "Software Engineer",
		Employee: User{
			FirstName: firstName,
			LastName:  "Doe",
			Email:     fmt.Sprintf("%s@gmail.com", firstName),
		},
	}

	jsonProfile, _ := json.Marshal(profile) // convert profile struct into JSON
	return jsonProfile
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/profiles", addItem).Methods("POST")
	router.HandleFunc("/profiles", getAllProfiles).Methods("GET")
	router.HandleFunc("/profiles/{id}", getProfile).Methods("GET")
	router.HandleFunc("/profiles/{id}", updateProfile).Methods("PUT")
	router.HandleFunc("/profiles/{id}", deleteProfile).Methods("DELETE")

	go http.ListenAndServe(":8080", router)

	// Parallelized API calls
	startTime := time.Now()

	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup
	numCalls := 100
	wg.Add(numCalls)

	for i := 0; i < numCalls; i++ {
		firstName := fmt.Sprintf("John%d", i+1)      // Generate unique first name
		person := generateRandJSONProfile(firstName) // Generate JSON person with unique first name
		go makeAPICallParallelized("/profiles", person, &wg)
	}

	// Wait for all API calls to finish
	wg.Wait()

	totalTime := time.Since(startTime)
	fmt.Println("Time taken for", numCalls, "API calls parallelized:", totalTime)

	// Unparallelized API calls
	startTime2 := time.Now()

	for i := 0; i < numCalls; i++ {
		firstName := fmt.Sprintf("John%d", i+1)      // Generate unique first name
		person := generateRandJSONProfile(firstName) // Generate JSON payload with unique first name
		makeAPICallNotParallelized("/profiles", person)
	}

	timeTaken2 := time.Since(startTime2)
	fmt.Println("Time taken for", numCalls, "API calls unparallelized:", timeTaken2)

	fmt.Println("Press Ctrl+C to stop the server.")
	select {}

}
