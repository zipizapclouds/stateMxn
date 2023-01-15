package stateMxn

import "fmt"

/*
StateMxnSimpleflow

This is a special implementation of a statemxn, where each state should have 2 trasitions:

	[0] "Ok" state
	[-1] "Nok" state

and each state should be precreated with some execution-handler, to do some processing.

The state machine will take each state, execute its handlers with state.Activate() and
  - if there is no error, it will change to "Ok" state which assumed to be sm.transitionsMap[state.GetName()][0], and progress from there
  - if there is an error, it will change to "Nok" state which is assumed to be sm.transitionsMap[state.GetName()][-1]

So overall, this statemachine should take some precreated states, each with 2 transitions and 0-or-more-handlersExec, and will automatically
progress the execution from state to state, until it reaches a finalstate, or an error occurs.

The smsf.ChangeToInitialStateAndAutoprogressToOtherStates() method should be used to start the state machine, and it will automatically progress through
the states, until it reaches a final state or an error occurs.

The smsf.Change() method should not be used - use smsf.ChangeToInitialStateAndAutoprogressToOtherStates() instead.

Examples: see main.go example4, .. and example8
*/
type StateMxnSimpleflow struct {
	*StateMxnGeneric
}

// Will create a new StateMxnSimpleflow
func NewStateMxnSimpleFlow(smxName string, transitionsMap map[string][]string, precreatedStates map[string]StateIfc) (*StateMxnSimpleflow, error) {
	// call constructor for StateMxnGeneric
	smg, err := NewStateMxnGeneric(smxName, transitionsMap, precreatedStates)
	smsf := &StateMxnSimpleflow{
		StateMxnGeneric: smg,
	}
	if err != nil {
		return smsf, err
	}

	return smsf, err
}

// This function will automatically progress through the states, until it reaches a final state or an error occurs
func (smsf *StateMxnSimpleflow) ChangeToInitialStateAndAutoprogressToOtherStates(initialstateName string) error {
	hasOkNokTransitionsFunc := func(stateName string) (hasOkNokTransitions bool, OkStatename string, NokStatename string) {
		tMap := smsf.GetTransitionsMap()
		if len(tMap[stateName]) < 2 {
			return false, "", ""
		}
		OkStatename = tMap[stateName][0]
		NokStatename = tMap[stateName][len(tMap[stateName])-1]
		return true, OkStatename, NokStatename
	}

	a_state := initialstateName
	for {
		hasOkNokTransitions, OkStatename, NokStatename := hasOkNokTransitionsFunc(a_state)
		serr := smsf.StateMxnGeneric.Change(a_state)
		// NOTE: serr may come from handler-error or another-error. We assume its a handler-error without additional checks
		if serr == nil {
			if !hasOkNokTransitions {
				// this state has no transitions defined, so this is a final-state
				// return smfs.GetError(), which may be nil or might be an error set in a previous state
				return smsf.GetError()
			}
			a_state = OkStatename
		} else {
			// serr is not nil, Change returned error

			// Save error into smx-data
			smsf.setError(serr)
			if !hasOkNokTransitions {
				// this state has no transitions defined, so this is a final-state
				// return smfs.GetError(), which in this case will be serr just stored before
				return smsf.GetError()
			}
			a_state = NokStatename
		}
	}
}

func (smf *StateMxnSimpleflow) Change(stateName string) error {
	return fmt.Errorf("Change() method is not allowed for StateMxnSimpleflow. Use ChangeToInitialStateAndAutoprogressToOtherStates() instead")
}
