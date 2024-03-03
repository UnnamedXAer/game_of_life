package main

import (
	"math"
)

type nodeChildren struct {
	nw, ne, sw, se *node
}

type node struct {
	children   nodeChildren // children nodes, nil if leaf node
	level      int          // 1 for leaf, 2 for node consisting of four leave, >= 3 for others
	state      cellState    // only matters if a leaf node
	population int          // numbers of active nodes expected to be 0..4
	size       int          // size of the square, that this node represents
	// next       *node        // nil if not computed yet
}

var (
	aliveLeaf = &node{state: aliveCell, level: 1, population: 1, size: 1}
	deadLeaf  = &node{state: deadCell, level: 1, population: 0, size: 1}
)

func newNode(children nodeChildren, level int, size int) *node {
	n, ok := cacheNodes[children]
	if ok {
		// fmt.Printf("\ngot hit: %+v", n)
		return n
	}
	population := children.nw.population + children.ne.population + children.sw.population + children.se.population

	state := deadCell
	if population > 0 {
		state = aliveCell
	}

	n = &node{
		children:   children,
		level:      level,
		state:      state,
		population: population,
		size:       size,
	}

	cacheNodes[children] = n

	return n
}

const (
	nwPos = iota
	nePos
	swPos
	sePos
)

func (g *GOL) buildTree() {
	f := math.Log2(float64(g.gridSize))
	g.root = g.buildNode(int(f)+1, g.gridSize, 0, 0)
}

var cacheNodes map[nodeChildren]*node = make(map[nodeChildren]*node, 0)

func (g *GOL) buildNode(level int, size int, y, x int) *node {

	children := nodeChildren{}

	if level == 2 {
		return newNode(g.buildLeafs(y, x), level, size)
	}

	halfSize := size / 2
	i := 0
	for yy := y; yy < y+2; yy++ {
		for xx := x; xx < x+2; xx++ {

			switch i {
			case nwPos:
				children.nw = g.buildNode(level-1, halfSize, y, x)
			case nePos:
				children.ne = g.buildNode(level-1, halfSize, y, x+halfSize)
			case swPos:
				children.sw = g.buildNode(level-1, halfSize, y+halfSize, x)
			case sePos:
				children.se = g.buildNode(level-1, halfSize, y+halfSize, x+halfSize)
			}

			i++
		}
	}

	return newNode(children, level, size)
}

// buildLeafs returns children for nodes, the children are leaf nodes either alive or dead
func (g *GOL) buildLeafs(y, x int) nodeChildren {

	children := nodeChildren{}
	var tmp *node
	i := 0
	for yy := y; yy < y+2; yy++ {
		for xx := x; xx < x+2; xx++ {

			if g.grid[yy][xx] == aliveCell {
				tmp = aliveLeaf
			} else {
				tmp = deadLeaf
			}

			switch i {
			case nwPos:
				children.nw = tmp
			case nePos:
				children.ne = tmp
			case swPos:
				children.sw = tmp
			case sePos:
				children.se = tmp
			}

			i++
		}
	}

	return children
}
