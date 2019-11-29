package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	s "strings"
	"text/template"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomStringURLSafe returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomStringURLSafe(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}

func boardCheck(w http.ResponseWriter, r *http.Request, board *GameBoard) {
	board.XCount = s.Count(board.Board, "X")
	board.OCount = s.Count(board.Board, "O")

	// Check for expected number of O's
	if board.OCount != (2 * board.Round) {
		board.isCheating = true

		// Check for expected number of X's
	} else if board.XCount != (board.Round) {
		board.isCheating = true

		// Check for board size >= 9
	}
	for len(board.Board) < 9 {
		//board.isCheating = true
		board.Board += "-"
		// default case
	}
	// Check for any characters other than X, O, and -
	if restrictedCharacters(board.Board) {
		board.isCheating = true
	}

}

func restrictedCharacters(s string) bool {
	for _, r := range s {
		if r != 'X' && r != 'O' && r != '-' {
			return true
		}
	}
	return false
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
	m := make(map[string]string)

	m["1"], _ = GenerateRandomStringURLSafe(16)
	m["2"], _ = GenerateRandomStringURLSafe(16)
	m["3"], _ = GenerateRandomStringURLSafe(16)
	m["4"], _ = GenerateRandomStringURLSafe(16)
	m["5"], _ = GenerateRandomStringURLSafe(16)
	m["6"], _ = GenerateRandomStringURLSafe(16)
	m["7"], _ = GenerateRandomStringURLSafe(16)
	m["8"], _ = GenerateRandomStringURLSafe(16)
	m["9"], _ = GenerateRandomStringURLSafe(16)

	r1 := m["1"]
	r2 := m["2"]
	r3 := m["3"]
	r4 := m["4"]
	r5 := m["5"]
	r6 := m["6"]
	r7 := m["7"]
	r8 := m["8"]
	r9 := m["9"]

	switch r.Method {
	case "GET":
		http.SetCookie(w, &http.Cookie{Name: "Round", Value: r1})
		b.Board = "---------"
		b.Round = 0
		boardCheck(w, r, &b)
		checkPlayerWin(&b)
		// Place server moves
		b.Board = s.Replace(b.Board, "-", "O", 2)
		checkServerWin(&b)
	case "POST":
		// Populate board
		r.ParseForm()
		b.Board = r.Form.Get("String")
		// Get cookie value and set for next round
		cookie, _ := r.Cookie("Round")
		if cookie.Value == r1 {
			b.Round = 1
			cookie.Value = r2
			http.SetCookie(w, cookie)
		} else if cookie.Value == r2 {
			b.Round = 2
			cookie.Value = r3
			http.SetCookie(w, cookie)
		} else if cookie.Value == r3 {
			b.Round = 3
			cookie.Value = r4
			http.SetCookie(w, cookie)
		} else if cookie.Value == r4 {
			b.Round = 4
			cookie.Value = r5
			http.SetCookie(w, cookie)
		} else if cookie.Value == r5 {
			b.Round = 5
			cookie.Value = r6
			http.SetCookie(w, cookie)
		} else if cookie.Value == r6 {
			b.Round = 6
			cookie.Value = r7
			http.SetCookie(w, cookie)
		} else if cookie.Value == r7 {
			b.Round = 7
			cookie.Value = r8
			http.SetCookie(w, cookie)
		} else if cookie.Value == r8 {
			b.Round = 8
			cookie.Value = r9
			http.SetCookie(w, cookie)
		} else if cookie.Value == r9 {
			// Force error to cheaing screen
			b.Round = 100
			cookie.Value = r9
			http.SetCookie(w, cookie)
		}
		// Check for proper board
		boardCheck(w, r, &b)
		checkPlayerWin(&b)
		// Place server moves
		b.Board = s.Replace(b.Board, "-", "O", 2)
		checkServerWin(&b)
		switch {
		case b.isCheating:
			t, _ = template.ParseFiles("ischeating.gtpl")
		case b.playerVictory:
			t, _ = template.ParseFiles("playervictory.gtpl")
		case b.serverVictory:
			t, _ = template.ParseFiles("servervictory.gtpl")
		}
	default:
		t, _ = template.ParseFiles("ischeating.gtpl")
	}
	t.Execute(w, b)

	// Cookie Cleanup
	/*
		if t.Tree.Name != "game.gtpl" {
			cookie, _ := r.Cookie("Round")
			cookie.Value = ""
			http.SetCookie(w, cookie)
		}*/
}

func main() {
	http.HandleFunc("/game", game)

	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
