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
			node = node.BestChild(140)
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
	// If the node is nil or the game is already in a terminal state, return the reward for this state.
	if node == nil || node.state.IsTerminal() {
		return node.state.GetReward()
	}

	// Start simulation with the current state of the node.
	currentState := node.state

	count := 0
	// Continue the simulation until a terminal state is reached.
	for !currentState.IsTerminal() {

		// Iterate over all players to simulate their actions in turn.
		for _, player := range currentState.Snakes {
			// Get possible actions for the current player.
			possibleActions := currentState.GetPossibleActionsForPlayer(player.ID)

			// If no possible actions are available, return the reward for the current state.
			if len(possibleActions) == 0 {
				return currentState.GetReward()
				// continue
			}

			// Select a random action from the possible actions for the current player.
			action := possibleActions[rand.Intn(len(possibleActions))]

			// Apply the selected action to get the next state.
			currentState = *currentState.TakeActionForPlayer(action, player.ID)
		}
		count++
	}
	// Return the reward for the final state after the simulation reaches a terminal state.
	// fmt.Println("Simulation depth: ", count)
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
	for currentTime.Sub(startTime) < 200*time.Millisecond {
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
