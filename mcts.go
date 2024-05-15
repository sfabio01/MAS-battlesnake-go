package main

import (
	"fmt"
	"math/rand"
	"time"
)

const NUM_ITERATIONS = 100

func Traverse(node *Node) *Node {
	for len(node.children) > 0 {
		if node.isFullyExpanded {
			node = node.BestChild(14)
		} else {
			return Expand(node)
		}
	}
	return node
}

func Expand(node *Node) *Node {
	if node == nil || node.state.IsTerminal() {
		return node
	}
	possibleActions := node.state.GetPossibleActions()
	if len(possibleActions) == 0 {
		return node
	}
	for _, action := range possibleActions {
		if _, exists := node.children[action]; !exists {
			newState := node.state.TakeAction(action)
			return node.AddChild(*newState, action)
		}
	}
	node.isFullyExpanded = true
	randomActionIndex := rand.Intn(len(node.children))
	for _, child := range node.children {
		if randomActionIndex == 0 {
			return child
		}
		randomActionIndex--
	}
	return nil
}

func Simulate(node *Node) int {
	if node == nil || node.state.IsTerminal() {
		return node.state.GetReward()
	}
	currentState := node.state
	for !currentState.IsTerminal() {
		possibleActions := currentState.GetPossibleActions()
		if len(possibleActions) == 0 {
			return currentState.GetReward()
		}
		action := possibleActions[rand.Intn(len(possibleActions))]
		currentState = *currentState.TakeAction(action)
	}
	return currentState.GetReward()
}

func Backpropagate(node *Node, reward float64) {
	for node != nil {
		node.Update(reward)
		node = node.parent
	}
}

func GetNextMove(gameState GameState) string {
	currentGameState := NewState(gameState)
	rootNode := NewNode(*currentGameState, nil)
	i := 0
	currentTime := time.Now()
	startTime := time.Now()
	for currentTime.Sub(startTime) < 360*time.Millisecond {
		lastNode := Traverse(rootNode)
		newNode := Expand(lastNode)
		reward := Simulate(newNode)
		Backpropagate(newNode, float64(reward))
		currentTime = time.Now()
		i++
	}

	fmt.Println("Time: ", currentTime.Sub(startTime))
	fmt.Println("Num iterations: ", i)

	return rootNode.BestAction()
}
