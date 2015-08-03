package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const boardWidth = 40
const boardHeight = 40
const numStates = 2
const cellWidth = 2
const bornFrames = 5
const dyingFrames = 5
const backgroundColour = termbox.ColorBlack
const titleColour = termbox.ColorYellow

// Text in the UI
const title = "GAME OF LIFE"

const animationSpeed = 1 * time.Millisecond

const defaultMarginWidth = 2
const defaultMarginHeight = 1
const boardStartX = defaultMarginWidth
const boardStartY = defaultMarginHeight + 2
const boardEndX = boardStartX + boardWidth*cellWidth
const boardEndY = boardStartY + boardHeight

type CellState int

const (
	Empty CellState = 0
	Alive CellState = 1
)

type AnimState int

const (
	EmptyAnim AnimState = 0
	AliveAnim AnimState = 1
	BornAnim  AnimState = 2
	DyingAnim AnimState = 3
)

var lifeColours = []termbox.Attribute{
	termbox.ColorBlue,   // Empty
	termbox.ColorGreen,  // Alive
	termbox.ColorRed,    // Born
	termbox.ColorYellow, // Dying
}

var (
	generations = 0
	births      = 0
	deaths      = 0
)

type Cell struct {
	state     CellState
	animState AnimState
	animCount int
}

type Game struct {
	board [][]Cell
}

// NewGame returns a fully-initialized game.
func NewGame() *Game {
	g := new(Game)
	g.reset()
	return g
}

// Reset the game in order to play again.
func (g *Game) reset() {
	g.board = make([][]Cell, boardHeight)
	for y := 0; y < boardHeight; y++ {
		g.board[y] = make([]Cell, boardWidth)
		for x := 0; x < boardWidth; x++ {

			rnd := rand.Int() % 50
			if rnd < 6 {
				g.board[y][x] = Cell{
					state:     Alive,
					animState: AliveAnim,
					animCount: 0,
				}
				births++
			} else {
				g.board[y][x] = Cell{
					state:     Empty,
					animState: EmptyAnim,
					animCount: 0,
				}

			}
		}
	}

}

// This takes care of rendering everything.
func (g *Game) render() {
	termbox.Clear(backgroundColour, backgroundColour)
	tbprint(35, 0, titleColour, backgroundColour, title)
	tbprint(22, 1, titleColour, backgroundColour, "Red: birth  Green: Alive  Yellow: Dying")
	tbprint(22, 2, titleColour, backgroundColour, fmt.Sprintf("Generation: %d Births: %d Deaths: %d", generations, births, deaths))

	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			cell := g.board[y][x]
			absCellValue := int(cell.animState)
			cellColor := lifeColours[absCellValue]
			for i := 0; i < cellWidth; i++ {
				termbox.SetCell(boardStartX+cellWidth*x+i, boardStartY+y, ' ', cellColor, cellColor)
			}
		}
	}

	termbox.Flush()
}

// update game - apply rules of game of life to board
func (g *Game) update() {

	generations++

	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {

			g.updateCellAnimation(x, y)
			g.updateCellLives(x, y)

		}
	}

}

func (g *Game) updateCellLives(x, y int) {
	// count live cells
	liveCellsAround := g.getLiveCellCount(x, y)
	cell := g.getCell(x, y)

	if cell.state == Empty && liveCellsAround == 3 {
		// cell is born
		g.board[y][x] = Cell{
			state:     Alive,
			animState: BornAnim,
			animCount: bornFrames,
		}
		births++
		return
	}
	if cell.state == Alive && liveCellsAround < 2 {
		// cell dies
		g.board[y][x] = Cell{
			state:     Empty,
			animState: DyingAnim,
			animCount: dyingFrames,
		}
		deaths++
		return
	}
	if cell.state == Alive && (liveCellsAround == 2 || liveCellsAround == 3) {
		// cell lives
		g.board[y][x] = Cell{
			state:     Alive,
			animState: AliveAnim,
			animCount: 0,
		}
		return
	}
	if cell.state == Alive && liveCellsAround > 3 {
		// cell dies - overcrowding
		g.board[y][x] = Cell{
			state:     Empty,
			animState: DyingAnim,
			animCount: dyingFrames,
		}
		deaths++
		return
	}
}

func (g *Game) updateCellAnimation(x, y int) {
	cell := g.getCell(x, y)
	if cell.animCount > 0 {
		cell.animCount--
		g.board[y][x] = cell
	}

	if cell.animCount != 0 {
		return
	}

	if cell.animState == DyingAnim {
		g.board[y][x] = Cell{
			state:     Empty,
			animState: EmptyAnim,
			animCount: 0,
		}
		return
	}
	if cell.animState == BornAnim {
		g.board[y][x] = Cell{
			state:     Alive,
			animState: AliveAnim,
			animCount: 0,
		}
		return
	}
}

func (g *Game) getCell(x, y int) Cell {
	if x < 0 || x > boardWidth-1 || y < 0 || y > boardHeight-1 {
		return Cell{state: Empty} // return empty cell if off board
	}

	return g.board[y][x]
}

func (g *Game) getLiveCellCount(x, y int) int {
	// check all surrounding cells
	count := 0

	for cy := y - 1; cy <= (y + 1); cy++ {
		for cx := x - 1; cx <= (x + 1); cx++ {
			if cx == x && cy == y {
				// current cell, don't check
				continue
			}
			cell := g.getCell(cx, cy)
			if cell.state == Alive {
				count++
			}
		}
	}

	return count
}

// Function tbprint draws a string.
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}
