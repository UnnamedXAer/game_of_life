package main

type nodeChildren struct {
	nw, ne, sw, se *node
}

type node struct {
	children   nodeChildren // children nodes, nil if leaf node
	level      uint         // 1 for leaf, 2 for node consisting of four leave, >= 3 for others
	state      cellState    // only matters if a leaf node
	population int          // numbers of active nodes expected to be 0..4
	next       *node        // nil if not computed yet
}

type nodeMap map[nodeChildren]*node

var (
	OnLeaf  = &node{state: aliveCell, level: 1}
	OffLeaf = &node{state: deadCell, level: 1}
)

func newNode(nw, ne, sw, se *node, level uint) *node {
	population := nw.population + ne.population + sw.population + se.population

	state := deadCell
	if population > 0 {
		state = aliveCell
	}

	return &node{
		children: nodeChildren{
			nw: nw,
			ne: ne,
			sw: sw,
			se: se,
		},
		level:      level,
		state:      state,
		population: population,
	}

}

func newLeafNode(state cellState) *node {
	population := 0

	if state == aliveCell {
		population = 1
	}

	return &node{
		level:      0,
		state:      state,
		population: population,
	}
}

func (n *node) createCenterSubnode() *node {

	// return newNode(n.nw, n.ne, n.sw, n.se, n.level-1)
	return nil
}
