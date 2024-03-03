package main

import "fmt"

type cellState = byte

const aliveCell cellState = 1
const deadCell cellState = 0
const fastforward = false

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
	for i, v := range g.grid { // is this necessary?
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

// tree base implementation of conway's game of life

type evolveResult = *node

var cacheEvolveResults map[*node]evolveResult = make(map[*node]evolveResult, 100)

// evolve performs evolution of the node (or at least the center of it)
// it always return the center of the given node (the center is always half of
// the size of th node).
func evolve(n *node) evolveResult {

	// we are reusing nodes, so if we once calculate the result of a node with given children
	// then we can reuse it next time we see node with the same shape (children), no matter
	// what level that node is.
	r, ok := cacheEvolveResults[n]
	if ok {
		return r
	}

	if n.level == 3 {
		r = evolveGol(n)
		cacheEvolveResults[n] = r
		return r
	}

	a1 := newNode(nodeChildren{
		n.children.nw.children.nw,
		n.children.nw.children.ne,
		n.children.nw.children.sw,
		n.children.nw.children.se,
	}, n.level-1, n.size/2)

	a2 := newNode(nodeChildren{
		n.children.nw.children.ne,
		n.children.ne.children.nw,
		n.children.nw.children.se,
		n.children.ne.children.sw,
	}, n.level-1, n.size/2)

	a3 := newNode(nodeChildren{
		n.children.ne.children.nw,
		n.children.ne.children.ne,
		n.children.ne.children.sw,
		n.children.ne.children.se,
	}, n.level-1, n.size/2)

	a4 := newNode(nodeChildren{
		n.children.nw.children.sw,
		n.children.nw.children.se,
		n.children.sw.children.nw,
		n.children.sw.children.ne,
	}, n.level-1, n.size/2)

	a5 := newNode(nodeChildren{
		n.children.nw.children.se,
		n.children.ne.children.sw,
		n.children.sw.children.ne,
		n.children.se.children.nw,
	}, n.level-1, n.size/2)

	a6 := newNode(nodeChildren{
		n.children.ne.children.sw,
		n.children.ne.children.se,
		n.children.se.children.nw,
		n.children.se.children.ne,
	}, n.level-1, n.size/2)

	a7 := newNode(nodeChildren{
		n.children.sw.children.nw,
		n.children.sw.children.ne,
		n.children.sw.children.sw,
		n.children.sw.children.se,
	}, n.level-1, n.size/2)

	a8 := newNode(nodeChildren{
		n.children.sw.children.ne,
		n.children.se.children.nw,
		n.children.sw.children.se,
		n.children.se.children.sw,
	}, n.level-1, n.size/2)

	a9 := newNode(nodeChildren{
		n.children.se.children.nw,
		n.children.se.children.ne,
		n.children.se.children.sw,
		n.children.se.children.se,
	}, n.level-1, n.size/2)

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

	var center *node
	// when fastforward is on, we jump a few generation each time we call evolve
	if fastforward {
		center = assembleCenterNode(
			evolve(res1),
			evolve(res2),
			evolve(res3),
			evolve(res4),
		)
	} else {
		center = assembleCenterNode(
			getCenterNode(res1),
			getCenterNode(res2),
			getCenterNode(res3),
			getCenterNode(res4),
		)
	}

	cacheEvolveResults[n] = center

	return center
}

// evolveGol implements the real evolutions of nodes(cells)
// it converts node to matrix, performs the "evolution" base on rules,
// then converts the matrix back to node. The resulting node is a center (i.e. 4x4 -> 2x2)
// of the given node after evolution,
func evolveGol(n *node) evolveResult {
	m := convertNodeToAuxMatrix(n)
	mNext := nextGeneration(m)
	nextGenNode := convertAuxMatrixToNode(mNext)
	return nextGenNode
}

func getCanonical(state cellState) *node {
	if state == deadCell {
		return deadLeaf
	}
	return aliveLeaf
}

func convertNodeToAuxMatrix(n *node) [6][6]cellState {

	if n.level != 3 {
		panic(fmt.Sprintf("\ncannot create 6x6 matrix from node at lvl: %d\n node: %+v", n.level, n))
	}

	m := [6][6]cellState{}

	m[1][1] = n.children.nw.children.nw.state
	m[1][2] = n.children.nw.children.ne.state
	m[1][3] = n.children.ne.children.nw.state
	m[1][4] = n.children.ne.children.ne.state

	m[2][1] = n.children.nw.children.sw.state
	m[2][2] = n.children.nw.children.se.state
	m[2][3] = n.children.ne.children.sw.state
	m[2][4] = n.children.ne.children.se.state

	m[3][1] = n.children.sw.children.nw.state
	m[3][2] = n.children.sw.children.ne.state
	m[3][3] = n.children.se.children.nw.state
	m[3][4] = n.children.se.children.ne.state

	m[4][1] = n.children.sw.children.sw.state
	m[4][2] = n.children.sw.children.se.state
	m[4][3] = n.children.se.children.sw.state
	m[4][4] = n.children.se.children.se.state

	return m
}

func convertAuxMatrixToNode(m [6][6]cellState) *node {
	children := nodeChildren{
		newNode(nodeChildren{
			getCanonical(m[1][1]),
			getCanonical(m[1][2]),
			getCanonical(m[2][1]),
			getCanonical(m[2][2]),
		},
			2,
			2,
		),

		newNode(nodeChildren{
			getCanonical(m[1][3]),
			getCanonical(m[1][4]),
			getCanonical(m[2][3]),
			getCanonical(m[2][4]),
		},
			2,
			2,
		),

		newNode(nodeChildren{
			getCanonical(m[3][1]),
			getCanonical(m[3][2]),
			getCanonical(m[4][1]),
			getCanonical(m[4][2]),
		},
			2,
			2,
		),

		newNode(nodeChildren{
			getCanonical(m[3][3]),
			getCanonical(m[3][4]),
			getCanonical(m[4][3]),
			getCanonical(m[4][4]),
		},
			2,
			2,
		),
	}

	n := newNode(children, 3, 4)
	n = getCenterNode(n) // could we just compose the center ourself?

	return n
}

// generateCanonical0 creates a tree of nodes (which leaves are dead) where root will
// be of the given level
func generateCanonical0(level int) *node {
	if level == 0 {
		return nil
	}

	n := deadLeaf
	for i := 2; i <= level; i++ {
		n = newNode(nodeChildren{n, n, n, n}, n.level+1, n.size*2)
	}

	return n
}

// addBorder surrounds the given node with border (of nodes with dead leaves) which level
// will be one higher and the given node will end up as the returned node's center.
func addBorder(n *node) *node {
	level := n.level

	nodeBorder := generateCanonical0(level - 1)

	nw := newNode(nodeChildren{
		nodeBorder, nodeBorder, nodeBorder, n.children.nw,
	}, level, n.size)
	ne := newNode(nodeChildren{
		nodeBorder, nodeBorder, n.children.ne, nodeBorder,
	}, level, n.size)
	sw := newNode(nodeChildren{
		nodeBorder, n.children.sw, nodeBorder, nodeBorder,
	}, level, n.size)
	se := newNode(nodeChildren{
		n.children.se, nodeBorder, nodeBorder, nodeBorder,
	}, level, n.size)

	return newNode(nodeChildren{nw, ne, sw, se}, level+1, n.size*2)
}

func assembleCenterNode(nw, ne, sw, se evolveResult) evolveResult {
	children := nodeChildren{nw, ne, sw, se}

	n := newNode(children, nw.level+1, 2*nw.size)

	return n
}

func getCenterNode(n *node) *node {
	children := nodeChildren{
		n.children.nw.children.se,
		n.children.ne.children.sw,
		n.children.sw.children.ne,
		n.children.se.children.nw,
	}

	return newNode(children, n.level-1, n.size/2)
}

func nextGeneration(m [6][6]cellState) [6][6]cellState {
	const gridSize = 6
	grid := [gridSize][gridSize]cellState{}

	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			c := m[y][x]
			aliveNeighbours := countAliveNeighbours(m, y, x)
			grid[y][x] = getNextGenerationState(aliveNeighbours, c)
		}
	}

	return grid
}

// countAliveNeighbours is a helpers function that counts alive cells in an aux grid
// created from node at level 3 (of size 4x4) with padding of dead cells to simplify
// operations
func countAliveNeighbours(m [6][6]cellState, y int, x int) int {
	const gridSize = 6

	yyStart := max(0, y-1)
	yyEnd := min(gridSize, y+2)
	xxStart := max(0, x-1)
	xxEnd := min(gridSize, x+2)

	aliveNeighbours := 0

	for yy := yyStart; yy < yyEnd; yy++ {
		for xx := xxStart; xx < xxEnd; xx++ {
			if y == yy && x == xx {
				continue
			}
			if m[yy][xx] == aliveCell {
				aliveNeighbours++
			}
		}
	}

	return aliveNeighbours
}

// getNextGenerationState returns state of a cell in the next generation.
// those are the basic rules for Conway's game of life.
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
