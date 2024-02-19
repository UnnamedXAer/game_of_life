package main

type nodeChildren struct {
	nw, ne, sw, se *node
}

type node struct {
	children   nodeChildren // children nodes, nil if leaf node
	level      int          // 1 for leaf, 2 for node consisting of four leave, >= 3 for others
	state      cellState    // only matters if a leaf node
	population int          // numbers of active nodes expected to be 0..4
	next       *node        // nil if not computed yet
}

type nodeMap map[nodeChildren]*node

var (
	aliveLeaf = &node{state: aliveCell, level: 1, population: 1}
	deadLeaf  = &node{state: deadCell, level: 1, population: 0}
)

func newNode(children nodeChildren, level int) *node {
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
	}

}

// func newLeafNode(state cellState) *node {
// 	population := 0

// 	if state == aliveCell {
// 		population = 1
// 	}

// 	return &node{
// 		level:      1,
// 		state:      state,
// 		population: population,
// 	}
// }

const (
	nwPos = iota
	nePos
	swPos
	sePos
)

func (g *GOL) buildNode(level int, size int, y, x int) *node {

	children := nodeChildren{}

	if level == 2 {
		return newNode(g.buildLeafs(y, x), level-1)
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

	return newNode(children, level)
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
