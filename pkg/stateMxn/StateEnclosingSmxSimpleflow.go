package stateMxn

import "fmt"

type StateEnclosingSmxSimpleflow struct {
	*State
}

// See example8
func NewStateEnclosingSmxSimpleflow(stateName string, smxInnerSf *StateMxnSimpleflow, smxInitialStateName string) *StateEnclosingSmxSimpleflow {
	se := &StateEnclosingSmxSimpleflow{
		State: NewState(stateName),
	}
	se.setEnclosedSmx(smxInnerSf)
	se.AddHandlerExec(
		func(inputs StateInputs, outputs StateOutputs, stateData StateData, smData StateMxnData) error {
			// smxInner: progress the state-changes
			fmt.Println("HHHHHHHHHEEEEEEEEEEEEEERRRRRRRRRRREEEEEEEEEEEEE")
			smxInnerSf := stateData["enclosedSmx"].(*StateMxnSimpleflow)
			err := smxInnerSf.ChangeToInitialStateAndAutoprogressToOtherStates(smxInitialStateName)
			return err
		})
	return se
}

func (se *StateEnclosingSmxSimpleflow) GetEnclosedSmx() *StateMxnSimpleflow {
	return se.GetData()["enclosedSmx"].(*StateMxnSimpleflow)
}
func (se *StateEnclosingSmxSimpleflow) setEnclosedSmx(smxsf *StateMxnSimpleflow) {
	se.GetData()["enclosedSmx"] = smxsf
}
