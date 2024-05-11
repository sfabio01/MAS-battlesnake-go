package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Node struct {
	state           State
	parent          *Node
	children        map[string]*Node
	visitCount      int
	winScore        float64
	isFullyExpanded bool
}

func NewNode(state State, parent *Node) *Node {
	return &Node{
		state:           state,
		parent:          parent,
		children:        make(map[string]*Node),
		visitCount:      0,
		winScore:        0,
		isFullyExpanded: false,
	}
}

func (n *Node) AddChild(childState State, action string) *Node {
	childNode := NewNode(childState, n)
	n.children[action] = childNode
	return childNode
}

func (n *Node) Update(reward float64) {
	n.visitCount++
	n.winScore += reward
}

func (n *Node) BestChild(cParam float64) *Node {
	bestValue := math.Inf(-1)
	var bestNodes []*Node
	for _, child := range n.children {
		nodeValue := child.winScore/float64(child.visitCount) + cParam*math.Sqrt(2*math.Log(float64(n.visitCount))/float64(child.visitCount))
		if nodeValue > bestValue {
			bestValue = nodeValue
			bestNodes = []*Node{child}
		} else if nodeValue == bestValue {
			bestNodes = append(bestNodes, child)
		}
	}
	return bestNodes[rand.Intn(len(bestNodes))]
}

func (n *Node) BestAction() string {
	var maxAction string
	maxScore := math.Inf(-1)
	fmt.Println(n.children)
	for action, child := range n.children {
		score := child.winScore / float64(child.visitCount)
		if score > maxScore {
			maxScore = score
			maxAction = action
		}
	}
	return maxAction
}

func (n *Node) String() string {
	return fmt.Sprintf("Visits: %d Score: %f", n.visitCount, n.winScore/float64(n.visitCount))
}
