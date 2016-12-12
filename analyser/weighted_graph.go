package analyser

import (
	"sync"
	"time"
)

type WeightedGraph struct {
	matrix map[string]map[string]int64
	cache  []cacheEntry
	lock   *sync.RWMutex
}

func newWeightedGraph() *WeightedGraph {
	return &WeightedGraph{
		make(map[string]map[string]int64),
		[]cacheEntry{},
		&sync.RWMutex{},
	}
}

type cacheEntry struct {
	timestamp time.Time
	source    string
	target    string
}

func (g *WeightedGraph) add(source string, target string) {
	g.lock.Lock()
	defer g.lock.Unlock()
	if len(g.matrix[source]) == 0 {
		g.matrix[source] = make(map[string]int64)
	}
	count := g.matrix[source][target]
	g.matrix[source][target] = count + 1
	g.addToCache(source, target)
}

func (g *WeightedGraph) addToCache(source string, target string) {
	if len(g.cache) == 1000 {
		g.cache = g.cache[1:]
	}
	g.cache = append(g.cache, cacheEntry{time.Now(), source, target})
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

func (g *WeightedGraph) Matrix() map[string]map[string]int64 {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.matrix
}
