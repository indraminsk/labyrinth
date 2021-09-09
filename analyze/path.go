package analyze

import (
	"github.com/yourbasic/graph"
)

const (
	OpenTilePoint  = 0
	WallPoint      = 1
	PitTrapPoint   = 2
	ArrowTrapPoint = 3
	HeroPoint      = 4
)

func MinSurvivablePathLen(labyrinthData [][]int) (mspLength int) {
	var (
		model    map[int]map[int]int
		lastNode int

		edges       *graph.Mutable
		edgesSorted *graph.Immutable

		dist []int
	)

	model, lastNode = buildGraph(labyrinthData)
	edges = buildEdges(model)
	edgesSorted = graph.Sort(edges)

	dist = make([]int, edgesSorted.Order())
	graph.BFS(edgesSorted, 0, func(v, w int, _ int64) {
		dist[w] = dist[v] + 1
	})

	mspLength = dist[lastNode]

	return mspLength
}

func buildGraph(labyrinthData [][]int) (vertices map[int]map[int]int, lastNode int) {
	var (
		num int
	)

	vertices = make(map[int]map[int]int)
	num = 0

	for y := len(labyrinthData) - 1; y > -1; y-- {
		var (
			line map[int]int
		)

		line = make(map[int]int)

		for x := 0; x <= len(labyrinthData[0])-1; x++ {
			if labyrinthData[y][x] == HeroPoint {
				line[x] = 0
			} else {
				if labyrinthData[y][x] == WallPoint {
					line[x] = -1
				} else {
					num++
					line[x] = num
				}
			}
		}

		vertices[y] = line
	}

	return vertices, num
}

func buildEdges(vertices map[int]map[int]int) (model *graph.Mutable) {
	var (
		height, width int
	)

	height = len(vertices)
	width = len(vertices[0])

	model = graph.New(height * width)

	for y := height - 1; y > 0; y-- {
		for x := 0; x <= (width-1)-1; x++ {
			var (
				currentVertex, upVertex, leftVertex int
			)

			currentVertex = vertices[y][x]
			leftVertex = vertices[y][x+1]
			upVertex = vertices[y-1][x]

			if currentVertex == -1 {
				continue
			}

			if leftVertex != -1 {
				model.AddBoth(currentVertex, leftVertex)
			}

			if upVertex != -1 {
				model.AddBoth(currentVertex, upVertex)
			}
		}
	}

	return model
}
