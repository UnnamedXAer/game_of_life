package main

import (
	"fmt"
	"strings"
)

const printableAliveCell cellCharacter = '+'
const printableDeadCell cellCharacter = ' '

func (g GOL) dump() {
	sb := strings.Builder{}

	bh := make([]byte, g.grigSize+2, g.grigSize+2)
	i := 0
	bh[i] = '|'
	for i++; i < g.grigSize+1; i++ {
		bh[i] = '-'
	}
	bh[i] = '|'

	sb.WriteByte('\n')
	sb.Write(bh)

	for _, row := range g.grid {
		sb.WriteByte('\n')
		sb.WriteByte('|')
		for _, cell := range row {
			if cell == 0 {
				sb.WriteByte(printableDeadCell)
			} else {
				sb.WriteByte(printableAliveCell)
			}

		}
		sb.WriteByte('|')
	}
	sb.WriteByte('\n')
	sb.Write(bh)
	sb.WriteByte('\n')
	fmt.Print(sb.String())
	// return sb.String()
}
