package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

func setPatternsFromFile(fn string, g *GOL) {
	fnClean := filepath.Clean(fn)
	b, err := os.ReadFile(fnClean)
	if err != nil {
		panic(fmt.Errorf("input file: %w", err))
	}

	patterData := bytes.Split(b, []byte{'\n'})

	for i := 0; i < len(patterData); i++ {
		line := patterData[i]
		lineSize := len(line)
		if lineSize == 0 {
			continue
		}
		if line[0] == '!' {
			// just a comment
			continue
		}

		var p point
		if line[0] == '>' {
			p = extractStartPosition(lineSize, line, i)
			i++
		} else {
			panic(fmt.Sprintf("\nunexpected character %q in line: %d, ot pos 0", string(line[0]), i))
		}

		// extract actual pattern
		for ; i < len(patterData); i++ {
			lineSize = len(patterData[i])
			if lineSize != 0 {
				break
			}
		}

		patternStartLine := i
		for ; i < len(patterData); i++ {
			lineSize = len(patterData[i])
			if lineSize == 0 {
				break
			}
		}

		if patternStartLine == i {
			panic(fmt.Sprintf("\nempty pattern, line: %d", i))
		}

		pattern := make([][]byte, i-patternStartLine, i-patternStartLine)
		for w, pline := range patterData[patternStartLine:i] {
			pattern[w] = make([]byte, len(pline), len(pline))
			for u, cell := range pline {
				if cell == printableAliveCell {
					pattern[w][u] = aliveCell
				} else if cell == printableDeadCell {
					pattern[w][u] = deadCell
				} else {
					panic(fmt.Sprintf("\nunexpected cell character: %q in pattern. row: %d, cell: %d", string(cell), w, u))
				}
			}
		}
		// copy(pattern, patterData[patternStartLine:i])

		setPattern(pattern, p, g)
	}
}

// coordinates where the pattern should be placed on the grid (its left top corner)
// , expected syntax: `> 43 44`
func extractStartPosition(lineSize int, line []byte, i int) point {

	y := 0
	x := 0

	if lineSize < 3 {
		panic(fmt.Sprintf("\ninvalid line: %d", i))
	}

	k := 1
	for ; k < lineSize; k++ {
		if line[k] != ' ' {
			break
		}
	}

	// extract y
	size := 0
	for ; k < lineSize; k++ {
		if line[k] == ' ' {
			break
		}
		if line[k] >= '0' && line[k] <= '9' {
			size++
		} else {
			panic(fmt.Sprintf("\nunexpected character %q in line: %d, at pos %d", string(line[0]), i, k))
		}
	}

	if size == 0 {
		panic(fmt.Sprintf("\nmissing y coordinates in line: %d", i))
	}
	multiplier := 1
	for j := k - 1; j >= k-size; j-- {
		y += int(line[j]-'0') * multiplier
		multiplier *= 10
	}

	// extract x
	size = 0
	for ; k < lineSize; k++ {
		if size == 0 && line[k] == ' ' {
			continue
		}

		if line[k] >= '0' && line[k] <= '9' {
			size++
		} else {
			panic(fmt.Sprintf("\nunexpected character %q in line: %d, at pos %d", string(line[0]), i, k))
		}
	}

	if size == 0 {
		panic(fmt.Sprintf("\nmissing x coordinates in line: %d", i))
	}
	multiplier = 1
	for j := k - 1; j >= k-size; j-- {
		x += int(line[j]-'0') * multiplier
		multiplier *= 10
	}

	return point{y, x}
}

func setPattern(pattern [][]byte, startPos point, g *GOL) {

	y := startPos.y
	x := startPos.x
	for i, pline := range pattern {
		copy(g.grid[y+i][x:x+len(pline)], pline)
	}
}

func setGridFromFile(fn string, g *GOL) {

	fnClean := filepath.Clean(fn)
	b, err := os.ReadFile(fnClean)
	if err != nil {
		panic(fmt.Errorf("input file: %w", err))
	}

	grid := golGrid(bytes.Split(b, []byte{'\n'}))
	grid = make(golGrid, 0, 8)

	grid = append(grid, make([]byte, 0, 8))
	i := 0
	for _, v := range b {
		if v == '\n' {
			grid = append(grid, make([]byte, 0, len(grid[i])))
			i++
			continue
		}
		grid[i] = append(grid[i], v)
	}

	gridSize := len(grid) // max from no of rows or the longest row
	for y, row := range grid {
		for x, v := range row {
			if v == printableAliveCell {
				grid[y][x] = aliveCell
				continue
			}

			if v == printableDeadCell {
				grid[y][x] = deadCell
				continue
			}
			panic(fmt.Sprintf("\nunknown character: %q, at: %d x %d", string(v), y, x))
		}
		if len(row) > gridSize {
			gridSize = len(row)
		}
	}

	expectedGridSize := 1
	for {
		if expectedGridSize >= gridSize {
			break
		}

		expectedGridSize *= 2
	}
	gridSize = expectedGridSize

	if gridSize > len(grid) {
		oldGrid := grid
		grid = make(golGrid, gridSize, gridSize)
		copy(grid, oldGrid)
	}

	for y := 0; y < gridSize; y++ {
		if len(grid[y]) < gridSize {
			row := grid[y]
			grid[y] = make([]byte, gridSize, gridSize)
			copy(grid[y], row)
		}
	}

	g.grid = grid
	g.gridSize = gridSize
}
