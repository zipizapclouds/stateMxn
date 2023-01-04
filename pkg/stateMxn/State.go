package stateMxn

import "regexp"

type State struct {
	name string
}

func NewState(name string) *State {
	return &State{name: name}
}

func (s *State) Name() string {
	return s.name
}

// Is returns true if the state name matches the given regexp
//
// stateNameRegexp - is a regexp RE2 as described at https://golang.org/s/re2syntax
// the same as used by regexp.MatchString()
//
// Examples:
//   - regexp "Init" 						matches Name "Init"
//   - regexp "Init|Running" 				matches Name "Init"
//   - regexp "Init|Running|FinishedOk" 	matches Name "Init"
//   - regexp "Finished" 					matches Name "FinisedOk" or "FinishedNok"
func (s *State) Is(stateNameRegexp string) (bool, error) {
	return regexp.MatchString(stateNameRegexp, s.name)
}

func (s *State) deepcopy() *State {
	// NOTE: deepcopy libs like https://github.com/barkimedes/go-deepcopy
	//       dont copy unexported fields - so we need to define our own deepcopy() method
	// 	     for the type, in the package where the type is defined

	// all truct fields, both exported and unexported, need to be copied here
	stateCopy := &State{
		name: s.name,
	}
	return stateCopy
}
