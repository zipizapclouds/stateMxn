package stateMxn

// type StateEnclosingSmx struct {
// 	*State
// 	// data["enclosedSmx"] *StateMxnGeneric
// }

// func NewStateEnclosingSmx(stateName string, smx *StateMxnGeneric) *StateEnclosingSmx {
// 	se := &StateEnclosingSmx{
// 		State: NewState(stateName),
// 	}
// 	se.data["enclosedSmx"] = smx

// 	// Add handler to execute enclosedSmx
// 	se.AddHandlerExec(
// 		func(inputs StateInputs, outputs StateOutputs, stateData StateData, smachineData StateMxnData) error {
// 			enclosedSmx := stateData["enclosedSmx"].(*StateMxnGeneric)
// 			_, err := enclosedSmx.Execute(smachineData, inputs)
// 			return
// 		},
// 	)

// 	return se
// }

// func (se *StateEnclosingSmx) GetEnclosedSmx() *StateMxnGeneric {
// 	return se.data["enclosedSmx"].(*StateMxnGeneric)
// }
