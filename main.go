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
	currentStateName := smg.CurrentState().Name()
	fmt.Println("currentStateName:", currentStateName) // "Init"

	// Change state
	err = smg.Change("Running")
	if err != nil {
		log.Fatal(err)
	}
	// Get current state name
	currentStateName = smg.CurrentState().Name()
	fmt.Println("currentStateName:", currentStateName) // "Running"

	// Change state
	err = smg.Change("FinishedOk")
	if err != nil {
		log.Fatal(err)
	}
	// Get current state name
	currentStateName = smg.CurrentState().Name()
	fmt.Println("currentStateName:", currentStateName) // "FinishedOk"

	// Check if current state matches any ^Finished state
	isFinished, err := smg.CurrentState().Is("^Finished")
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
	fmt.Println("currentStateName:", smg.CurrentState().Name()) // "Init_TriggerB"
	f := func() (isInit, isRunning, isFinished bool) {
		isInit, _ = smg.Is("^Init")
		isRunning, _ = smg.Is("^Running")
		isFinished, _ = smg.Is("^Finished")
		return isInit, isRunning, isFinished
	}
	isInit, isRunning, isFinished := f()
	fmt.Printf("isInit: %v, isRunning: %v, isFinished: %v\n", isInit, isRunning, isFinished) // true, false, false

	_ = smg.Change("Running_ProcessZeta")
	fmt.Println("currentStateName:", smg.CurrentState().Name()) // "Running_ProcessZeta"
	isInit, isRunning, isFinished = f()
	fmt.Printf("isInit: %v, isRunning: %v, isFinished: %v\n", isInit, isRunning, isFinished) // false, true, false

	_ = smg.Change("Running_ProcessTau")
	fmt.Println("currentStateName:", smg.CurrentState().Name()) // "Running_ProcessTau"
	isInit, isRunning, isFinished = f()
	fmt.Printf("isInit: %v, isRunning: %v, isFinished: %v\n", isInit, isRunning, isFinished) // false, true, false

	_ = smg.Change("FinishedOk")
	fmt.Println("currentStateName:", smg.CurrentState().Name()) // "FinishedOk"
	isInit, isRunning, isFinished = f()
	fmt.Printf("isInit: %v, isRunning: %v, isFinished: %v\n", isInit, isRunning, isFinished) // false, false, true
}

func main() {
	stateMxnGeneric_example2()
}
