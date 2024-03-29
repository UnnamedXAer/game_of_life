package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	// "runtime/pprof"
)

type controlAction byte

const (
	next     controlAction = 'n'
	previous controlAction = 'p'
	exit     controlAction = 'e'
	nothing  controlAction = 0
)

func main() {
	// debug.SetGCPercent(-1)
	// debug.SetMemoryLimit(math.MaxInt64)

	// // Create a memory profile file
	// memProfileFile, err := os.Create("mem.prof")
	// if err != nil {
	// 	panic(err)
	// }
	// defer memProfileFile.Close()

	fmt.Printf("\nmain\n")

	const gridSize int = 64
	g := newGOL(gridSize)

	if strings.Contains(os.Args[1], "patterns") {
		setPatternsFromFile(os.Args[1], g)
	} else {
		setGridFromFile(os.Args[1], g)
	}

	// add a little empty space around
	// g.root = addBorder(g.root)
	g.root = addBorder(g.root)
	g.gridSize = g.root.size

	g.dumpTreeRecursive()
	fmt.Println()
	// goLife(g)
	goLifeUserTriggered(g)

	// // Write memory profile to file
	// if err := pprof.WriteHeapProfile(memProfileFile); err != nil {
	// 	panic(err)
	// }
	fmt.Println("end of main")
}

func goLife(g *GOL) {
	for i := 0; i < 100; i++ {
		fmt.Println(i)

		g.root = evolve(addBorder(g.root))
		g.dumpTreeRecursive()
		fmt.Println()
		fmt.Println("done: ", i)
	}
	g.dumpTreeRecursive()

	fmt.Println("\nout of loop, returning from goLife")
}

func goLifeUserTriggered(g *GOL) {
	actionStream := make(chan controlAction)
	proceed := make(chan bool)
	go readInput(actionStream, proceed)
	proceed <- true

	i := 0
	for action := range actionStream {
		// if i == 5 {
		// 	close(actionStream)
		// 	close(proceed)
		// 	break
		// }
		switch action {
		case exit:
			close(actionStream)
			proceed <- false
			close(proceed)
		case next:
			i++
			fmt.Printf("\n%d\n", i)
			// g.nextGeneration()

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
			continue
		}

		input := b[0]

		switch b[0] {
		case '\n', '\r':
			action <- next
			continue
		case 'd', 'w':
			action <- next
			continue
		case 'a', 's':
			action <- previous
			continue

		default:
			fmt.Printf("\n u pressed something that doesn't make sense :). key: %v ,%q", input, string(input))
			action <- nothing
		}
	}
}
