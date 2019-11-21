package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	s "strings"
	"text/template"
)

func boardCheck(w http.ResponseWriter, r *http.Request, board *GameBoard) {
	board.XCount = s.Count(board.Board, "X")
	board.OCount = s.Count(board.Board, "O")

	if board.OCount != (2 * board.Round) {
		// TODO Create error for wrong number of O's
		board.isCheating = true
		fmt.Println("Incorrect number of O's")
		// TODO: Can I force the browser to a Cheating Page Here?
	}
	if board.XCount != (board.Round) {
		// TODO Create error for wrong number of X's
		board.isCheating = true
		fmt.Println("Incorrect number of X's")
		//r.RequestURI = "/xischeating"
		// TODO: Can I force the browser to a Cheating Page Here?
		//t, _ := template.ParseFiles("xischeating.gtpl")
		//t.Execute(w, board)
	}

	board.Round++
}

func checkPlayerWin(board *GameBoard) {

	// Capture relevant board string
	//var subBoard = board[:9]

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

	// Capture relevant board string
	//var subBoard = board[:9]

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

func ticTacTest(w http.ResponseWriter, r *http.Request, board *GameBoard) {

	// Check for valid board
	boardCheck(w, r, board)

	// Check for player victory
	checkPlayerWin(board)
	if board.playerVictory == true {
		// TODO player victory
		fmt.Println("Player victory: ", board.playerVictory)
	}

	// Place server moves
	board.Board = s.Replace(board.Board, "-", "O", 2)

	// Check for server victory
	checkServerWin(board)
	if board.serverVictory == true {
		// TODO server victory
		fmt.Println("Server victory: ", board.serverVictory)
	}

	//board.Round++
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
		ticTacTest(w, r, &b)
	case "POST":
		r.ParseForm()
		b.Board = r.Form.Get("String")
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
		b.XCount = strings.Count(b.Board, "X")
		b.OCount = strings.Count(b.Board, "O")
		// Check if cheating is detected
		if isCheating(b) {
			// TODO Change template to static html
			t, _ := template.ParseFiles("xischeating.gtpl")
			t.Execute(w, b)
			return
		}
		fmt.Println("Beginning Round: ", b.Round)
		fmt.Println("Cookie value: ", cookie.Value)
		ticTacTest(w, r, &b)
	default:
	}

	t.Execute(w, b)
}

func main() {
	//b := "---------"
	//r := 0
	//reader := bufio.NewReader(os.Stdin)
	//var playerVictory, serverVictory bool

	http.HandleFunc("/game", game)

	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	//ticTacTest(board)
	// fmt.Println("Board: ", b2)

	// for playerVictory == false && serverVictory == false {
	// 	fmt.Println("Board: ", b2)
	// 	fmt.Print("Enter the next board -> ")
	// 	text, _ := reader.ReadString('\n')
	// 	// convert CRLF to LF
	// 	text = strings.Replace(text, "\n", "", -1)
	// 	b2 = text
	// 	b2, r2, playerVictory, serverVictory = ticTacTest(b2, r2)
	// }
	// fmt.Println("Round: ", r2)
	// fmt.Println("Player victory: ", playerVictory)
	// fmt.Println("Server victory: ", serverVictory)

}
