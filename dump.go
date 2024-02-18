package main

import (
	"fmt"
	"strings"
)

const printableAliveCell cellCharacter = '@'
const printableDeadCell cellCharacter = '.'

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
