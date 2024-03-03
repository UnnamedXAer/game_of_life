package main

import (
	"math"
)

type nodeChildren struct {
	nw, ne, sw, se *node
}

type node struct {
	children   nodeChildren // children nodes, all nil if it is a leaf node
	level      int          // 1 for leaf, 2 for node consisting of four leaves, >= 3 for others
	state      cellState    // only matters if a leaf node, if not leaf it is 1 if at least one of it's leaves is alive
	population int          // numbers of alive leaves in this square
	size       int          // size of the square, that this node represents
}

// could the leaves be just zeros or ones value of the pointer?
var (
	aliveLeaf = &node{state: aliveCell, level: 1, population: 1, size: 1}
	deadLeaf  = &node{state: deadCell, level: 1, population: 0, size: 1}
)

func newNode(children nodeChildren, level int, size int) *node {
	// each "pattern" will be represented by exactly one node, pattern can be composed of a single leaf or a tree ending with leaves
	n, ok := cacheNodes[children]
	if ok {
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

var cacheNodes map[nodeChildren]*node = make(map[nodeChildren]*node, 100)

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

// buildLeafs returns node children base on the grid and given left top position of the node,
// the children are leaf nodes either alive or dead
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
