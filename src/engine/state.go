// VoV engine
package engine

// Game state interface
type State interface {
	OnInit() bool
	OnQuit() bool

	String() string

	HandleEvents()
	Update()
	Draw()
}

// State machine
type StateMachine struct {
	states []State
}

// Push state
func (g *StateMachine) Push(state State) {
	g.states = append(g.states, state)

	g.states[g.State()].OnInit()
}

// Pop state
func (g *StateMachine) Pop() {
	if g.Size() != 0 {
		if g.states[g.State()].OnQuit() {
			tmp, states := g.states[g.State()], g.states[:g.State()]
			g.states = append(states, tmp)
		}
	}
}

// Change state
func (g *StateMachine) Change(state State) {
	if g.Size() != 0 {
		if g.states[g.State()].String() == state.String() {
			return
		}

		if g.states[g.State()].OnQuit() {
			tmp, states := g.states[g.State()], g.states[:g.State()]
			g.states = append(states, tmp)
		}
	}

	g.Push(state)
}

// Returns size
func (g *StateMachine) Size() int {
	return len(g.states)
}

// Returns state
func (g *StateMachine) State() int {
	return len(g.states) - 1
}

// Handles state events
func (g *StateMachine) HandleEvents() {
	if g.Size() != 0 {
		g.states[g.State()].HandleEvents()
	}
}

// Updates state
func (g *StateMachine) Update() {
	if g.Size() != 0 {
		g.states[g.State()].Update()
	}
}

// Draws state
func (g *StateMachine) Draw() {
	if g.Size() != 0 {
		g.states[g.State()].Draw()
	}
}
