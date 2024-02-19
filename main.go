package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

type controlAction byte

const (
	next     controlAction = 'n'
	previous controlAction = 'p'
	exit     controlAction = 'e'
)

func main() {
	fmt.Printf("\nmain")

	const gridSize int = 64
	g := newGOL(gridSize)

	if strings.Contains(os.Args[1], "patterns") {
		setPatternsFromFile(os.Args[1], g)
	} else {
		setGridFromFile(os.Args[1], g)
	}
	g.dump()

	f := math.Log2(float64(g.gridSize))
	root := g.buildNode(int(f)+1, g.gridSize, 0, 0)

	// TODO: display grid from the tree

	fmt.Printf("\nroot: %v", root)

	// return

	goLife(g)
}

func goLife(g *GOL) {
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
