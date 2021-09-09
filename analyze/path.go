// package analyze calculate minimal survivable path
package analyze

import (
	"github.com/yourbasic/graph"
)

const (
	OpenTilePoint  = 0 // labyrinth level essence - open tile
	WallPoint      = 1 // labyrinth level essence - wall
	PitTrapPoint   = 2 // labyrinth level essence - kind of trap
	ArrowTrapPoint = 3 // labyrinth level essence - kind of trap
	HeroPoint      = 4 // labyrinth level essence - hero marker
)

// convert labyrinth level data from request into graph, set of edges.
// return length of minimal survivable path.
func MinSurvivablePathLen(labyrinthData [][]int) (mspLength int) {
	var (
		model    map[int]map[int]int
		lastNode int

		edges       *graph.Mutable
		edgesSorted *graph.Immutable

		dist []int
	)

	// convert input data into graph object
	model, lastNode = buildGraph(labyrinthData)
	edges = buildEdges(model)
	edgesSorted = graph.Sort(edges)

	// calculate distance between vertices
	dist = make([]int, edgesSorted.Order())
	graph.BFS(edgesSorted, 0, func(v, w int, _ int64) {
		dist[w] = dist[v] + 1
	})

	// get length for minimal survivable path
	mspLength = dist[lastNode]

	return mspLength
}

// numerated vertices of input labyrinth level data
// return generated object and number top vertex
func buildGraph(labyrinthData [][]int) (vertices map[int]map[int]int, lastNode int) {
	var (
		num int
	)

	vertices = make(map[int]map[int]int)
	num = 0

	// loop in backward order labyrinth level's lines
	for y := len(labyrinthData) - 1; y > -1; y-- {
		var (
			line map[int]int
		)

		line = make(map[int]int)

		// loop each line and numerating vertices
		for x := 0; x <= len(labyrinthData[0])-1; x++ {
			if labyrinthData[y][x] == HeroPoint {
				// hero position is always start position
				line[x] = 0
			} else {
				if labyrinthData[y][x] == WallPoint {
					// marking wall for don't use on next step
					line[x] = -1
				} else {
					// set valid number as name and increase for next vertex
					num++
					line[x] = num
				}
			}
		}

		vertices[y] = line
	}

	return vertices, num
}

// convert numerated vertices into graph object.
// return graph object
func buildEdges(vertices map[int]map[int]int) (model *graph.Mutable) {
	var (
		height, width int
	)

	// physical size of labyrinth level
	height = len(vertices)
	width = len(vertices[0])

	model = graph.New(height * width)

	// loop in backward order labyrinth level's lines
	for y := height - 1; y > 0; y-- {
		// loop each line and numerating vertices
		for x := 0; x <= (width-1)-1; x++ {
			var (
				currentVertex, upVertex, leftVertex int
			)

			currentVertex = vertices[y][x]
			leftVertex = vertices[y][x+1]
			upVertex = vertices[y-1][x]

			// do nothing if current essence is wall
			if currentVertex == -1 {
				continue
			}

			// add edge if right essence is not wall
			if leftVertex != -1 {
				model.AddBoth(currentVertex, leftVertex)
			}

			// add edge if up essence is not wall
			if upVertex != -1 {
				model.AddBoth(currentVertex, upVertex)
			}
		}
	}

	return model
}
