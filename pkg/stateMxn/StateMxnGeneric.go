package stateMxn

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

/*
USAGE:

	// ===== States and Tansitions =====
	// Init
	//   |--> Running
	//           |------------> FinishedOk
	//		     |------------> FinishedNok
	//   |--------------------> FinishedNok

	// Create a StateMxnGeneric
	transitionsMap := map[string][]string{
			"Init": []string{"Running", "FinishedNok"},
			"Running": []string{"FinishedOk", "FinishedNok"},
		},
	initialStateName := "Init"

	smg, err := NewStateMxnGeneric(transitionsMap, initialStateName, nil)

	// Get current state name
	currentStateName := smg.CurrentState().Name()	// "Init"
	// Change state
	err := smg.Change("Running")
	// Get current state name
	currentStateName := smg.CurrentState().Name()	// "Running"
*/
type StateMxnGeneric struct {
	transitionsMap   map[string][]string
	precreatedStates map[string]*State
	currentState     *State
	historyOfStates  []*State
}

func NewStateMxnGeneric(transitionsMap map[string][]string, initialStateName string, precreatedStates map[string]*State) (*StateMxnGeneric, error) {
	smg := &StateMxnGeneric{}

	// FutureImprovement: Assure transitionsMap is valid

	// FutureImprovement: Assure precreatedStates is valid

	// define smg.precreatedStates from precreatedStates or from a new map
	if precreatedStates == nil {
		smg.precreatedStates = make(map[string]*State)
	} else {
		smg.precreatedStates = precreatedStates
	}

	// define initialState from precreatedStates or create-and-store a new state for it
	initialState, err := smg.getOrCreateState(initialStateName)
	if err != nil {
		return nil, err
	}

	smg.transitionsMap = transitionsMap
	smg.currentState = initialState

	return smg, nil
}

// TODO: fix this
func (smg *StateMxnGeneric) Change(newStateName string) error {
	// Assure newStateName is in the transitions map keys
	if !smg.verifyIfStatenameIsInTransitionsMap(newStateName) {
		return fmt.Errorf("state '%s' is not in the transitions map keys:\n%s", newStateName, spew.Sdump(smg.transitionsMap))
	}

	// Assure newStateName is in the transitions map values
	if !smg.verifyIfStatenameIsInTransitionsMapValues(newStateName) {
		return fmt.Errorf("state '%s' is not in the transitions map values:\n%s", newStateName, spew.Sdump(smg.transitionsMap))
	}

	// Get newState from precreatedStates or create-and-store a new state for it
	newState, err := smg.getOrCreateState(newStateName)
	if err != nil {
		return err
	}

	// Change state
	smg.currentState = newState
	smg.historyOfStates = append(smg.historyOfStates, newState)

	return nil
}

func (smg *StateMxnGeneric) TransitionsMap() map[string][]string {
	return smg.transitionsMap
}

func (smg *StateMxnGeneric) CurrentState() *State {
	return smg.currentState
}

func (smg *StateMxnGeneric) HistoryOfStates() []*State {
	return smg.historyOfStates
}

func (smg *StateMxnGeneric) verifyIfStatenameIsInTransitionsMap(stateName string) bool {
	_, ok := smg.transitionsMap[stateName]
	return ok
}

// return the state from:
// - precreatedStates if it exists there
// - or create-and-store a new state into precreatedStates
func (smg *StateMxnGeneric) getOrCreateState(stateName string) (*State, error) {
	// Assure stateName is in the transitions map keys
	if !smg.verifyIfStatenameIsInTransitionsMap(stateName) {
		return nil, fmt.Errorf("state '%s' is not in the transitions map keys:\n%s", stateName, spew.Sdump(smg.transitionsMap))
	}

	// Get stateName from precreatedStates or create-and-store a new state for it
	if _, ok := smg.precreatedStates[stateName]; !ok {
		smg.precreatedStates[stateName] = NewState(stateName)
	}
	state := smg.precreatedStates[stateName]

	return state, nil
}
