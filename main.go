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
	nothing  controlAction = 0
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

	g.dumpTreeRecursive()
	fmt.Println()
	goLife(g)
}

func goLife(g *GOL) {
	actionStream := make(chan controlAction)
	proceed := make(chan bool)
	go readInput(actionStream, proceed)
	proceed <- true

	i := 0
	for action := range actionStream {
		switch action {
		case exit:
			close(actionStream)
			proceed <- false
			close(proceed)
		case next:
			i++
			fmt.Printf("\n%d\n", i)
			// g.nextGeneration()
			g.root = addBorder(g.root)
			g.gridSize = g.root.size

			g.root = evolve(addBorder(g.root))
			g.dumpTreeRecursive()
			proceed <- true
		case previous:
			fmt.Printf("\n prev is not supported using tree")
			// g.prevGeneration()
			// g.dump()
			proceed <- true
		case nothing:
			proceed <- true
		default:
			fmt.Printf("unsupported action: %v", action)
			proceed <- true
		}
	}
}

func readInput(action chan<- controlAction, proceed chan bool) {
	reader := bufio.NewReader(os.Stdin)
	for <-proceed {
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
			action <- nothing
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
		action <- nothing
	}

}
