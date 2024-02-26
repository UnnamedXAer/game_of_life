package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

	fmt.Println()

	g.dumpTreeRecursive()
	n := generateCanonical0(g.root.level)
	g.root = n
	fmt.Println()
	fmt.Println()
	g.dumpTreeRecursive()
	fmt.Println()
	fmt.Println()
	g.root = addBorder(g.root) // we got root node with level 6 but we are missing level 5.
	g.dumpTreeRecursive()
	fmt.Println()

	// evolve(g.root)
	// return
	// g.dump()

	// g.dumpTreeRecursive()
	// goLife(g)
}

func goLife(g *GOL) {
	actionStream := make(chan controlAction)
	go readInput(actionStream)

	i := 0
	for action := range actionStream {
		switch action {
		case exit:
			close(actionStream)
		case next:
			i++
			fmt.Printf("\n%d\n", i)
			// g.nextGeneration()
			g.root = evolve(g.root)
			g.dumpTreeRecursive()
		case previous:
			fmt.Printf("\n prev is not supported using tree")
			// g.prevGeneration()
			// g.dump()

		default:
		}

	}
}

func readInput(action chan<- controlAction) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\nwaiting for key: (wsad): \n")

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
