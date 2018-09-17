package main

import "sort"

func main() {

	println(bruteForce(4, 5, []int{0, 1, 1, 2, 3}, []int{1, 2, 3, 4, 4}, []int{3, 2, 4, 2, 1}))
}

func bruteForce(n, m int, x, y, w []int) int {
	parent := make([]int, n+1)
	for i := 0; i <= n; i++ {
		parent[i] = i
	}

	edges := make([]Edge, 0, m)
	for i := 0; i < m; i++ {
		edges = append(edges, Edge{
			start:  x[i],
			end:    y[i],
			weight: w[i],
		})
	}
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].weight < edges[j].weight
	})

	for _, edge := range edges {
		if !Connected(edge.start, edge.end, parent) {
			Connect(edge.start, edge.end, parent)
			if Connected(0, n, parent) {
				return edge.weight
			}
		}
	}
	println("finish")
	return -1

}

func findHead(parent []int, index int) int {
	for parent[index] != index {
		parent[index] = parent[parent[index]]
		index = parent[index]
	}
	return parent[index]
}

type Edge struct {
	start  int
	end    int
	weight int
}

func Connect(i, j int, parent []int) {
	parent[findHead(parent, i)] = findHead(parent, j)
}
func Connected(i, j int, parent []int) bool {
	return findHead(parent, i) == findHead(parent, j)
}
