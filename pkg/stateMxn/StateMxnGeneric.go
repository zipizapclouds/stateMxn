package stateMxn

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

/*
Features that exist in any StateMxnGeneric:

- per-stage-handlers: each stage can have a handlerBegin, handlerExec and handlerEnd. The execution order is: handlerBegin, handlerExec, handlerEnd

- stage-output-input chaining: the *ouput* from a state is copied to the *input* of the next state

- stage-data: each state has a data map[string]interface{} where you can store any data meaningfull for that state

# Use `smg.GetHistoryOfStates().DisplayStatesFlow()` to display the state-flow of the state-machine

Use `smg.Is("^Finished"")` to check if the state-machine is in a specific state (regexp)

----
Ideas for future improvements:
- StateMxnGeneric.data: to share/persist data inter-states
*/
type StateMxnGeneric struct {
	transitionsMap map[string][]string

	// map[<statename>]*State
	precreatedStates map[string]*State
	currentState     *State
	historyOfStates  HistoryOfStates
}

// precreatedStates can be nil
func NewStateMxnGeneric(transitionsMap map[string][]string, precreatedStates map[string]*State) (*StateMxnGeneric, error) {
	smg := &StateMxnGeneric{}

	// FutureImprovement: Assure transitionsMap is valid

	// FutureImprovement: Assure precreatedStates is valid

	// Define smg.transitionsMap
	smg.transitionsMap = transitionsMap

	// Define smg.precreatedStates from precreatedStates or from a new map
	if precreatedStates == nil {
		smg.precreatedStates = make(map[string]*State)
	} else {
		smg.precreatedStates = precreatedStates
	}
	return smg, nil
}

// Changes from current state to nextStateName
func (smg *StateMxnGeneric) Change(nextStateName string) error {
	// Note: this function may be called to set initialstate in which case smg.currentState is nil
	//
	// Performs safety-validations:
	// - check if its valid the transition change from currentState to nextStateName
	//
	// and execute the change, updating currentState, historyOfStates and possibly precreatedStates, by:
	// - creating a nextState, from a copy-or-a-new-state in precreatedStates
	// - appending nextState to historyOfStates
	// - setting currentState = nextState
	// - call currentState.Activate(inputs)
	//---------------------------------------------------------------------------------------------

	// Performs safety-validations:
	// - check if its valid the transition change from currentState to nextStateName
	{
		// -- check if nextStateName is a valid stateName
		// -- check if currentState is a valid sourcestate
		// -- check if nextState is a valid destinationstate, from currentState

		// - check if nextStateName is a valid stateName
		err := smg.verifyIfValidStatename(nextStateName)
		if err != nil {
			return err
		}

		if smg.currentState == nil {
			// When smg.currentState == nil this function is called to set initialstate, and then
			// .we accept any nextStateName as valid (dont check if valid sourcestate or valid transition)
			// .
		} else {
			// -- check if currentState is a valid sourcestate
			err = smg.verifyIfValidSourcestate(smg.currentState.GetName())
			if err != nil {
				return err
			}
			// -- check if nextState is a valid destinationstate, from currentState
			err = smg.verifyIfValidTransition(smg.currentState.GetName(), nextStateName)
			if err != nil {
				return err
			}
		}
	}

	// and execute the change, by:
	// - creating a nextState, from a copy-or-a-new-state in precreatedStates
	nextState, err := smg.getStatecopyFromPrecreatedstatesOrNew(nextStateName)
	if err != nil {
		return err
	}
	oldState := smg.currentState
	var inputs StateInputs
	if oldState == nil {
		// When smg.currentState == nil this function is called to set initialstate
		inputs = make(StateInputs)
	} else {
		oldState_outputs := oldState.GetOutputs()
		inputs = oldState_outputs.Convert2Inputs()
	}

	// - appending nextState to historyOfStates
	smg.historyOfStates = append(smg.historyOfStates, nextState)
	// - setting currentState = nextState
	smg.currentState = nextState

	// - call currentState.Activate(inputs)
	_, err = smg.currentState.Activate(inputs)
	if err != nil {
		return err
	}

	return nil
}

// Is returns true if the smg.CurrentState.Name() matches the given regexp
//
// currentStateNameRegexp - is a regexp RE2 as described at https://golang.org/s/re2syntax
// the same as used by regexp.MatchString()
//
// Examples:
//   - regexp "Init" 						matches Name "Init"
//   - regexp "Init|Running" 				matches Name "Init"
//   - regexp "Init|Running|FinishedOk" 	matches Name "Init"
//   - regexp "Finished" 					matches Name "FinisedOk" or "FinishedNok"
func (smg *StateMxnGeneric) Is(currentStateNameRegexp string) (bool, error) {
	return smg.GetCurrentState().Is(currentStateNameRegexp)
}

func (smg *StateMxnGeneric) GetTransitionsMap() map[string][]string {
	return smg.transitionsMap
}
func (smg *StateMxnGeneric) GetCurrentState() *State {
	return smg.currentState
}

// NOTE: historyOfStates[-1] == currentState
func (smg *StateMxnGeneric) GetHistoryOfStates() HistoryOfStates {
	return smg.historyOfStates
}

// Returns if stateName is a valid source or destination state (ie, either in the transitions map keys or in the transitions map values)
func (smg *StateMxnGeneric) verifyIfValidStatename(stateName string) error {
	// Check if stateName is in the transitions map keys
	if err := smg.verifyIfValidSourcestate(stateName); err == nil {
		// found in transitions map keys
		return nil
	}

	// Check if stateName is in the transitions map values
	for _, stateNames := range smg.transitionsMap {
		for _, stateNameInValues := range stateNames {
			if stateNameInValues == stateName {
				// found in transitions map values
				return nil
			}
		}
	}
	return fmt.Errorf("state '%s' is unrecognized, invalid! The transitionsMap is:\n%s", stateName, spew.Sdump(smg.transitionsMap))
}

func (smg *StateMxnGeneric) verifyIfValidSourcestate(stateName string) error {
	// Check if stateName is in the transitions map keys
	if _, ok := smg.transitionsMap[stateName]; ok {
		return nil
	} else {
		return fmt.Errorf("stateName '%s' is not a valid sourcestate! The transitionsMap is:\n%s", stateName, spew.Sdump(smg.transitionsMap))
	}
}

func (smg *StateMxnGeneric) verifyIfValidTransition(source_stateName string, destination_stateName string) error {
	// Check if destionationstateName is in the transitions map values, from sourcestateName
	if possibleDeststates, ok := smg.transitionsMap[source_stateName]; ok {
		for _, a_possibleDeststate := range possibleDeststates {
			if a_possibleDeststate == destination_stateName {
				return nil
			}
		}
	}
	return fmt.Errorf("transition from sourcestate '%s' to destinationstate '%s' is not valid! The transitionsMap is:\n%s", source_stateName, destination_stateName, spew.Sdump(smg.transitionsMap))
}

// Performs some safety-validations:
// - if stateName is valid
//
// and then return the state (correspoding to stateName), from:
// - if precreatedStates contains that state, then return a copy of it
// or
// - create-and-store into precreatedStates a new state, and then return a copy of it
func (smg *StateMxnGeneric) getStatecopyFromPrecreatedstatesOrNew(stateName string) (*State, error) {
	// Performs some safety-validations:
	// - if stateName is valid
	{
		err := smg.verifyIfValidStatename(stateName)
		if err != nil {
			return nil, err
		}
	}

	// and then return the state (corresponding to stateName), from:
	// - if precreatedStates contains that state, then return a copy of it
	// or
	// - create-and-store into precreatedStates a new state, and then return a copy of it
	if stateCandidate, ok := smg.precreatedStates[stateName]; ok {
		// precreatedStates contains that state, lets return a copy of it
		stateCopy := stateCandidate.deepcopy()
		return stateCopy, nil
	} else {
		// precreatedStates does not contain that state
		// - create-and-store into precreatedStates a new state, and then return a copy of it
		state := NewState(stateName)
		smg.precreatedStates[stateName] = state

		stateCopy := state.deepcopy()
		return stateCopy, nil
	}
}
