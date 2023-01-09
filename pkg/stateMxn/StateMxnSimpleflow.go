package stateMxn

/*
StateMxnSimpleflow

This is a special implementation of a statemxn, where each state should have 2 trasitions:

	[0] "Ok" state
	[-1] "Nok" state

and each state should be precreated with some execution-handler, to do some processing.

The state machine will take each state, execute its handlers with state.Activate() and
  - if there is no error, it will change to "Ok" state which assumed to be sm.transitionsMap[state.GetName()][0], and progress from there
  - if there is an error, it will change to "Nok" state which is assumed to be sm.transitionsMap[state.GetName()][-1]

So overall, this statemachine should take some precreated states, each with 2 transitions and 1-or-more-handlersExec, and will automatically
progress the execution from state to state, until it reaches a finalstate, or an error occurs.

There is no smsf.Change() method, because the state machine will start executing the initial state when NewStateMxnSimpleFlow() is called, and
then will automatically progress through the states, until it reaches a final state or an error occurs.

Examples: see main.go example4
*/
type StateMxnSimpleflow struct {
	smg *StateMxnGeneric
}

// Will create a new StateMxnSimpleflow, and will automatically progress through the states, until it reaches a final state or an error occurs.
func NewStateMxnSimpleFlow(transitionsMap map[string][]string, precreatedStates map[string]*State) (*StateMxnSimpleflow, error) {
	smsf := &StateMxnSimpleflow{}

	smg, err := NewStateMxnGeneric(transitionsMap, precreatedStates)
	smsf.smg = smg
	if err != nil {
		return smsf, err
	}
	return smsf, nil
}

// This function will automatically progress through the states, until it reaches a final state or an error occurs.
// It will return the error, if any.
func (smsf *StateMxnSimpleflow) ChangeToInitialStateAndAutoprogressToOtherStates(initialstateName string) error {
	hasOkNokTransitionsFunc := func(stateName string) (hasOkNokTransitions bool, OkStatename string, NokStatename string) {
		if len(smsf.smg.GetTransitionsMap()[stateName]) < 2 {
			return false, "", ""
		}
		OkStatename = smsf.smg.GetTransitionsMap()[stateName][0]
		NokStatename = smsf.smg.GetTransitionsMap()[stateName][len(smsf.smg.GetTransitionsMap()[stateName])-1]
		return true, OkStatename, NokStatename
	}

	a_state := initialstateName
	for {
		hasOkNokTransitions, OkStatename, NokStatename := hasOkNokTransitionsFunc(a_state)
		serr := smsf.smg.Change(a_state)
		// NOTE: serr may come from handler-error or another-error. We assume its a handler-error without additional checks
		if serr == nil {
			if !hasOkNokTransitions {
				// this state has no transitions defined, so this is a final-state, return nil
				return nil
			}
			a_state = OkStatename
		} else {
			if !hasOkNokTransitions {
				// this state has no transitions defined, so this is a final-state, return serr
				return serr
			}
			a_state = NokStatename
		}
	}
}
func (smsf *StateMxnSimpleflow) GetCurrentState() *State {
	return smsf.smg.GetCurrentState()
}
func (smsf *StateMxnSimpleflow) GetHistoryOfStates() HistoryOfStates {
	return smsf.smg.GetHistoryOfStates()
}
func (smsf *StateMxnSimpleflow) GetTransitionsMap() map[string][]string {
	return smsf.smg.GetTransitionsMap()
}
func (smsf *StateMxnSimpleflow) Is(stateName string) (bool, error) {
	return smsf.smg.Is(stateName)
}
