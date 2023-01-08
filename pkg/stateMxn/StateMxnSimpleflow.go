package stateMxn

/*
StateMxnSimpleflow

This is a special implementation of a statemxn, where each state should have 2 trasitions:

	[0] nextStateOk
	[-1] FailedNok

and each state should be precreated with some execution-handler, to do some processing.

The state machine will take each state, execute its handlers with state.Activate() and
  - if there is no error, it will change to nextStateOk which assumed to be sm.transitionsMap[state.GetName()][0], and progress from there
  - if there is an error, it will change to FailedNok which is assumed to be sm.transitionsMap[state.GetName()][-1]

So overall, this statemachine should take some precreated states, each with 2 transitions and 1-or-more-handlers, and will automatically
progress the execution from state to state, until it reaches a finalstate, or an error occurs.

There is no smsf.Change() method, because the state machine will start executing the initial state when NewStateMxnSimpleFlow() is called, and
then will automatically progress through the states, until it reaches a final state, or an error occurs.
*/
type StateMxnSimpleflow struct {
	smg *StateMxnGeneric
}

func NewStateMxnSimpleFlow(transitionsMap map[string][]string, initialStateName string, precreatedStates map[string]*State) (*StateMxnSimpleflow, error) {
	smsf := &StateMxnSimpleflow{}

	// NewStateMxnGeneric will Change() into initialstate, and return its error
	// The initialstate should never fail, because if would not allow the state machine to be created
	smg, err := NewStateMxnGeneric(transitionsMap, initialStateName, precreatedStates)
	if err != nil {
		return nil, err
	}
	// At this point, the initialstate is was executed, and the state machine is ready to progress
	nextStateName, ok := smg.GetTransitionsMap()
	/// CONTINUE HERE
	// loop over the state, to the nextstateOk or failedStateNok

	smg.GetCurrentState().GetName()

	smsf.smg = smg
	return smsf, nil
}

func (smsf *StateMxnSimpleflow) GetCurrentState() *State {
	return smsf.smg.GetCurrentState()
}
func (smsf *StateMxnSimpleflow) GetHistoryOfStates() []*State {
	return smsf.smg.GetHistoryOfStates()
}
func (smsf *StateMxnSimpleflow) GetTransitionsMap() map[string][]string {
	return smsf.smg.GetTransitionsMap()
}
func (smsf *StateMxnSimpleflow) Is(stateName string) bool {
	return smsf.smg.Is(stateName)
}
