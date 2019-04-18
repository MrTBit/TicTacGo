package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"math"
	"os"
)

type Position struct {
	side int
	x int
	y int
	played int
}


func main() {

	encoding.Register()

	s, e := tcell.NewScreen()
	if e != nil{
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil{
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.Clear()

	spaces := drawGrid(s, defStyle)

	s.Show()


	player1 := true
	win := false

	for {
		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
		case *tcell.EventMouse:
			mx, my := ev.Position()

			if ev.Buttons() == tcell.Button1 {
				clicked := checkValidClick(spaces, mx, my)
				if clicked > -1 {
					win = move(s, spaces, clicked, player1)

					if !win {
						if player1 {
							player1 = false
						} else {
							player1 = true
						}
					} else {
						//win condition
						fmt.Println("win")
					}
				}
			}
		}
		s.Show()
	}

}


//calculates and draws grid, returning space locations
func drawGrid(s tcell.Screen, style tcell.Style) []Position {
	w, h := s.Size()

	//leave 1/4 of screen area for info
	//grid will be square
	start := w / 4
	gridSide := 0
	if h < (w - start){
		gridSide = h
	} else {
		gridSide = w - start
	}

	gridStart := w/3
	gridSpacing := int(math.Round(float64(gridSide) / 3.0))
	vLinePos := [2]int{gridStart + gridSpacing, gridStart + gridSpacing * 2}
	hLinePos := [2]int{gridSpacing, gridSpacing * 2}

	for col := gridStart + 1; col <= gridSide + gridStart; col++ {
		for row := 0; row <= gridSide; row++ {
			if (col == vLinePos[0] || col == vLinePos[1]) && (row == hLinePos[0] || row == hLinePos[1]) { //intersection
				s.SetContent(col, row, tcell.RunePlus, nil, style)
			} else if col == vLinePos[0] || col == vLinePos[1] {
				s.SetContent(col, row, tcell.RuneVLine, nil, style)
			} else if row == hLinePos[0] || row == hLinePos[1] {
				s.SetContent(col, row, tcell.RuneHLine, nil, style)
			}  else {
				s.SetCell(col, row, style, ' ')
			}
		}
	}


	spaces := make([]Position, 9)
	spaces[0] = Position{x:gridStart, y:0, side:gridSpacing-2, played:0}
	spaces[1] = Position{x:gridStart + gridSpacing + 1, y:0, side: gridSpacing - 2, played:0}
	spaces[2] = Position{x: gridStart + gridSpacing * 2 + 1, y:0, side: gridSpacing - 2,played:0}
	spaces[3] = Position{x: gridStart, y:gridSpacing + 1, side: gridSpacing - 2, played:0}
	spaces[4] = Position{x: gridStart + gridSpacing + 1, y:gridSpacing + 1, side: gridSpacing - 2, played:0}
	spaces[5] = Position{x: gridStart + gridSpacing*2 + 1, y: gridSpacing + 1, side: gridSpacing - 2, played:0}
	spaces[6] = Position{x: gridStart, y: gridSpacing * 2 + 1, side: gridSpacing - 2, played:0}
	spaces[7] = Position{x: gridStart + gridSpacing + 1,y: gridSpacing * 2 + 1, side: gridSpacing - 2, played:0}
	spaces[8] = Position{x: gridStart + gridSpacing*2 + 1, y: gridSpacing * 2 + 1, side: gridSpacing - 2, played:0}

	return spaces


}

func move(s tcell.Screen, spaces []Position, clickedSpace int ,player1 bool)(win bool){
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)

	if player1 {
		spaces[clickedSpace].played = 1
		style = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorDarkBlue)
	} else {
		spaces[clickedSpace].played = 4
		style = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
	}

	for col := spaces[clickedSpace].x; col <= spaces[clickedSpace].x + spaces[clickedSpace].side; col++{
		for row := spaces[clickedSpace].y; row <= spaces[clickedSpace].y + spaces[clickedSpace].side; row++{
			if spaces[clickedSpace].played == 1 {
				if (col == spaces[clickedSpace].x || col == spaces[clickedSpace].x+spaces[clickedSpace].side) && (row == spaces[clickedSpace].y || row == spaces[clickedSpace].y+spaces[clickedSpace].side) { //corners
					s.SetContent(col, row, tcell.RunePlus, nil, style)
				} else if col == spaces[clickedSpace].x || col == spaces[clickedSpace].x+spaces[clickedSpace].side {
					r, _, _, _ := s.GetContent(col, row)

					if r == tcell.RunePlus {
						continue
					} else {
						s.SetContent(col, row, tcell.RuneVLine, nil, style)
					}
				} else if row == spaces[clickedSpace].y || row == spaces[clickedSpace].y+spaces[clickedSpace].side {
					r, _, _, _ := s.GetContent(col, row)

					if r == tcell.RunePlus {
						continue
					} else {
						s.SetContent(col, row, tcell.RuneHLine, nil, style)
					}
				}
			} else if spaces[clickedSpace].played == 4 {
				s.SetCell(col, row, style, 'X')
			}
		}
	}

	return checkWin(spaces)
}

func checkWin(spaces []Position)(win bool){

	//test vertical wins
	for col := 0; col < 3; col++{
		test := 0
		for row := 0; row < 3; row++{
			if spaces[col + row*3].played != 0 {
				test += spaces[col + row*3].played
			}
		}
		if test == 3 || test == 12{
			return true
		}
	}

	//test horizontal wins
	for row := 0; row < 3; row++{
		test := 0
		for col := 0; col < 3; col++{
			if spaces[col + row*3].played != 0 {
				test += spaces[col + row*3].played
			}
		}
		if test == 3 || test == 12 {
			return true
		}
	}

	//test diagonals
	if(spaces[0].played == 1 && spaces[4].played == 1 && spaces[8].played == 1) || (spaces[2].played == 1 && spaces[4].played == 1 && spaces[6].played == 1){
		return true
	} else if (spaces[0].played == 4 && spaces[4].played == 4 && spaces[8].played == 4) || (spaces[2].played == 4 && spaces[4].played == 4 && spaces[6].played == 4){
		return true
	}

	return false
}

func checkValidClick(spaces []Position, x, y int)(clickedSpace int){

	for i, space := range spaces {
		if (x > space.x && x < space.x + space.side) && (y > space.y && y < space.y + space.side){
			if space.played == 0 {
				return i
			}
		}
	}
	return -1
}