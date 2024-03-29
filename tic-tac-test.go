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
func GenerateRandomStringURLSafe(n int) string {
	b, _ := GenerateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b)
}

func boardCheck(board *GameBoard) {
	board.XCount = s.Count(board.Board, "X")
	board.OCount = s.Count(board.Board, "O")

	if board.Round > 9 ||
		restrictedCharacters(board.Board) || // Check for any characters other than X, O, and -
		board.OCount != (2*board.Round) || // Check for expected number of O's
		board.XCount != (board.Round) { // Check for expected number of X's
		board.IsCheating = true
	}
	// Check for board size > 9
	for len(board.Board) <= 9 {
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
	fm := template.FuncMap{
		"combine": func(i int, s string) string {
			return string(rune(i)) + s
		},
		"ReplaceIndexWX": func(s string) string {
			return s[1:int(s[0])+1] + "X" + s[int(s[0])+2:]
		},
		"idx": func(s string) string {
			return string(s[int(s[0])+1])
		},
	}
	t, err := template.New("game.htpl").Funcs(fm).ParseFiles("game.htpl")
	if err != nil {
		panic(err)
	}

	NewBoard := func(cVal string, round int) {
		c := &http.Cookie{Name: "SESSION", Value: cVal}
		http.SetCookie(w, c)
		b.Round = round
		cmap.Store(c.Value, b.Round)
		// Place server moves
		b.Board = s.Replace("---------", "-", "O", 2)
	}

	switch r.Method {
	case "GET":
		c, err := r.Cookie("SESSION")
		if err == http.ErrNoCookie {
			NewBoard(GenerateRandomStringURLSafe(32), 0)
		} else {
			cmap.Store(c.Value, 0)
		}
	case "POST":
		b.IsPlaying = true
		// Populate board
		err = r.ParseForm()
		if err != nil {
			panic(err)
		}
		c, err := r.Cookie("SESSION")
		if err == http.ErrNoCookie {
			NewBoard(GenerateRandomStringURLSafe(32), 0)
			break
		} else {
			roundi, found := cmap.Load(c.Value)
			if !found {
				NewBoard(c.Value, 0)
				break
			}
			switch roundi.(int) {
			case 9:
				b.Round = 100
			default:
				b.Round = roundi.(int) + 1
			}
			cmap.Store(c.Value, b.Round)
		}
		b.Board = r.Form.Get("Board")
		if b.Board == "" {
			NewBoard(c.Value, 0)
			break
		}
		// Get cookie value and increment internal value
		// Check for proper board
		boardCheck(&b)
		checkWin(&b, 'X')
		// Place server moves
		b.Board = s.Replace(b.Board, "-", "O", 2)
		checkWin(&b, 'O')
		if b.IsCheating || b.PlayerVictory || b.ServerVictory {
			cmap.Store(c.Value, 0)
		}
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
