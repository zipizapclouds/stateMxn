package main

import (
	"fmt"
	"log"

	"github.com/zipizapclouds/stateMxn/pkg/stateMxn"
)

func check_stateMxnGeneric() {
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

}

func main() {
	check_stateMxnGeneric()
}
