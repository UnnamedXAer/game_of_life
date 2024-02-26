package main

import (
	"fmt"
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
	label string // temporary for debugging info
}

type nodeMap map[nodeChildren]*node

var (
	aliveLeaf = &node{state: aliveCell, level: 1, population: 1, size: 1, label: "aliveLeaf"}
	deadLeaf  = &node{state: deadCell, level: 1, population: 0, size: 1, label: "deadLeaf"}
)

func newNode(children nodeChildren, level int, size int, label string) *node {
	population := children.nw.population + children.ne.population + children.sw.population + children.se.population

	state := deadCell
	if population > 0 {
		state = aliveCell
	}

	return &node{
		children:   children,
		level:      level,
		state:      state,
		population: population,
		size:       size,
		label:      label,
	}
}

const (
	nwPos = iota
	nePos
	swPos
	sePos
)

func (g *GOL) buildTree() {
	f := math.Log2(float64(g.gridSize))
	g.root = g.buildNode(int(f)+1, g.gridSize, 0, 0, "root")

	// i := 0
	// for children, cnt := range cache {
	// 	i++
	// 	fmt.Printf("\n%3d. %+v: %4d", i, children, cnt)
	// }
	// fmt.Println()
}

var cache map[nodeChildren]int = make(map[nodeChildren]int)

func (g *GOL) buildNode(level int, size int, y, x int, label string) *node {

	children := nodeChildren{}

	if level == 2 {
		return newNode(g.buildLeafs(y, x), level, size, label)
	}

	halfSize := size / 2
	i := 0
	for yy := y; yy < y+2; yy++ {
		for xx := x; xx < x+2; xx++ {

			switch i {
			case nwPos:
				children.nw = g.buildNode(level-1, halfSize, y, x, "nw")
			case nePos:
				children.ne = g.buildNode(level-1, halfSize, y, x+halfSize, "ne")
			case swPos:
				children.sw = g.buildNode(level-1, halfSize, y+halfSize, x, "sw")
			case sePos:
				children.se = g.buildNode(level-1, halfSize, y+halfSize, x+halfSize, "se")
			}

			i++
		}
	}

	cache[children]++

	return newNode(children, level, size, label)
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

			tmp = cloneNode(tmp)

			switch i {
			case nwPos:
				tmp.label += fmt.Sprintf("_nw_%d_%d", yy, xx)
				children.nw = tmp
			case nePos:
				tmp.label += fmt.Sprintf("_ne_%d_%d", yy, xx)
				children.ne = tmp
			case swPos:
				tmp.label += fmt.Sprintf("_sw_%d_%d", yy, xx)
				children.sw = tmp
			case sePos:
				tmp.label += fmt.Sprintf("_se_%d_%d", yy, xx)
				children.se = tmp
			}

			i++
		}
	}

	cache[children]++

	return children
}

func cloneNode(n *node) *node {
	return &node{
		children:   n.children,
		level:      n.level,
		state:      n.state,
		population: n.population,
		size:       n.size,
		label:      n.label,
	}
}
