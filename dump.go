package main

import (
	"fmt"
	"strings"
)

const printableAliveCell cellState = '@'
const printableDeadCell cellState = '.'

func (g *GOL) dump() {
	sb := strings.Builder{}

	border := make([]byte, g.gridSize+2, g.gridSize+2)
	i := 0
	border[i] = '|'
	for i++; i < g.gridSize+1; i++ {
		border[i] = '-'
	}
	border[i] = '|'

	sb.WriteByte('\n')
	sb.Write(border)

	for _, row := range g.grid {
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

func (g *GOL) dumpTree() {

	s := stack[*node]{
		data: make([]*node, 0, g.gridSize),
	}
	q := queue[*node]{
		data: make([]*node, 0, g.gridSize),
	}

	upDown := byte('U')

	leafsCnt := 0

	s.push(g.root)
	for !s.isEmpty() || !q.isEmpty() {
		if s.isEmpty() {
			for !q.isEmpty() {
				// the order after this operation is wrong, we are staring with last node of level 2
				s.push(q.pop())
			}
		}

		if leafsCnt == g.gridSize {
			if upDown == 'U' {
				upDown = 'D'
			} else {
				upDown = 'U'
			}
			leafsCnt = 0
			fmt.Printf("\n")
		}

		// n := s.top()
		n := s.pop()
		if n.level == 1 {
			leafsCnt++
			b := printableDeadCell
			if n.state == aliveCell {
				b = printableAliveCell
			}
			fmt.Print(string(b))
			// s.pop()
			continue
		}

		if upDown == 'U' {
			s.push(n.children.ne)
			s.push(n.children.nw)

			q.push(n)
			continue
		}

		// s.pop()

		s.push(n.children.se)
		s.push(n.children.sw)

	}

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
