package main

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"net/http"
	s "strings"
	"sync"
)

type GameBoard struct {
	Board         string
	Round         int
	XCount        int
	OCount        int
	PlayerVictory bool
	ServerVictory bool
	IsCheating    bool
	IsPlaying     bool
}

var cmap *sync.Map

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

	if board.Round > 9 ||
		restrictedCharacters(board.Board) || // Check for any characters other than X, O, and -
		board.OCount != (2*board.Round) || // Check for expected number of O's
		board.XCount != (board.Round) { // Check for expected number of X's
		board.IsCheating = true
	}
	// Check for board size >= 9
	for len(board.Board) < 9 {
		board.Board += "-"
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

func checkWin(board *GameBoard, r byte) {
	var victory *bool
	if r == 'X' {
		victory = &board.PlayerVictory
	} else if r == 'O' {
		victory = &board.ServerVictory
	}
	// Check for all win conditions
	if (board.Board[0] == r && board.Board[1] == r && board.Board[2] == r) ||
		(board.Board[3] == r && board.Board[4] == r && board.Board[5] == r) ||
		(board.Board[6] == r && board.Board[7] == r && board.Board[8] == r) ||
		(board.Board[0] == r && board.Board[3] == r && board.Board[6] == r) ||
		(board.Board[1] == r && board.Board[4] == r && board.Board[7] == r) ||
		(board.Board[2] == r && board.Board[5] == r && board.Board[8] == r) ||
		(board.Board[0] == r && board.Board[4] == r && board.Board[8] == r) ||
		(board.Board[2] == r && board.Board[4] == r && board.Board[6] == r) {
		*victory = true
	} else {
		*victory = false
	}
}

func game(w http.ResponseWriter, r *http.Request) {
	b := GameBoard{
		Board: "---------",
	}
	t, err := template.ParseFiles("game.htpl")
	if err != nil {
		panic(err)
	}

	NewBoard := func() {
		v, _ := GenerateRandomStringURLSafe(32)
		c := &http.Cookie{Name: "SESSION", Value: v}
		http.SetCookie(w, c)
		b.Round = 0
		cmap.Store(c.Value, b.Round)
		// Place server moves
		b.Board = s.Replace("---------", "-", "O", 2)
	}

	switch r.Method {
	case "GET":
		c, err := r.Cookie("SESSION")
		if err == http.ErrNoCookie {
			NewBoard()
		} else {
			cmap.Store(c.Name, 0)
		}
	case "POST":
		b.IsPlaying = true
		// Populate board
		err = r.ParseForm()
		if err != nil {
			panic(err)
		}
		b.Board = r.Form.Get("String")
		if b.Board == "" {
			NewBoard()
			break
		}
		// Get cookie value and increment internal value
		c, err := r.Cookie("SESSION")
		if err == http.ErrNoCookie {
			NewBoard()
			break
		} else {
			roundi, found := cmap.Load(c.Name)
			if !found {
				NewBoard()
				break
			}
			switch roundi.(int) {
			case 9:
				b.Round = 100
			default:
				b.Round = roundi.(int) + 1
			}
			cmap.Store(c.Name, b.Round)
		}
		// Check for proper board
		boardCheck(w, r, &b)
		checkWin(&b, 'X')
		// Place server moves
		b.Board = s.Replace(b.Board, "-", "O", 2)
		checkWin(&b, 'O')
	default:
		b.IsCheating = true
	}
	err = t.Execute(w, b)
	if err != nil {
		panic(err)
	}
}

func main() {
	cmap = new(sync.Map)
	// Static file handling
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	// Template handling
	http.HandleFunc("/", game)

	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		panic(err)
	}
}
