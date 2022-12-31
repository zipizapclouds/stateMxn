package stateMxn

type State struct {
	name string
}

func NewState(name string) *State {
	return &State{name: name}
}

func (s *State) Name() string {
	return s.name
}
