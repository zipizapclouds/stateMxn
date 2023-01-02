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

func (s *State) deepcopy() *State {
	// NOTE: deepcopy libs like https://github.com/barkimedes/go-deepcopy
	//       dont copy unexported fields - so we need to define our own deepcopy() method
	// 	     for the type, in the package where the type is defined

	// all truct fields, both exported and unexported, need to be copied
	stateCopy := &State{
		name: s.name,
	}
	return stateCopy
}
