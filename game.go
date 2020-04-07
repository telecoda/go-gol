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
const bornFrames = 2
const dyingFrames = 2
const backgroundColour = termbox.ColorBlack
const titleColour = termbox.ColorYellow

// Text in the UI
const title = "GAME OF LIFE"

const animationSpeed = 200 * time.Millisecond

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

type Board struct {
	cells [][]Cell
}

type Game struct {
	board Board
}

// NewGame returns a fully-initialized game.
func NewGame() *Game {
	g := &Game{
		board: newBoard(),
	}
	return g
}

func newBoard() Board {
	b := Board{}
	b.cells = make([][]Cell, boardHeight)
	for y := 0; y < boardHeight; y++ {
		b.cells[y] = make([]Cell, boardWidth)
		for x := 0; x < boardWidth; x++ {
			b.cells[y][x] = Cell{}

			rnd := rand.Int() % 50
			if rnd < 6 {
				b.cells[y][x] = Cell{
					state:     Alive,
					animState: AliveAnim,
					animCount: 0,
				}
				births++
			} else {
				b.cells[y][x] = Cell{
					state:     Empty,
					animState: EmptyAnim,
					animCount: 0,
				}
			}
		}
	}

	return b
}

// This takes care of rendering everything.
func (g *Game) render() {
	termbox.Clear(backgroundColour, backgroundColour)
	tbprint(35, 0, titleColour, backgroundColour, title)
	tbprint(22, 1, titleColour, backgroundColour, "Red: birth  Green: Alive  Yellow: Dying")
	tbprint(22, 2, titleColour, backgroundColour, fmt.Sprintf("Generation: %d Births: %d Deaths: %d", generations, births, deaths))

	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			cell := g.board.cells[y][x]
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

	// copy board
	boardCopy := g.board.copy()

	g.board.updateCells(boardCopy)

	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			g.board.updateCellAnimation(x, y)
		}
	}
}

// updateCells updates the cells based on a copy of the board so it doesn change
func (b *Board) updateCells(boardCopy Board) {

	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			b.updateCellLives(boardCopy, x, y)
		}
	}

}

func (b *Board) updateCellLives(bc Board, x, y int) {
	// count live cells
	liveCellsAround := bc.getLiveCellCount(x, y)
	cell := b.getCell(x, y)

	// Any live cell with two or three neighbors survives.
	if cell.state == Alive && (liveCellsAround == 2 || liveCellsAround == 3) {
		// cell lives
		b.cells[y][x] = Cell{
			state:     Alive,
			animState: AliveAnim,
			animCount: 0,
		}
		return
	}

	// Any dead cell with three live neighbors becomes a live cell.
	if cell.state == Empty && liveCellsAround == 3 {
		// cell is born
		b.cells[y][x] = Cell{
			state:     Alive,
			animState: BornAnim,
			animCount: bornFrames,
		}
		births++
		return
	}

	// All other live cells die in the next generation.
	// Similarly, all other dead cells stay dead.
	if cell.state == Alive {
		// cell dies
		b.cells[y][x] = Cell{
			state:     Empty,
			animState: DyingAnim,
			animCount: dyingFrames,
		}
		deaths++
		return
	}
}

func (b *Board) updateCellAnimation(x, y int) {
	cell := b.getCell(x, y)
	if cell.animCount > 0 {
		cell.animCount--
		b.cells[y][x] = cell
	}

	if cell.animCount != 0 {
		return
	}

	if cell.animState == DyingAnim {
		b.cells[y][x] = Cell{
			state:     Empty,
			animState: EmptyAnim,
			animCount: 0,
		}
		return
	}
	if cell.animState == BornAnim {
		b.cells[y][x] = Cell{
			state:     Alive,
			animState: AliveAnim,
			animCount: 0,
		}
		return
	}
}

func (b *Board) copy() Board {
	newBoard := Board{
		cells: make([][]Cell, len(b.cells)),
	}
	// copy rows
	for r, row := range b.cells {
		newBoard.cells[r] = make([]Cell, len(row))
		for c, cell := range row {
			newBoard.cells[r][c] = cell
		}
	}
	return newBoard
}

func (b *Board) getCell(x, y int) Cell {
	if x < 0 || x > boardWidth-1 || y < 0 || y > boardHeight-1 {
		return Cell{state: Empty} // return empty cell if off board
	}

	return b.cells[y][x]
}

func (b *Board) getLiveCellCount(x, y int) int {
	// check all surrounding cells
	count := 0

	for cy := y - 1; cy <= (y + 1); cy++ {
		for cx := x - 1; cx <= (x + 1); cx++ {
			if cx == x && cy == y {
				// current cell, don't check
				continue
			}
			cell := b.getCell(cx, cy)
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
