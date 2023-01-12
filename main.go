package main

import (
	"fmt"
	"log"
	"time"

	"github.com/zipizapclouds/stateMxn/pkg/stateMxn"
)

func logFatalIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func stateMxnGeneric_example1() {
	// Basic usage of StateMxnGeneric
	//
	// ===== States and Tansitions =====
	// Init
	//   |--> Running
	//           |------------> FinishedOk
	//		     |------------> FinishedNok
	//   |--------------------> FinishedNok
	fmt.Println("===== stateMxnGeneric_example1 =====")

	// Create a StateMxnGeneric
	transitionsMap := map[string][]string{
		"Init":    {"Running", "FinishedNok"},
		"Running": {"FinishedOk", "FinishedNok"},
	}
	initialStateName := "Init"
	smg, err := stateMxn.NewStateMxnGeneric(transitionsMap, nil)
	logFatalIfError(err)

	// Start by changing to the initial state
	err = smg.Change(initialStateName)
	logFatalIfError(err)

	// Get current state name
	currentStateName := smg.GetCurrentState().GetName()
	fmt.Println("currentStateName:", currentStateName) // "Init"

	// Change state
	err = smg.Change("Running")
	logFatalIfError(err)
	// Get current state name
	currentStateName = smg.GetCurrentState().GetName()
	fmt.Println("currentStateName:", currentStateName) // "Running"

	// Change state
	err = smg.Change("FinishedOk")
	logFatalIfError(err)
	// Get current state name
	currentStateName = smg.GetCurrentState().GetName()
	fmt.Println("currentStateName:", currentStateName) // "FinishedOk"

	// Check if current state matches any ^Finished state
	isFinished, err := smg.GetCurrentState().Is("^Finished")
	logFatalIfError(err)
	fmt.Println("isFinished:", isFinished) // true
}

func stateMxnGeneric_example2() {
	// More complex state machine, with multiple states and transitions
	//
	// ===== States and Tansitions =====
	// Init_TriggerA -------+-> Running_ProcessZeta -----> Running_ProcessTau  ---------> FinishedOk
	//  | Init_TriggerB ---/      |                          |
	//	|  |                      |                          |
	//	|   \                      \                          \
	//	 \---+----------------------+--------------------------+------------------------> FinishedNok
	fmt.Println("===== stateMxnGeneric_example2 =====")

	// Create a StateMxnGeneric
	transitionsMap := map[string][]string{
		"Init_TriggerA": {
			"Running_ProcessZeta",
			"FinishedNok"},
		"Init_TriggerB": {
			"Running_ProcessZeta",
			"FinishedNok"},
		"Running_ProcessZeta": {
			"Running_ProcessTau",
			"FinishedNok"},
		"Running_ProcessTau": {
			"FinishedOk",
			"FinishedNok"},
	}
	initialStateName := "Init_TriggerB"
	smg, err := stateMxn.NewStateMxnGeneric(transitionsMap, nil)
	logFatalIfError(err)
	err = smg.Change(initialStateName)
	logFatalIfError(err)

	fmt.Println("currentStateName:", smg.GetCurrentState().GetName()) // "Init_TriggerB"
	f := func() (isInit, isRunning, isFinished bool) {
		isInit, _ = smg.Is("^Init")
		isRunning, _ = smg.Is("^Running")
		isFinished, _ = smg.Is("^Finished")
		return isInit, isRunning, isFinished
	}
	isInit, isRunning, isFinished := f()
	fmt.Printf("isInit: %v, isRunning: %v, isFinished: %v\n", isInit, isRunning, isFinished) // true, false, false

	err = smg.Change("Running_ProcessZeta")
	logFatalIfError(err)
	fmt.Println("currentStateName:", smg.GetCurrentState().GetName()) // "Running_ProcessZeta"
	isInit, isRunning, isFinished = f()
	fmt.Printf("isInit: %v, isRunning: %v, isFinished: %v\n", isInit, isRunning, isFinished) // false, true, false

	err = smg.Change("Running_ProcessTau")
	logFatalIfError(err)
	fmt.Println("currentStateName:", smg.GetCurrentState().GetName()) // "Running_ProcessTau"
	isInit, isRunning, isFinished = f()
	fmt.Printf("isInit: %v, isRunning: %v, isFinished: %v\n", isInit, isRunning, isFinished) // false, true, false

	err = smg.Change("FinishedOk")
	logFatalIfError(err)
	fmt.Println("currentStateName:", smg.GetCurrentState().GetName()) // "FinishedOk"
	isInit, isRunning, isFinished = f()
	fmt.Printf("isInit: %v, isRunning: %v, isFinished: %v\n", isInit, isRunning, isFinished) // false, false, true
}

func stateMxnGeneric_example3() {
	// State machine, with state-handlers
	//
	// ===== States and Tansitions =====
	// Init
	//   |--> Running
	//           |------------> FinishedOk
	//		     |------------> FinishedNok
	//   |--------------------> FinishedNok
	fmt.Println("===== stateMxnGeneric_example3 =====")

	transitionsMap := map[string][]string{
		"Init":    {"Running", "FinishedNok"},
		"Running": {"FinishedOk", "FinishedNok"},
	}
	initialStateName := "Init"

	// When the statemachine changes to a state containing handlers, the handlers are called
	// To pass the handlers of states, we need to pre-create the states objects beforehand and add handlers to them
	runningState := stateMxn.NewState("Running")
	{
		runningState.AddHandlerBegin(
			func(inputs stateMxn.StateInputs, outputs stateMxn.StateOutputs, stateData stateMxn.StateData, smData stateMxn.StateMxnData) error {
				fmt.Println("+ inside runningState handlerBegin")
				return nil
			})
		runningState.AddHandlerExec(
			func(inputs stateMxn.StateInputs, outputs stateMxn.StateOutputs, stateData stateMxn.StateData, smData stateMxn.StateMxnData) error {
				fmt.Println("+ inside runningState handlerExec")
				// sleep 500 miliseconds
				time.Sleep(500 * time.Millisecond)
				return nil
			})
		runningState.AddHandlerEnd(
			func(inputs stateMxn.StateInputs, outputs stateMxn.StateOutputs, stateData stateMxn.StateData, smData stateMxn.StateMxnData) error {
				fmt.Println("+ inside runningState handlerEnd")
				return nil
			})
	}
	precreatedStates := map[string]stateMxn.StateIfc{
		"Running": runningState,
	}

	// Now lets create the statemachine passing the precreated states
	smg, err := stateMxn.NewStateMxnGeneric(transitionsMap, precreatedStates)
	logFatalIfError(err)
	err = smg.Change(initialStateName)
	logFatalIfError(err)

	// Get current state name
	currentStateName := smg.GetCurrentState().GetName()
	fmt.Println("currentStateName:", currentStateName) // "Init"

	// Change state
	err = smg.Change("Running")
	logFatalIfError(err)
	// Get current state name
	currentStateName = smg.GetCurrentState().GetName()
	fmt.Println("currentStateName:", currentStateName) // "Running"

	// Change state
	err = smg.Change("FinishedOk")
	logFatalIfError(err)
	// Get current state name
	currentStateName = smg.GetCurrentState().GetName()
	fmt.Println("currentStateName:", currentStateName) // "FinishedOk"

	// Check if current state matches any ^Finished state
	isFinished, err := smg.GetCurrentState().Is("^Finished")
	logFatalIfError(err)
	fmt.Println("isFinished:", isFinished) // true

	fmt.Println("............................................")
	fmt.Println("timeElapsed in each state:")
	for _, state := range smg.GetHistoryOfStates() {
		fmt.Printf("  %s:\t%s\n", state.GetName(), state.GetData()["timeElapsed"])
	}
}

func stateMxnGeneric_example4() {
	// This example shows how to use StateMxnSimpleflow:
	//  - each state should have 2 (at least) destination-transition-states: "Ok" and "Nok" (except final-states)
	//  - each state can have zero-or-more handlerExec functions
	// When a there is a change to a state, that state handlers are called in order, and then
	// the transition to another state is done automatically based on the return value of the handlers of current-state:
	//  - if all the handlers return nil, the transition is done to the "Ok" state
	//  - if any of the handlers returns an error, the transition is done to the "Nok" state
	// The "Ok" state is considered the stateMxn.GetTransitionMap()[currentStateName][0] state
	// The "Nok" state is considered the stateMxn.GetTransitionMap()[currentStateName][-1] state (where -1 represents the index of last-element)
	//
	// Notice that this example if very similar to example3, but:
	// - uses StateMxnSimpleflow instead of StateMxnGeneric
	// - the states-changes progress automatically, that is done by the state-machine itself (without requiring external calls to sm.Change())

	//
	// ===== States and Tansitions =====
	// Init  --> RunningAlpha --> RunningBeta --> FinishedOk
	// |----------------------------------------> FinishedNok
	//		     |------------------------------> FinishedNok
	//		                      |-------------> FinishedNok
	//
	//

	fmt.Println("===== stateMxnGeneric_example4 =====")

	transitionsMap := map[string][]string{
		"Init":         {"RunningAlpha", "FinishedNok"}, // first-transition is "Ok" transition, second-transition is "Nok" transition
		"RunningAlpha": {"RunningBeta", "FinishedNok"},  // first-transition is "Ok" transition, second-transition is "Nok" transition
		"RunningBeta":  {"FinishedOk", "FinishedNok"},   // first-transition is "Ok" transition, second-transition is "Nok" transition
	}
	initialStateName := "Init"

	// When the statemachine changes to a state containing handlers, the handlers are called
	// To pass the handlers of states, we need to pre-create the states objects beforehand and add handlers to them
	// In the StateMxnSimpleflow, if the handlers return nil then the transition is done to the "Ok" state, otherwise the transition is done to the "Nok" state
	runningAlphaState := stateMxn.NewState("RunningAlpha")
	runningAlphaState.AddHandlerExec(
		func(inputs stateMxn.StateInputs, outputs stateMxn.StateOutputs, stateData stateMxn.StateData, smData stateMxn.StateMxnData) error {
			// do alpha processing...
			outputs["fromAlpha"] = "This is the output from alpha"
			return nil
		})
	runningBetaState := stateMxn.NewState("RunningBeta")
	runningBetaState.AddHandlerExec(
		func(inputs stateMxn.StateInputs, outputs stateMxn.StateOutputs, stateData stateMxn.StateData, smData stateMxn.StateMxnData) error {
			// do beta processing...
			fmt.Println("inputs[\"fromAlpha\"]:", inputs["fromAlpha"])
			return fmt.Errorf("error created by handlerExec of RunningBeta to force transition to FinishedNok")
		})

	// Precreated states
	precreatedStates := map[string]stateMxn.StateIfc{
		runningAlphaState.GetName(): runningAlphaState,
		runningBetaState.GetName():  runningBetaState,
	}

	// Now lets create the statemachine passing the precreated states
	smsf, err := stateMxn.NewStateMxnSimpleFlow(transitionsMap, precreatedStates)
	logFatalIfError(err)
	err = smsf.ChangeToInitialStateAndAutoprogressToOtherStates(initialStateName)
	logFatalIfError(err)

	// At this point, the statemachine has already autoprogressed until a final state (FinishedOk or FinishedNok) was reached
	// Get final state name, which is the current state name
	currentStateName := smsf.GetCurrentState().GetName()
	finalStateName := currentStateName
	fmt.Println("finalStateName:", finalStateName) // "FinishedNok"

	fmt.Println("............................................")
	fmt.Println(smsf.GetHistoryOfStates().DisplayStatesFlow())
}

func main() {
	stateMxnGeneric_example1()
	stateMxnGeneric_example2()
	stateMxnGeneric_example3()
	stateMxnGeneric_example4()

}
