package analyser

type WeightedGraph struct {
	matrix map[string]map[string]int64
}

func newWeightedGraph() *WeightedGraph {
	return &WeightedGraph{make(map[string]map[string]int64)}
}

func (g *WeightedGraph) add(source string, target string) {
	count := g.matrix[source][target]
	g.matrix[source][target] = count + 1
}

func (g *WeightedGraph) remove(source string, target string) {
	count := g.matrix[source][target]
	if count > 1 {
		g.matrix[source][target] = count - 1
	} else {
		delete(g.matrix[source], target)
		if len(g.matrix[source]) == 0 {
			delete(g.matrix, source)
		}
	}
}
