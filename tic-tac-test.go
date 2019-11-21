package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	s "strings"
	"text/template"
)

func boardCheck(board string, round int) {
	var numO, numX int
	numO = s.Count(board, "O")
	numX = s.Count(board, "X")

	if numO != (2 * round) {
		// TODO Create error for wrong number of O's
		fmt.Println("Incorrect number of O's")
	}
	if numX != (round) {
		// TODO Create error for wrong number of X's
		fmt.Println("Incorrect number of X's")
	}
}

func checkPlayerWin(board string) bool {
	// Initialize bool for player vistory state
	var playerVictory = false

	// Capture relevant board string
	//var subBoard = board[:9]

	// Check for player victory
	if board[0] == 'X' && board[1] == 'X' && board[2] == 'X' {
		playerVictory = true
	} else if board[3] == 'X' && board[4] == 'X' && board[5] == 'X' {
		playerVictory = true
	} else if board[6] == 'X' && board[7] == 'X' && board[8] == 'X' {
		playerVictory = true
	} else if board[0] == 'X' && board[3] == 'X' && board[6] == 'X' {
		playerVictory = true
	} else if board[1] == 'X' && board[4] == 'X' && board[7] == 'X' {
		playerVictory = true
	} else if board[2] == 'X' && board[5] == 'X' && board[8] == 'X' {
		playerVictory = true
	} else if board[0] == 'X' && board[4] == 'X' && board[8] == 'X' {
		playerVictory = true
	} else if board[2] == 'X' && board[4] == 'X' && board[6] == 'X' {
		playerVictory = true
	} else {
		// base case
		playerVictory = false
	}

	return playerVictory
}

func checkServerWin(board string) bool {
	// Initialize bool for player vistory state
	var serverVictory = false

	// Capture relevant board string
	//var subBoard = board[:9]

	// Check for server victory
	if board[0] == 'O' && board[1] == 'O' && board[2] == 'O' {
		serverVictory = true
	} else if board[3] == 'O' && board[4] == 'O' && board[5] == 'O' {
		serverVictory = true
	} else if board[6] == 'O' && board[7] == 'O' && board[8] == 'O' {
		serverVictory = true
	} else if board[0] == 'O' && board[3] == 'O' && board[6] == 'O' {
		serverVictory = true
	} else if board[1] == 'O' && board[4] == 'O' && board[7] == 'O' {
		serverVictory = true
	} else if board[2] == 'O' && board[5] == 'O' && board[8] == 'O' {
		serverVictory = true
	} else if board[0] == 'O' && board[4] == 'O' && board[8] == 'O' {
		serverVictory = true
	} else if board[2] == 'O' && board[4] == 'O' && board[6] == 'O' {
		serverVictory = true
	} else {
		// base case
		serverVictory = false
	}

	return serverVictory
}

func ticTacTest(board string, round int) (string, int, bool, bool) {
	// initialize variables
	var playerVictory bool
	var serverVictory bool

	// Check for valid board
	boardCheck(board, round)

	// Check for player victory
	playerVictory = checkPlayerWin(board)
	if playerVictory == true {
		// TODO player victory
		fmt.Println("Player victory: ", playerVictory)
	}

	// Place server moves
	board = s.Replace(board, "-", "O", 2)

	// Check for server victory
	serverVictory = checkServerWin(board)
	if serverVictory == true {
		// TODO server victory
		fmt.Println("Server victory: ", serverVictory)
	}

	round++
	return board, round, playerVictory, serverVictory
}

type GameBoard struct {
	Board         string
	Round         int
	XCount        int
	OCount        int
	playerVictory bool
	serverVictory bool
}

func (b *GameBoard) Check() bool {

}

func game(w http.ResponseWriter, r *http.Request) {
	b := GameBoard{
		Board: "---------",
	}
	t, _ := template.ParseFiles("game.gtpl")
	switch r.Method {
	case "GET":
	case "POST":
		r.ParseForm()
		b.Board = r.Form.Get("String")
		b.XCount = strings.Count(b.Board, "X")
		b.OCount = strings.Count(b.Board, "O")
		if b.Check() {

			return
		}
	default:
	}
	t.Execute(w, b)
}

func main() {
	b := "---------"
	r := 0
	reader := bufio.NewReader(os.Stdin)
	//var playerVictory, serverVictory bool

	http.HandleFunc("/game", game)

	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	var b2, r2, playerVictory, serverVictory = ticTacTest(b, r)
	// fmt.Println("Board: ", b2)

	for playerVictory == false && serverVictory == false {
		fmt.Println("Board: ", b2)
		fmt.Print("Enter the next board -> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		b2 = text
		b2, r2, playerVictory, serverVictory = ticTacTest(b2, r2)
	}
	// fmt.Println("Round: ", r2)
	// fmt.Println("Player victory: ", playerVictory)
	// fmt.Println("Server victory: ", serverVictory)

}
