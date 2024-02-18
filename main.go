package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// GameResult struct to hold data sent back to the frontend
type GameResult struct {
	PlayerScore    int    `json:"playerScore"`
	ComputerScore  int    `json:"computerScore"`
	Message        string `json:"message"`
	ComputerChoice string `json:"computerChoice"` // Field to hold the computer's choice
}

// Declare the playerScore and computerScore variables as global
var playerScore int
var computerScore int

// Middleware to handle CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")                                // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE") // Allowed methods
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass down the request to the next handler
		next.ServeHTTP(w, r)
	})
}

func generateComputerChoice() string {
	choices := []string{"rock", "paper", "scissors"}
	randomIndex := rand.Intn(len(choices))
	return choices[randomIndex]
}

func determineWinner(playerChoice, computerChoice string) (string, string) {
	// Messages for different game outcomes
	tieMessages := []string{
		"It's a tie!",
		"Nobody wins this round.",
		"You both chose the same thing.",
	}

	playerWinMessages := []string{
		"You win! You have bested the computer with your skills.",
		"The computer is no match for you. You rock!",
		"You have defeated the computer. Congratulations!",
	}

	computerWinMessages := []string{
		"The computer wins this round. Better luck next time!",
		"You have been defeated by the computer. Try again!",
	}

	// Logic to determine the winner and select an appropriate message
	var message string
	switch playerChoice {
	case "rock":
		switch computerChoice {
		case "rock":
			message = tieMessages[rand.Intn(len(tieMessages))]
		case "paper":
			computerScore++
			message = computerWinMessages[rand.Intn(len(computerWinMessages))]
		case "scissors":
			playerScore++
			message = playerWinMessages[rand.Intn(len(playerWinMessages))]
		}
	case "paper":
		switch computerChoice {
		case "rock":
			playerScore++
			message = playerWinMessages[rand.Intn(len(playerWinMessages))]
		case "paper":
			message = tieMessages[rand.Intn(len(tieMessages))]
		case "scissors":
			computerScore++
			message = computerWinMessages[rand.Intn(len(computerWinMessages))]
		}
	case "scissors":
		switch computerChoice {
		case "rock":
			computerScore++
			message = computerWinMessages[rand.Intn(len(computerWinMessages))]
		case "paper":
			playerScore++
			message = playerWinMessages[rand.Intn(len(playerWinMessages))]
		case "scissors":
			message = tieMessages[rand.Intn(len(tieMessages))]
		}
	default:
		message = "Invalid choice. Please choose rock, paper, or scissors."
	}
	return message, computerChoice
}

func handleRPSMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get player move from the request
	playerMove := r.FormValue("move")
	if playerMove == "" {
		http.Error(w, "Missing 'move' parameter", http.StatusBadRequest)
		return
	}

	// Game logic
	computerChoice := generateComputerChoice()
	message, _ := determineWinner(playerMove, computerChoice) // Ignore the returned computerChoice

	// Create game result response
	gameResult := GameResult{
		PlayerScore:    playerScore,
		ComputerScore:  computerScore,
		Message:        message,
		ComputerChoice: computerChoice, // Include the computer's choice here
	}

	// Encode as JSON and send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gameResult)
}

func main() {
	// Seed the random number generator once at the beginning
	rand.Seed(time.Now().UnixNano())

	http.Handle("/rps", enableCORS(http.HandlerFunc(handleRPSMove)))
	fmt.Println("Rock Paper Scissors server starting on port 8080")
	http.ListenAndServe(":8080", nil)
}
