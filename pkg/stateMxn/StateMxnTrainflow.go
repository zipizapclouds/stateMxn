package stateMxn

import "fmt"

/*
StateMxnTrainflow

The StateMxnTrainsflow is a refinement over StateMxnSimpleflow, where:
- more user-friendly: the user easily specifies an ordered train-of-ministates, each composed of name and function-handler
- the initial state is automatically set to the first indicated state, and while each state returns no error, the state machine will automatically
  progress to the next state, until the end when its set to a final state "FinishedOk".
  If any state returns an error, the state machine will jump to a "FinishedNok" state.
  The states "FinishedOk" and "FinishedNok" are auto-created and should not be defined by user

The overall idea is to make it easy easy easy, for the user to define the train-of-mninistates, and then auto-progress to "FinishedO"/"FinishedNok"

Examples: see main.go example9
*/
type StateMxnTrainflow struct {
	*StateMxnSimpleflow
	trainOfMinistates []TrainMinistate
}

type TrainMinistate struct {
	StateName   string
	HandlerFunc StateHandler
}

func NewStateMxnTrainFlow(smxName string, trainOfMinistates []TrainMinistate) (*StateMxnTrainflow, error) {

	// create smtf
	var smtf *StateMxnTrainflow
	{
		var transitionsMap map[string][]string
		{
			/*
				transitionsMap := map[string][]string{
					"Init":    {"Running", "FinishedNok"},
					"Running": {"FinishedOk", "FinishedNok"},
				}
			*/

			var statesNames []string
			{
				for _, a_ministate := range trainOfMinistates {
					a_stateName := a_ministate.StateName
					statesNames = append(statesNames, a_stateName)
				}
			}

			for i, i_stateName := range statesNames {
				curState := i_stateName
				nextState := ""
				{
					if i < len(statesNames)-1 {
						nextState = statesNames[i+1]
					} else {
						nextState = "FinishedOk"
					}
				}
				transitionsMap[curState] = []string{nextState, "FinishedNok"}
			}
		}

		var precreatedStates map[string]StateIfc
		{
			/*
				var smxInnerRunningState *stateMxn.State
				{
					smxInnerRunningState = stateMxn.NewState("Running")
					smxInnerRunningState.AddHandlerExec(func(inputs stateMxn.StateInputs, outputs stateMxn.StateOutputs, stateData stateMxn.StateData, smachineData stateMxn.StateMxnData) error {
						return fmt.Errorf("simulating error from smxInnerRunningState, to change state to FinishedNok")
						// return nil
					})
				}
				precreatedStates := map[string]stateMxn.StateIfc{
					smxInnerRunningState.GetName(): smxInnerRunningState,
				}
			*/
			for _, a_ministate := range trainOfMinistates {
				a_stateName := a_ministate.StateName
				a_stateHandler := a_ministate.HandlerFunc
				a_state := NewState(a_stateName)
				a_state.AddHandlerExec(a_stateHandler)
				precreatedStates[a_stateName] = a_state
			}
		}

		smsf, err := NewStateMxnSimpleFlow(smxName, transitionsMap, precreatedStates)
		smtf := &StateMxnTrainflow{
			StateMxnSimpleflow: smsf,
			trainOfMinistates:  trainOfMinistates,
		}
		if err != nil {
			return smtf, err
		}
	} // ATP: smtf is created and ready to be used
	return smtf, nil
}

func (smtf *StateMxnTrainflow) ChangeToInitialStateAndAutoprogressToOtherStates() error {
	return smtf.StateMxnSimpleflow.ChangeToInitialStateAndAutoprogressToOtherStates(smtf.trainOfMinistates[0].StateName)
}

func (smtf *StateMxnTrainflow) Change(stateName string) error {
	return fmt.Errorf("Change() method is not allowed for StateMxnTrainflow. Use ChangeToInitialStateAndAutoprogressToOtherStates() instead")
}
