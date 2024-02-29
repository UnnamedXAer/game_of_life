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

func nextGeneration(m [6][6]cellState) [6][6]cellState {

	const gridSize = 6
	grid := [gridSize][gridSize]cellState{}

	// for i, v := range m {
	// 	copy(grid[i][:], v[:])
	// }

	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			c := m[y][x]
			aliveNeighbours := countAliveNeighbours(m, y, x)
			grid[y][x] = getNextGenerationState(aliveNeighbours, c)
		}
	}

	return grid
}

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

	if n.level == 3 {
		return evolveGol(n)
	}

	a1 := newNode(nodeChildren{
		n.children.nw.children.nw,
		n.children.nw.children.ne,
		n.children.nw.children.sw,
		n.children.nw.children.se,
	}, n.level-1, n.size/2, "a1")

	a2 := newNode(nodeChildren{
		n.children.nw.children.ne,
		n.children.ne.children.nw,
		n.children.nw.children.se,
		n.children.ne.children.sw,
	}, n.level-1, n.size/2, "a2")

	a3 := newNode(nodeChildren{
		n.children.ne.children.nw,
		n.children.ne.children.ne,
		n.children.ne.children.sw,
		n.children.ne.children.se,
	}, n.level-1, n.size/2, "a3")

	a4 := newNode(nodeChildren{
		n.children.nw.children.sw,
		n.children.nw.children.se,
		n.children.sw.children.nw,
		n.children.sw.children.ne,
	}, n.level-1, n.size/2, "a4")

	a5 := newNode(nodeChildren{
		n.children.nw.children.se,
		n.children.ne.children.sw,
		n.children.sw.children.ne,
		n.children.se.children.nw,
	}, n.level-1, n.size/2, "a5")

	a6 := newNode(nodeChildren{
		n.children.ne.children.sw,
		n.children.ne.children.se,
		n.children.se.children.nw,
		n.children.se.children.ne,
	}, n.level-1, n.size/2, "a6")

	a7 := newNode(nodeChildren{
		n.children.sw.children.nw,
		n.children.sw.children.ne,
		n.children.sw.children.sw,
		n.children.sw.children.se,
	}, n.level-1, n.size/2, "a7")

	a8 := newNode(nodeChildren{
		n.children.sw.children.ne,
		n.children.se.children.nw,
		n.children.sw.children.se,
		n.children.se.children.sw,
	}, n.level-1, n.size/2, "a8")

	a9 := newNode(nodeChildren{
		n.children.se.children.nw,
		n.children.se.children.ne,
		n.children.se.children.sw,
		n.children.se.children.se,
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

func getCanonical(state cellState) *node {
	if state == deadCell {
		return deadLeaf
	}
	return aliveLeaf
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
			"from matrix - nw",
		),

		newNode(nodeChildren{
			getCanonical(m[1][3]),
			getCanonical(m[1][4]),
			getCanonical(m[2][3]),
			getCanonical(m[2][4]),
		},
			2,
			2,
			"from matrix - ne",
		),

		newNode(nodeChildren{
			getCanonical(m[3][1]),
			getCanonical(m[3][2]),
			getCanonical(m[4][1]),
			getCanonical(m[4][2]),
		},
			2,
			2,
			"from matrix - sw",
		),

		newNode(nodeChildren{
			getCanonical(m[3][3]),
			getCanonical(m[3][4]),
			getCanonical(m[4][3]),
			getCanonical(m[4][4]),
		},
			2,
			2,
			"from matrix - se",
		),
	}

	n := newNode(children, 3, 4, "from matrix")
	n = getCenterNode(n)

	return n
}

func evolveGol(n *node) evolveResult {
	m := convertNodeToAuxMatrix(n)
	mNext := nextGeneration(m)
	nextGenNode := convertAuxMatrixToNode(mNext)
	return nextGenNode
	// return nextGenNode

	fmt.Printf("\nnode: %+v", nextGenNode)

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
