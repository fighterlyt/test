package main

import (
	"github.com/davecgh/go-spew/spew"
	"log"
	"strconv"
	"strings"
)

var (
	data = `-9,6,-9,3,-8,9,-2,-0,-4,12`
)

func main() {
	input := parse(data)
	spew.Dump(longest(input))
}

func parse(data string) []*node {
	elements := strings.Split(data, ",")

	result := make([]*node, 0, len(elements))
	for _, element := range elements {

		if element[0:1] == "-" {
			value, _ := strconv.ParseInt(element[1:], 10, 64)
			result = append(result, &node{
				value:   int(value),
				goRight: false,
			})
		} else {
			value, _ := strconv.ParseInt(element, 10, 64)
			result = append(result, &node{
				value:   int(value),
				goRight: true,
			})
		}
	}

	return result
}

type node struct {
	value   int
	goRight bool
	To      []int
}

func longest(input []*node) []int {
	for i, node := range input {
		if node.goRight {
			for j := i + 1; j < len(input); j++ {
				if input[j].value > node.value {
					input[i].To = append(input[i].To, j)
				}
			}
		} else {
			for j := i - 1; j >= 0; j-- {
				if input[j].value > node.value {
					input[i].To = append(input[i].To, j)
				}
			}
		}
	}

	spew.Dump("afterRelation", input)
	result := make([]int, 0, 10)

	candidates := make([]int, 0, len(input))

	for i := range input {
		candidates = append(candidates, i)
	}
	best := 0
	for {

		log.Println("find", candidates)

		candidates = find(candidates, input)

		log.Println("findAfter", candidates)

		switch len(candidates) {
		case 0:
			return result
		case 1:
			best = candidates[0]
		default:
			best = findBest(candidates, input)
		}

		result = append(result, best)
		candidates = candidates[:]

		for i := range input {
			candidates = append(candidates, i)
		}
		candidates = remove(candidates, result)
	}
}

func findBest(candidate []int, input []*node) []int {

}

func find(candidates []int, input []*node) []int {

	processed := make(map[int]struct{}, 10)

	restart := false
	for {
		restart = false
		for i := 0; i < len(candidates); i++ {

			index := candidates[i]
			if _, exist := processed[index]; exist {
				continue
			}

			node := input[index]
			processed[index] = struct{}{}
			candidates = remove(candidates, node.To)
			if len(node.To) > 0 {
				restart = true
				break
			}
		}
		// log.Println("第2个for", restart)
		if !restart {
			break
		}
	}

	return remove(candidates, []int{0, 9})
}

func remove(src, target []int) []int {
	// log.Println("remove", src, target)
	result := make([]int, 0, len(src))

	for _, elem := range src {
		found := false
		for _, targetElem := range target {
			if elem == targetElem {
				found = true
				break
			}
		}
		if !found {
			result = append(result, elem)
		}
	}

	// log.Println("remove-Result", result)
	return result
}
