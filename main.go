package main

import (
	"fmt"
	"log"

	"github.com/zipizapclouds/stateMxn/pkg/stateMxn"
)

func stateMxnGeneric_example1() {
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
	smg, err := stateMxn.NewStateMxnGeneric(transitionsMap, initialStateName, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Get current state name
	currentStateName := smg.GetCurrentState().GetName()
	fmt.Println("currentStateName:", currentStateName) // "Init"

	// Change state
	err = smg.Change("Running")
	if err != nil {
		log.Fatal(err)
	}
	// Get current state name
	currentStateName = smg.GetCurrentState().GetName()
	fmt.Println("currentStateName:", currentStateName) // "Running"

	// Change state
	err = smg.Change("FinishedOk")
	if err != nil {
		log.Fatal(err)
	}
	// Get current state name
	currentStateName = smg.GetCurrentState().GetName()
	fmt.Println("currentStateName:", currentStateName) // "FinishedOk"

	// Check if current state matches any ^Finished state
	isFinished, err := smg.GetCurrentState().Is("^Finished")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("isFinished:", isFinished) // true
}

func stateMxnGeneric_example2() {
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
	smg, _ := stateMxn.NewStateMxnGeneric(transitionsMap, initialStateName, nil)
	fmt.Println("currentStateName:", smg.GetCurrentState().GetName()) // "Init_TriggerB"
	f := func() (isInit, isRunning, isFinished bool) {
		isInit, _ = smg.Is("^Init")
		isRunning, _ = smg.Is("^Running")
		isFinished, _ = smg.Is("^Finished")
		return isInit, isRunning, isFinished
	}
	isInit, isRunning, isFinished := f()
	fmt.Printf("isInit: %v, isRunning: %v, isFinished: %v\n", isInit, isRunning, isFinished) // true, false, false

	_ = smg.Change("Running_ProcessZeta")
	fmt.Println("currentStateName:", smg.GetCurrentState().GetName()) // "Running_ProcessZeta"
	isInit, isRunning, isFinished = f()
	fmt.Printf("isInit: %v, isRunning: %v, isFinished: %v\n", isInit, isRunning, isFinished) // false, true, false

	_ = smg.Change("Running_ProcessTau")
	fmt.Println("currentStateName:", smg.GetCurrentState().GetName()) // "Running_ProcessTau"
	isInit, isRunning, isFinished = f()
	fmt.Printf("isInit: %v, isRunning: %v, isFinished: %v\n", isInit, isRunning, isFinished) // false, true, false

	_ = smg.Change("FinishedOk")
	fmt.Println("currentStateName:", smg.GetCurrentState().GetName()) // "FinishedOk"
	isInit, isRunning, isFinished = f()
	fmt.Printf("isInit: %v, isRunning: %v, isFinished: %v\n", isInit, isRunning, isFinished) // false, false, true
}

func stateMxnGeneric_example3() {
	// ===== States and Tansitions =====
	// Init
	//   |--> Running
	//           |------------> FinishedOk
	//		     |------------> FinishedNok
	//   |--------------------> FinishedNok
	fmt.Println("===== stateMxnGeneric_example3 =====")

	// Create a StateMxnGeneric
	transitionsMap := map[string][]string{
		"Init":    {"Running", "FinishedNok"},
		"Running": {"FinishedOk", "FinishedNok"},
	}
	initialStateName := "Init"

	// When the statemachine changes to a state containing handlers, the handlers are called
	// So lets pre-create the states objects to add handlers to them
	runningState := stateMxn.NewState("Running")
	runningState.AddHandlerBegin(
		func(inputs stateMxn.StateInputs, outputs stateMxn.StateOutputs, data stateMxn.StateData) error {
			fmt.Println("+ inside runningState handlerBegin")
			return nil
		})
	runningState.AddHandlerExec(
		func(inputs stateMxn.StateInputs, outputs stateMxn.StateOutputs, data stateMxn.StateData) error {
			fmt.Println("+ inside runningState handlerExec")
			return nil
		})
	runningState.AddHandlerEnd(
		func(inputs stateMxn.StateInputs, outputs stateMxn.StateOutputs, data stateMxn.StateData) error {
			fmt.Println("+ inside runningState handlerEnd")
			return nil
		})

	precreatedStates := map[string]*stateMxn.State{
		"Running": runningState,
	}
	smg, err := stateMxn.NewStateMxnGeneric(transitionsMap, initialStateName, precreatedStates)
	if err != nil {
		log.Fatal(err)
	}
	// Get current state name
	currentStateName := smg.GetCurrentState().GetName()
	fmt.Println("currentStateName:", currentStateName) // "Init"

	// Change state
	err = smg.Change("Running")
	if err != nil {
		log.Fatal(err)
	}
	// Get current state name
	currentStateName = smg.GetCurrentState().GetName()
	fmt.Println("currentStateName:", currentStateName) // "Running"

	// Change state
	err = smg.Change("FinishedOk")
	if err != nil {
		log.Fatal(err)
	}
	// Get current state name
	currentStateName = smg.GetCurrentState().GetName()
	fmt.Println("currentStateName:", currentStateName) // "FinishedOk"

	// Check if current state matches any ^Finished state
	isFinished, err := smg.GetCurrentState().Is("^Finished")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("isFinished:", isFinished) // true
}

func main() {
	stateMxnGeneric_example1()
	stateMxnGeneric_example2()
	stateMxnGeneric_example3()

}
