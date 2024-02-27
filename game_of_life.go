package main

import "fmt"

type cellState = byte

const aliveCell cellState = 1
const deadCell cellState = 0

type point struct {
	y, x int
}

type golGrid = [][]byte

type GOL struct {
	gridSize int
	history  []golGrid
	grid     golGrid
	root     *node
}

func newGOL(n int) *GOL {
	g := GOL{
		gridSize: n,
		grid:     make(golGrid, n),
		history:  make([]golGrid, 0, 10),
	}

	for i := 0; i < n; i++ {
		g.grid[i] = make([]byte, n)
	}

	return &g
}

func (g *GOL) nextGeneration() {

	grid := make(golGrid, g.gridSize)
	for i, v := range g.grid {
		grid[i] = make([]byte, g.gridSize)
		copy(grid[i], v)
	}

	for y := 0; y < g.gridSize; y++ {
		for x := 0; x < g.gridSize; x++ {
			c := g.grid[y][x]
			aliveNeighbours := g.countAliveNeighbours(y, x)
			grid[y][x] = getNextGenerationState(aliveNeighbours, c)
		}
	}

	g.history = append(g.history, g.grid)

	g.grid = grid
}

func getNextGenerationState(aliveNeighbours int, c cellState) cellState {

	if c == aliveCell {

		if aliveNeighbours < 2 {
			return deadCell // dies
		}

		if aliveNeighbours <= 3 {
			return aliveCell // live to the next generation
		}

		return deadCell // dies by overpopulation
	} else if c == deadCell {
		if aliveNeighbours == 3 {
			return aliveCell // becomes alive by reproduction
		}

		return c
	} else {
		panic(fmt.Sprintf("\nbad cell state: %q", string(c)))
	}
}

func (g *GOL) countAliveNeighbours(y int, x int) int {

	yyStart := max(0, y-1)
	yyEnd := min(g.gridSize, y+2)

	xxStart := max(0, x-1)
	xxEnd := min(g.gridSize, x+2)

	aliveNeighbours := 0

	for yy := yyStart; yy < yyEnd; yy++ {
		for xx := xxStart; xx < xxEnd; xx++ {
			if y == yy && x == xx {
				continue
			}
			if g.grid[yy][xx] == aliveCell {
				aliveNeighbours++
			}
		}
	}

	return aliveNeighbours
}

func (g *GOL) prevGeneration() {
	if len(g.history) == 0 {
		fmt.Printf("\nno history")
		return
	}
	g.grid = g.history[len(g.history)-1]
	g.history[len(g.history)-1] = nil
	g.history = g.history[:len(g.history)-1]
}

////////////////////// tree base

type evolveResult = *node

func evolve(n *node) evolveResult {

	if n.level == 4 {
		return evolveGol(n)
	}

	a1 := newNode(nodeChildren{
		n.children.nw.children.nw.children.se,
		n.children.nw.children.ne.children.sw,
		n.children.nw.children.sw.children.ne,
		n.children.nw.children.se.children.nw,
	}, n.level-1, n.size/2, "a1")

	a2 := newNode(nodeChildren{
		n.children.nw.children.ne.children.se,
		n.children.ne.children.nw.children.sw,
		n.children.nw.children.se.children.ne,
		n.children.ne.children.sw.children.nw,
	}, n.level-1, n.size/2, "a2")

	a3 := newNode(nodeChildren{
		n.children.ne.children.nw.children.se,
		n.children.ne.children.ne.children.sw,
		n.children.ne.children.sw.children.ne,
		n.children.ne.children.se.children.nw,
	}, n.level-1, n.size/2, "a3")

	a4 := newNode(nodeChildren{
		n.children.nw.children.sw.children.se,
		n.children.nw.children.se.children.sw,
		n.children.sw.children.nw.children.ne,
		n.children.sw.children.ne.children.nw,
	}, n.level-1, n.size/2, "a4")

	a5 := newNode(nodeChildren{
		n.children.nw.children.se.children.se,
		n.children.ne.children.sw.children.sw,
		n.children.sw.children.ne.children.ne,
		n.children.se.children.nw.children.nw,
	}, n.level-1, n.size/2, "a5")

	a6 := newNode(nodeChildren{
		n.children.ne.children.sw.children.se,
		n.children.ne.children.se.children.sw,
		n.children.se.children.nw.children.ne,
		n.children.se.children.ne.children.nw,
	}, n.level-1, n.size/2, "a6")

	a7 := newNode(nodeChildren{
		n.children.sw.children.nw.children.se,
		n.children.sw.children.ne.children.sw,
		n.children.sw.children.sw.children.ne,
		n.children.sw.children.se.children.nw,
	}, n.level-1, n.size/2, "a7")

	a8 := newNode(nodeChildren{
		n.children.sw.children.ne.children.se,
		n.children.se.children.nw.children.sw,
		n.children.sw.children.se.children.ne,
		n.children.se.children.sw.children.nw,
	}, n.level-1, n.size/2, "a8")

	a9 := newNode(nodeChildren{
		n.children.se.children.nw.children.se,
		n.children.se.children.ne.children.sw,
		n.children.se.children.sw.children.ne,
		n.children.se.children.se.children.nw,
	}, n.level-1, n.size/2, "a9")

	r1 := evolve(a1)
	r2 := evolve(a2)
	r3 := evolve(a3)
	r4 := evolve(a4)
	r5 := evolve(a5)
	r6 := evolve(a6)
	r7 := evolve(a7)
	r8 := evolve(a8)
	r9 := evolve(a9)

	res1 := assembleCenterNode(r1, r2, r4, r5)
	res2 := assembleCenterNode(r2, r3, r5, r6)
	res3 := assembleCenterNode(r4, r5, r7, r8)
	res4 := assembleCenterNode(r5, r6, r8, r9)

	center := assembleCenterNode(
		getCenterNode(res1),
		getCenterNode(res2),
		getCenterNode(res3),
		getCenterNode(res4),
	)
	// center := assembleCenterNode(
	// 	evolve(res1),
	// 	evolve(res2),
	// 	evolve(res3),
	// 	evolve(res4),
	// )

	return center
}

var stateToCell = map[cellState]*node{deadCell: deadLeaf, aliveCell: aliveLeaf}

func evolveGol(n *node) evolveResult {

	var nw, ne, sw, se *node

	aliveNeighbours := 0
	for _, neighbour := range []*node{
		n.children.nw.children.nw,
		n.children.nw.children.ne,
		n.children.ne.children.nw,

		n.children.nw.children.sw,
		n.children.ne.children.sw,

		n.children.sw.children.nw,
		n.children.sw.children.ne,
		n.children.se.children.nw,
	} {
		if neighbour.state == aliveCell {
			aliveNeighbours++
		}
	}
	nw = stateToCell[getNextGenerationState(aliveNeighbours, n.children.nw.children.se.state)]

	aliveNeighbours = 0
	for _, neighbour := range []*node{
		n.children.nw.children.ne,
		n.children.ne.children.nw,
		n.children.ne.children.ne,

		n.children.nw.children.se,
		n.children.ne.children.se,

		n.children.sw.children.ne,
		n.children.se.children.nw,
		n.children.se.children.ne,
	} {
		if neighbour.state == aliveCell {
			aliveNeighbours++
		}
	}
	ne = stateToCell[getNextGenerationState(aliveNeighbours, n.children.ne.children.sw.state)]

	aliveNeighbours = 0
	for _, neighbour := range []*node{
		n.children.nw.children.sw,
		n.children.nw.children.se,
		n.children.ne.children.sw,

		n.children.sw.children.nw,
		n.children.ne.children.nw,

		n.children.sw.children.sw,
		n.children.sw.children.se,
		n.children.se.children.sw,
	} {
		if neighbour.state == aliveCell {
			aliveNeighbours++
		}
	}
	sw = stateToCell[getNextGenerationState(aliveNeighbours, n.children.sw.children.ne.state)]

	aliveNeighbours = 0
	for _, neighbour := range []*node{
		n.children.nw.children.se,
		n.children.ne.children.sw,
		n.children.ne.children.se,

		n.children.sw.children.ne,
		n.children.se.children.ne,

		n.children.sw.children.se,
		n.children.se.children.sw,
		n.children.se.children.se,
	} {
		if neighbour.state == aliveCell {
			aliveNeighbours++
		}
	}
	se = stateToCell[getNextGenerationState(aliveNeighbours, n.children.se.children.nw.state)]

	children := nodeChildren{nw, ne, sw, se}

	return newNode(children, n.level-1, n.size/2, "center from leaves of: "+n.label)
}

func generateCanonical0(level int) *node {
	if level == 0 {
		return nil
	}

	n := deadLeaf
	for i := 2; i <= level; i++ {
		n = newNode(nodeChildren{n, n, n, n}, n.level+1, n.size*2, fmt.Sprintf("canonical0 at depth: %d", n.level+1))
	}

	return n
}

func addBorder(n *node) *node {
	level := n.level

	nodeBorder := generateCanonical0(level - 1)

	nw := newNode(nodeChildren{
		nodeBorder, nodeBorder, nodeBorder, n.children.nw,
	}, level, n.size, fmt.Sprintf("nw - border node of level %d", level))
	ne := newNode(nodeChildren{
		nodeBorder, nodeBorder, n.children.ne, nodeBorder,
	}, level, n.size, fmt.Sprintf("ne - border node of level %d", level))
	sw := newNode(nodeChildren{
		nodeBorder, n.children.sw, nodeBorder, nodeBorder,
	}, level, n.size, fmt.Sprintf("sw - border node of level %d", level))
	se := newNode(nodeChildren{
		n.children.se, nodeBorder, nodeBorder, nodeBorder,
	}, level, n.size, fmt.Sprintf("se - border node of level %d", level))

	return newNode(nodeChildren{nw, ne, sw, se}, level+1, n.size*2, fmt.Sprintf("bordered node at level: %d", level+1))
}

func assembleCenterNode(nw, ne, sw, se evolveResult) evolveResult {
	children := nodeChildren{nw, ne, sw, se}

	n := newNode(children, nw.level+1, 2*nw.size, "center")

	return n

}

func getCenterNode(n *node) *node {
	children := nodeChildren{
		n.children.nw.children.se,
		n.children.ne.children.sw,
		n.children.sw.children.ne,
		n.children.se.children.nw,
	}

	return newNode(children, n.level-1, n.size/2, "center of "+n.label)
}
