package main

import (
	"fmt"
	"strings"
)

const printableAliveCell cellState = '@'
const printableDeadCell cellState = '.'

func dump(grid [][]cellState) {
	sb := strings.Builder{}

	gridSize := len(grid)
	border := make([]byte, gridSize+2)
	i := 0
	border[i] = '|'
	for i++; i < gridSize+1; i++ {
		border[i] = '-'
	}
	border[i] = '|'

	sb.WriteByte('\n')
	sb.Write(border)

	for _, row := range grid {
		sb.WriteByte('\n')
		sb.WriteByte('|')
		for _, cell := range row {
			if cell == aliveCell {
				sb.WriteByte(printableAliveCell)
			} else {
				sb.WriteByte(printableDeadCell)
			}

		}
		sb.WriteByte('|')
	}
	sb.WriteByte('\n')
	sb.Write(border)
	sb.WriteByte('\n')
	fmt.Print(sb.String())
}

func (g *GOL) dump() {
	dump(g.grid)
}

func dumpTreeRecHelper(n *node, grid [][]byte, y, x int) {
	if n.level == 1 {
		if n.state == aliveCell {
			grid[y][x] = printableAliveCell
		} else {
			grid[y][x] = printableDeadCell
		}
		return
	}

	dumpTreeRecHelper(n.children.nw, grid, y, x)
	dumpTreeRecHelper(n.children.ne, grid, y, x+n.children.ne.size)
	dumpTreeRecHelper(n.children.sw, grid, y+n.children.sw.size, x)
	dumpTreeRecHelper(n.children.se, grid, y+n.children.se.size, x+n.children.se.size)
}

func (g *GOL) dumpTreeRecursive() {
	dumpTreeRecursive(g.root)
}

func dumpTreeRecursive(n *node) {
	gridSize := n.size
	printableGrid := make([][]byte, gridSize)
	for i := 0; i < gridSize; i++ {
		printableGrid[i] = make([]byte, gridSize)
	}
	dumpTreeRecHelper(n, printableGrid, 0, 0)

	fmt.Println("|" + strings.Repeat("-", gridSize) + "|")
	for _, line := range printableGrid {
		fmt.Println("|" + string(line) + "|")
	}
	fmt.Println("|" + strings.Repeat("-", gridSize) + "|")
}

type stack[T any] struct {
	data []T
}

func (s *stack[T]) isEmpty() bool {
	return len(s.data) == 0
}

func (s *stack[T]) push(item T) {
	s.data = append(s.data, item)
}

func (s *stack[T]) pop() T {
	item := s.data[len(s.data)-1]
	s.data[len(s.data)-1] = *new(T)
	s.data = s.data[:len(s.data)-1]
	return item
}

func (s *stack[T]) top() T {
	return s.data[len(s.data)-1]
}

type queue[T any] struct {
	data []T
}

func (q *queue[T]) isEmpty() bool {
	return len(q.data) == 0
}

func (q *queue[T]) push(item T) {
	q.data = append(q.data, item)
}

func (q *queue[T]) pop() T {
	item := q.data[0]
	q.data[0] = *new(T)
	q.data = q.data[1:]
	return item
}

func (q *queue[T]) top() T {
	return q.data[0]
}
