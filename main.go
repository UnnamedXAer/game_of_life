package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

type controlAction byte

const (
	next     controlAction = 'n'
	previous controlAction = 'p'
	exit     controlAction = 'e'
)

func main() {
	fmt.Printf("\nmain")

	const gridSize int = 16
	g := newGOL(gridSize)

	setPatternsFromFile(os.Args[1], g)
	g.dump()

	actionStream := make(chan controlAction)
	go readInput(actionStream)

	i := 0
	for action := range actionStream {
		switch action {
		case exit:
			close(actionStream)
			break
		case next:
			i++
			fmt.Printf("\n%d", i)
			g.nextGeneration()
			g.dump()
		case previous:
			g.prevGeneration()
			g.dump()
			// fmt.Printf("\nprevious not implemented yet :(")

		default:
		}

	}
}

func readInput(action chan<- controlAction) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\nwaiting for key: (wsad): ")

		b, err := reader.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Printf("\nExiting... %v ", err)
				action <- exit
				return
			}
			fmt.Printf("\n read key: err: %v ", err)
			continue
		}
		if len(b) == 0 {
			fmt.Printf("\n read nothing")
			continue
		}

		input := b[0]

		if input == '\n' || input == '\r' {
			return
		}

		if input == 'd' || input == 'w' {
			action <- next
			continue
		}

		if input == 'a' || input == 's' {
			action <- previous
			continue
		}

		fmt.Printf("\n u pressed something that doesn't make sense :). key: %v ,%q", input, string(input))
	}

}

type cellCharacter = byte

const aliveCell cellCharacter = '+'
const deadCell cellCharacter = 0

type point struct {
	y, x int
}

type golGrid = [][]byte

type GOL struct {
	gridSize int
	history  []golGrid
	grid     golGrid
}

func newGOL(n int) GOL {
	g := GOL{
		gridSize: n,
		grid:     make(golGrid, n, n),
		history:  make([]golGrid, 0, 10),
	}

	for i := 0; i < n; i++ {
		g.grid[i] = make([]byte, n, n)
	}

	return g
}

func (g *GOL) nextGeneration() {

	grid := make(golGrid, g.gridSize, g.gridSize)
	for i, v := range g.grid {
		grid[i] = make([]byte, g.gridSize)
		copy(grid[i], v)
	}

	for y := 0; y < g.gridSize; y++ {
		for x := 0; x < g.gridSize; x++ {
			c := g.grid[y][x]
			aliveNeighbours := g.countAliveNeighbours(y, x)
			grid[y][x] = getNextGenerationState(aliveNeighbours, c)
		}
	}

	g.history = append(g.history, g.grid)

	g.grid = grid
}

func getNextGenerationState(aliveNeighbours int, c cellCharacter) cellCharacter {

	if c == aliveCell {

		if aliveNeighbours < 2 {
			return deadCell // dies
		}

		if aliveNeighbours == 2 {
			return aliveCell // live to the next generation
		}

		return deadCell // dies by overpopulation
	} else if c == deadCell {
		if aliveNeighbours == 3 {
			return aliveCell // becomes alive by reproduction
		}

		return c
	} else {
		panic(fmt.Sprintf("\nbad cell state: %q", string(c)))
	}
}

func (g *GOL) countAliveNeighbours(y int, x int) int {

	yyStart := max(0, y-1)
	yyEnd := min(g.gridSize, y+2)

	xxStart := max(0, x-1)
	xxEnd := min(g.gridSize, x+2)

	aliveNeighbours := 0

	for yy := yyStart; yy < yyEnd; yy++ {
		for xx := xxStart; xx < xxEnd; xx++ {
			if y == yy && x == xx {
				continue
			}
			if g.grid[yy][xx] == aliveCell {
				aliveNeighbours++
			}
		}
	}

	return aliveNeighbours
}

func (g *GOL) prevGeneration() {
	if len(g.history) == 0 {
		fmt.Printf("\nno history")
		return
	}
	g.grid = g.history[len(g.history)-1]
	g.history[len(g.history)-1] = nil
	g.history = g.history[:len(g.history)-1]
}
