package main

import "math"

type State struct {
	MyID          string
	Snake         []Coord
	Head          Coord
	Health        int
	Food          []Coord
	Turn          int
	Snakes        []Battlesnake
	Board         map[string]int
	CurrentPlayer int
}

func NewState(data GameState) *State {
	return &State{
		MyID:          data.You.ID,
		Snake:         data.You.Body,
		Head:          data.You.Head,
		Health:        data.You.Health,
		Food:          data.Board.Food,
		Turn:          data.Turn,
		Snakes:        data.Board.Snakes,
		Board:         map[string]int{"width": data.Board.Width, "height": data.Board.Height},
		CurrentPlayer: 0,
	}
}

func (s *State) GetCurrentPlayer() int {
	return s.CurrentPlayer
}

func (s *State) GetPossibleActions() []string {
	possibleActions := []string{"up", "down", "left", "right"}
	s.AvoidWalls(&possibleActions)
	s.AvoidBody(&possibleActions)
	return possibleActions
}

func (s *State) GetPossibleActionsForPlayer(playerID string) []string {
	possibleActions := []string{"up", "down", "left", "right"}
	s.AvoidWallsForPlayer(&possibleActions, playerID)
	s.AvoidBodyForPlayer(&possibleActions, playerID)
	return possibleActions
}

func (s *State) AvoidWalls(possibleActions *[]string) {
	if s.Board["width"]-s.Head.X == 1 {
		removeAction("right", possibleActions)
	}
	if s.Head.X == 0 {
		removeAction("left", possibleActions)
	}
	if s.Board["height"]-s.Head.Y == 1 {
		removeAction("up", possibleActions)
	}
	if s.Head.Y == 0 {
		removeAction("down", possibleActions)
	}
}

func (s *State) AvoidWallsForPlayer(possibleActions *[]string, playerID string) {
	var playerHead Coord
	for _, snake := range s.Snakes {
		if snake.ID == playerID {
			playerHead = snake.Head
			break
		}
	}
	if s.Board["width"]-playerHead.X == 1 {
		removeAction("right", possibleActions)
	}
	if playerHead.X == 0 {
		removeAction("left", possibleActions)
	}
	if s.Board["height"]-playerHead.Y == 1 {
		removeAction("up", possibleActions)
	}
	if playerHead.Y == 0 {
		removeAction("down", possibleActions)
	}
}

func (s *State) AvoidBody(possibleActions *[]string) {
	left := Coord{s.Head.X - 1, s.Head.Y}
	right := Coord{s.Head.X + 1, s.Head.Y}
	up := Coord{s.Head.X, s.Head.Y + 1}
	down := Coord{s.Head.X, s.Head.Y - 1}

	var removedActions []string
	for _, snake := range s.Snakes {
		snakeNoTail := snake.Body
		snakeNoTail = snakeNoTail[:len(snakeNoTail)-1]

		if contains(left, snakeNoTail) {
			removeAction("left", possibleActions)
		}
		if contains(right, snakeNoTail) {
			removeAction("right", possibleActions)
		}
		if contains(up, snakeNoTail) {
			removeAction("up", possibleActions)
		}
		if contains(down, snakeNoTail) {
			removeAction("down", possibleActions)
		}

		if snake.ID == s.MyID {
			continue
		} else {
			if len(snake.Body) >= len(s.Snake) {
				// Check left direction
				dLeft := math.Abs(float64(left.X-snake.Head.X)) + math.Abs(float64(left.Y-snake.Head.Y))
				if dLeft <= 1 {
					removedActions = append(removedActions, "left")
					removeAction("left", possibleActions)
				}

				// Check right direction
				dRight := math.Abs(float64(right.X-snake.Head.X)) + math.Abs(float64(right.Y-snake.Head.Y))
				if dRight <= 1 {
					removedActions = append(removedActions, "right")
					removeAction("right", possibleActions)
				}

				// Check up direction
				dUp := math.Abs(float64(up.X-snake.Head.X)) + math.Abs(float64(up.Y-snake.Head.Y))
				if dUp <= 1 {
					removedActions = append(removedActions, "up")
					removeAction("up", possibleActions)
				}

				// Check down direction
				dDown := math.Abs(float64(down.X-snake.Head.X)) + math.Abs(float64(down.Y-snake.Head.Y))
				if dDown <= 1 {
					removedActions = append(removedActions, "down")
					removeAction("down", possibleActions)
				}
			}
		}
	}

	if len(*possibleActions) == 0 {
		*possibleActions = removedActions
	}
}

func (s *State) AvoidBodyForPlayer(possibleActions *[]string, playerID string) {
	var playerHead Coord
	for _, snake := range s.Snakes {
		if snake.ID == playerID {
			playerHead = snake.Head
			break
		}
	}

	left := Coord{playerHead.X - 1, playerHead.Y}
	right := Coord{playerHead.X + 1, playerHead.Y}
	up := Coord{playerHead.X, playerHead.Y + 1}
	down := Coord{playerHead.X, playerHead.Y - 1}

	var removedActions []string
	for _, snake := range s.Snakes {
		snakeNoTail := snake.Body
		snakeNoTail = snakeNoTail[:len(snakeNoTail)-1]

		if contains(left, snakeNoTail) {
			removeAction("left", possibleActions)
		}
		if contains(right, snakeNoTail) {
			removeAction("right", possibleActions)
		}
		if contains(up, snakeNoTail) {
			removeAction("up", possibleActions)
		}
		if contains(down, snakeNoTail) {
			removeAction("down", possibleActions)
		}
	}

	if len(*possibleActions) == 0 {
		*possibleActions = removedActions
	}
}

func (s *State) TakeAction(action string) *State {
	nextState := DeepCopyState(s)
	switch action {
	case "up":
		nextState.Head.Y++
	case "down":
		nextState.Head.Y--
	case "right":
		nextState.Head.X++
	case "left":
		nextState.Head.X--
	default:
		panic("Invalid Action!")
	}

	nextState.Snake = append([]Coord{nextState.Head}, nextState.Snake...)
	nextState.Snake = nextState.Snake[:len(nextState.Snake)-1]

	if contains(nextState.Head, nextState.Food) {
		nextState.Health = 100
		nextState.Food = removePosition(nextState.Head, nextState.Food)
		nextState.Snake = append(nextState.Snake, nextState.Snake[len(nextState.Snake)-1])
	} else {
		nextState.Health--
	}

	nextState.Turn++
	nextState.CurrentPlayer = (nextState.CurrentPlayer + 1) % len(nextState.Snakes)

	return nextState
}

func (s *State) TakeActionForPlayer(action string, playerID string) *State {
	nextState := DeepCopyState(s)
	var playerIndex int
	for i, snake := range nextState.Snakes {
		if snake.ID == playerID {
			playerIndex = i
			switch action {
			case "up":
				nextState.Snakes[i].Head.Y++
			case "down":
				nextState.Snakes[i].Head.Y--
			case "right":
				nextState.Snakes[i].Head.X++
			case "left":
				nextState.Snakes[i].Head.X--
			default:
				panic("Invalid Action!")
			}

			nextState.Snakes[i].Body = append([]Coord{nextState.Snakes[i].Head}, nextState.Snakes[i].Body...)
			nextState.Snakes[i].Body = nextState.Snakes[i].Body[:len(nextState.Snakes[i].Body)-1]

			if contains(nextState.Snakes[i].Head, nextState.Food) {
				nextState.Snakes[i].Health = 100
				nextState.Food = removePosition(nextState.Snakes[i].Head, nextState.Food)
				nextState.Snakes[i].Body = append(nextState.Snakes[i].Body, nextState.Snakes[i].Body[len(nextState.Snakes[i].Body)-1])
			} else {
				nextState.Snakes[i].Health--
			}

			break
		}
	}

	nextState.Turn++
	nextState.CurrentPlayer = (playerIndex + 1) % len(nextState.Snakes)

	return nextState
}

func (s *State) IsTerminal() bool {
	if len(s.GetPossibleActions()) == 0 || s.Health == 0 {
		return true
	}
	return false
}

func (s *State) GetReward() int {
	if s.IsTerminal() {
		return s.Turn
	}
	return 1000000
}

func toPositionSlice(data []interface{}) []Coord {
	positions := make([]Coord, len(data))
	for i, item := range data {
		pos := item.(map[string]interface{})
		positions[i] = Coord{X: pos["x"].(int), Y: pos["y"].(int)}
	}
	return positions
}

func toPosition(data map[string]interface{}) Coord {
	return Coord{X: data["x"].(int), Y: data["y"].(int)}
}

func toMapSlice(data []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, len(data))
	for i, item := range data {
		result[i] = item.(map[string]interface{})
	}
	return result
}

func removeAction(action string, actions *[]string) {
	for i := range *actions {
		if (*actions)[i] == action {
			(*actions)[i] = (*actions)[len(*actions)-1]
			(*actions) = (*actions)[:len(*actions)-1]
			return
		}
	}
}

func DeepCopyState(s *State) *State {
	newState := *s
	newState.Snake = make([]Coord, len(s.Snake))
	copy(newState.Snake, s.Snake)
	newState.Food = make([]Coord, len(s.Food))
	copy(newState.Food, s.Food)
	newState.Snakes = make([]Battlesnake, len(s.Snakes))
	copy(newState.Snakes, s.Snakes)
	newState.Board = make(map[string]int)
	for k, v := range s.Board {
		newState.Board[k] = v
	}
	return &newState
}

func contains(pos Coord, positions []Coord) bool {
	for _, p := range positions {
		if p == pos {
			return true
		}
	}
	return false
}

func removePosition(pos Coord, positions []Coord) []Coord {
	for i, p := range positions {
		if p == pos {
			return append(positions[:i], positions[i+1:]...)
		}
	}
	return positions
}
