package stateMxn

import "fmt"

type StateMxnGeneric struct {
	transitionsMap  map[string][]string
	precreatedStates map[string]*State
	currentState    *State
	historyOfStates []*State
}

func NewStateMxnGeneric(transitionsMap map[string][]string, initialStateName string) (*StateMxnGeneric, error) {
	smg := &StateMxnGeneric{}

	// FutureImprovement: Assure transitionsMap is valid

	// Assure the initial state is in the transitions map keys
	if !smg.verifyIfStateIsInTransitionsMap(initialStateName) {
		return nil, fmt.Errorf("initial state '%s' is not in the transitions map keys:\n%s", initialStateName, spew.Sdump(transitionsMap))
	}
	initialState := NewState(initialStateName)

	smg.transitionsMap = transitionsMap
	smg.currentState = initialState

	return smg, nil
}

//func (smg *StateMxnGeneric) Change

func (smg *StateMxnGeneric) TransitionsMap() map[string][]string {
	return smg.transitionsMap
}

func (smg *StateMxnGeneric) CurrentState() *State {
	return smg.currentState
}

func (smg *StateMxnGeneric) HistoryOfStates() []*State {
	return smg.historyOfStates
}

func (smg *StateMxnGeneric) verifyIfStateIsInTransitionsMap(stateName string) bool {
	_, ok := smg.transitionsMap[stateName]
	return ok
}
