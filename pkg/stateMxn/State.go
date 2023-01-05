package stateMxn

import (
	"regexp"

	"github.com/mohae/deepcopy"
)

// read inputs, write outputs, read/write data
type StateHandler func(inputs StateInputs, outputs StateOutputs, data StateData) error
type StateInputs map[string]interface{}
type StateOutputs map[string]interface{}
type StateData map[string]interface{}

/*
Begin-handlers   >   Exec-handlers   >  End-handlers
*/
type State struct {
	name string

	inputs  StateInputs
	outputs StateOutputs
	data    StateData

	// handlers["begin"]
	// handlers["exec"]
	// handlers["end"]
	handlers map[string][]StateHandler
}

// inputs can be nil
func NewState(name string) *State {
	inputs := make(StateInputs)
	outputs := make(StateOutputs)
	data := make(StateData)
	handlers := make(map[string][]StateHandler)

	return &State{
		name:     name,
		inputs:   inputs,
		outputs:  outputs,
		data:     data,
		handlers: handlers,
	}
}

func (s *State) GetName() string {
	return s.name
}
func (s *State) SetInputs(inputs StateInputs) {
	s.inputs = inputs
}
func (s *State) GetOutputs() StateOutputs {
	return s.outputs
}
func (s *State) GetData() StateData {
	return s.data
}

// Appends a handler to the list of begin-handlers
func (s *State) AddHandlerBegin(handler StateHandler) {
	s.handlers["begin"] = append(s.handlers["begin"], handler)
}

// Appends a handler to the list of exec-handlers
func (s *State) AddHandlerExec(handler StateHandler) {
	s.handlers["exec"] = append(s.handlers["exec"], handler)
}

// Preprends a handler to the list of end-handlers
func (s *State) AddHandlerEnd(handler StateHandler) {
	s.handlers["end"] = append([]StateHandler{handler}, s.handlers["end"]...)
}

// Executes all handlers in the order: begin-handlers, exec-handlers, end-handlers
func (s *State) Activate() error {
	// Executes all begin-handlers
	for _, handler := range s.handlers["begin"] {
		err := handler(s.inputs, s.outputs, s.data)
		if err != nil {
			return err
		}
	}

	// Executes all exec-handlers
	for _, handler := range s.handlers["exec"] {
		err := handler(s.inputs, s.outputs, s.data)
		if err != nil {
			return err
		}
	}

	// Executes all end-handlers
	for _, handler := range s.handlers["end"] {
		err := handler(s.inputs, s.outputs, s.data)
		if err != nil {
			return err
		}
	}

	return nil
}

// Is returns true if the state name matches the given regexp
//
// stateNameRegexp - is a regexp RE2 as described at https://golang.org/s/re2syntax
// the same as used by regexp.MatchString()
//
// Examples:
//   - regexp "Init" 						matches Name "Init"
//   - regexp "Init|Running" 				matches Name "Init"
//   - regexp "Init|Running|FinishedOk" 	matches Name "Init"
//   - regexp "Finished" 					matches Name "FinisedOk" or "FinishedNok"
func (s *State) Is(stateNameRegexp string) (bool, error) {
	return regexp.MatchString(stateNameRegexp, s.name)
}

func (s *State) deepcopy() *State {
	// NOTE: deepcopy libs like https://github.com/barkimedes/go-deepcopy or https://github.com/mohae/deepcopy
	//       dont copy unexported fields - so we need to define our own deepcopy() method
	// 	     for the type, in the package where the type is defined

	// all truct fields, both exported and unexported, need to be copied here
	stateCopy := &State{
		name:     s.name,
		inputs:   deepcopy.Copy(s.inputs).(StateInputs),
		outputs:  deepcopy.Copy(s.outputs).(StateOutputs),
		data:     deepcopy.Copy(s.data).(StateData),
		handlers: deepcopy.Copy(s.handlers).(map[string][]StateHandler),
	}
	return stateCopy
}
