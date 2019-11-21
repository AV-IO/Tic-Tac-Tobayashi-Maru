package main

import (
	"log"
	"net/http"
	s "strings"
	"text/template"
)

func boardCheck(w http.ResponseWriter, r *http.Request, board *GameBoard) {
	board.XCount = s.Count(board.Board, "X")
	board.OCount = s.Count(board.Board, "O")

	if board.OCount != (2 * board.Round) {
		board.isCheating = true
	}
	if board.XCount != (board.Round) {
		board.isCheating = true
	}
}

func checkPlayerWin(board *GameBoard) {
	// Check for player victory
	if board.Board[0] == 'X' && board.Board[1] == 'X' && board.Board[2] == 'X' {
		board.playerVictory = true
	} else if board.Board[3] == 'X' && board.Board[4] == 'X' && board.Board[5] == 'X' {
		board.playerVictory = true
	} else if board.Board[6] == 'X' && board.Board[7] == 'X' && board.Board[8] == 'X' {
		board.playerVictory = true
	} else if board.Board[0] == 'X' && board.Board[3] == 'X' && board.Board[6] == 'X' {
		board.playerVictory = true
	} else if board.Board[1] == 'X' && board.Board[4] == 'X' && board.Board[7] == 'X' {
		board.playerVictory = true
	} else if board.Board[2] == 'X' && board.Board[5] == 'X' && board.Board[8] == 'X' {
		board.playerVictory = true
	} else if board.Board[0] == 'X' && board.Board[4] == 'X' && board.Board[8] == 'X' {
		board.playerVictory = true
	} else if board.Board[2] == 'X' && board.Board[4] == 'X' && board.Board[6] == 'X' {
		board.playerVictory = true
	} else {
		// base case
		board.playerVictory = false
	}
}

func checkServerWin(board *GameBoard) {
	// Check for server victory
	if board.Board[0] == 'O' && board.Board[1] == 'O' && board.Board[2] == 'O' {
		board.serverVictory = true
	} else if board.Board[3] == 'O' && board.Board[4] == 'O' && board.Board[5] == 'O' {
		board.serverVictory = true
	} else if board.Board[6] == 'O' && board.Board[7] == 'O' && board.Board[8] == 'O' {
		board.serverVictory = true
	} else if board.Board[0] == 'O' && board.Board[3] == 'O' && board.Board[6] == 'O' {
		board.serverVictory = true
	} else if board.Board[1] == 'O' && board.Board[4] == 'O' && board.Board[7] == 'O' {
		board.serverVictory = true
	} else if board.Board[2] == 'O' && board.Board[5] == 'O' && board.Board[8] == 'O' {
		board.serverVictory = true
	} else if board.Board[0] == 'O' && board.Board[4] == 'O' && board.Board[8] == 'O' {
		board.serverVictory = true
	} else if board.Board[2] == 'O' && board.Board[4] == 'O' && board.Board[6] == 'O' {
		board.serverVictory = true
	} else {
		// base case
		board.serverVictory = false
	}
}

type GameBoard struct {
	Board         string
	Round         int
	XCount        int
	OCount        int
	playerVictory bool
	serverVictory bool
	isCheating    bool
}

func game(w http.ResponseWriter, r *http.Request) {
	b := GameBoard{
		Board:         "---------",
		Round:         0,
		XCount:        0,
		OCount:        0,
		playerVictory: false,
		serverVictory: false,
		isCheating:    false,
	}
	t, _ := template.ParseFiles("game.gtpl")

	switch r.Method {
	case "GET":
		http.SetCookie(w, &http.Cookie{Name: "Round", Value: "One"})
		b.Board = "---------"
		b.Round = 0
		boardCheck(w, r, &b)
		checkPlayerWin(&b)
		// Place server moves
		b.Board = s.Replace(b.Board, "-", "O", 2)
		checkServerWin(&b)
		t.Execute(w, b)
	case "POST":
		// Populate board
		r.ParseForm()
		b.Board = r.Form.Get("String")
		// Get cookie value and set for next round
		cookie, _ := r.Cookie("Round")
		if cookie.Value == "One" {
			b.Round = 1
			cookie.Value = "Two"
			http.SetCookie(w, cookie)
		} else if cookie.Value == "Two" {
			b.Round = 2
			cookie.Value = "Three"
			http.SetCookie(w, cookie)
		} else if cookie.Value == "Three" {
			b.Round = 3
			cookie.Value = "Four"
			http.SetCookie(w, cookie)
		}
		// Check for proper board
		boardCheck(w, r, &b)
		checkPlayerWin(&b)
		// Place server moves
		b.Board = s.Replace(b.Board, "-", "O", 2)
		checkServerWin(&b)
		// Check if the player is cheating
		switch b.isCheating {
		case true:
			t, _ = template.ParseFiles("ischeating.gtpl")
			fallthrough
		case false:
			fallthrough
		default:
			// do nothing
		}
		// Check for server victory
		switch b.serverVictory {
		case true:
			t, _ = template.ParseFiles("servervictory.gtpl")
			fallthrough
		case false:
			fallthrough
		default:
			// do nothing
		}
		// Check for player victory
		switch b.playerVictory {
		case true:
			t, _ = template.ParseFiles("playervictory.gtpl")
			fallthrough
		case false:
			fallthrough
		default:
			// do nothing
		}
		t.Execute(w, b)
	default:
	}
}

func main() {
	http.HandleFunc("/game", game)

	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
