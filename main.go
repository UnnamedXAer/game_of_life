package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Printf("\nmain")

	const gridSize int = 16
	g := newGOL(gridSize)

	setPatternsFromFile(os.Args[1], g)

	i := 0
	for {
		i++
		fmt.Printf("\n%d", i)
		g.dump()
		g.nextGeneration()
		time.Sleep(time.Second * 3)
	}

}

type cellCharacter = byte

const aliveCell cellCharacter = '+'
const deadCell cellCharacter = 0

type point struct {
	y, x int
}

type GOL struct {
	grigSize int
	grid     [][]byte
}

func newGOL(n int) GOL {
	g := GOL{
		grigSize: n,
		grid:     make([][]byte, n, n),
	}

	for i := 0; i < n; i++ {
		g.grid[i] = make([]byte, n, n)
	}

	return g
}

func (g GOL) nextGeneration() {

	grid := make([][]byte, g.grigSize, g.grigSize)
	copy(grid, g.grid)

	for y := 0; y < g.grigSize; y++ {
		for x := 0; x < g.grigSize; x++ {
			c := g.grid[y][x]
			aliveNeighbours := 0

			yyStart := max(0, y-1)
			yyEnd := min(g.grigSize, y+2)
			for yy := yyStart; yy < yyEnd; yy++ {
				xxStart := max(0, x-1)
				xxEnd := min(g.grigSize, x+2)
				for xx := xxStart; xx < xxEnd; xx++ {
					if y == yy && x == xx {
						continue
					}
					if g.grid[yy][xx] == aliveCell {
						aliveNeighbours++
					}
				}
			}

			if c == aliveCell {

				if aliveNeighbours < 2 {
					// dies
					grid[y][x] = deadCell
					continue
				}

				if aliveNeighbours <= 2 {
					// live to the next generation
					continue
				}

				grid[y][x] = deadCell // dies by overpopulation
			} else if c == deadCell {
				if aliveNeighbours == 3 {
					grid[y][x] = aliveCell // becomes alive by reproduction
				}

			} else {
				panic(fmt.Sprintf("\nbad cell state: %q, at: %3d, %3d", string(c), y, x))
			}

		}
	}

	g.grid = grid
}
