package stateMxn

import (
	"regexp"
	"time"

	"github.com/mohae/deepcopy"
)

// NOTE:
//
//	    StateIfc should be a pointer-to-a-state-object
//		Ex: (*State) implements a StateIfc but (State) does not implement a StateIfc
type StateIfc interface {
	GetName() string
	GetInputs() StateInputs
	GetOutputs() StateOutputs
	GetError() error
	GetData() StateData
	AddHandlerBegin(handler StateHandler)
	AddHandlerExec(handler StateHandler)
	AddHandlerEnd(handler StateHandler)
	activate(smData StateMxnData, inputs StateInputs) (outputs StateOutputs, err error)
	Is(stateNameRegexp string) (bool, error)
	deepcopy() StateIfc
}

// read inputs, write outputs, read/write data
type StateHandler func(inputs StateInputs, outputs StateOutputs, stateData StateData, smachineData StateMxnData) error

type StateInputs map[string]interface{}

type StateOutputs map[string]interface{}

func (so *StateOutputs) Convert2Inputs() StateInputs {
	si := make(StateInputs)
	for k, v := range *so {
		si[k] = v
	}
	return si
}

type StateData map[string]interface{}

type HistoryOfStates []StateIfc

// Returns string with ordered states
func (hos HistoryOfStates) DisplayStatesFlow() string {
	var str string
	for _, state := range hos {
		str += state.GetName() + "\t[" + state.GetData()["timeElapsed"].(time.Duration).String() + "]"
		if serr := state.GetError(); serr != nil {
			str += "\t!ERROR: " + serr.Error()
		}
		str += "\n"
	}
	return str
}

/*
Begin-handlers   >   Exec-handlers   >  End-handlers
*/
type State struct {
	name string

	inputs  StateInputs  // input from previous state
	outputs StateOutputs // output to next state. See also state.GetOutputs()
	err     error        // error from any handler

	// data is a map where handlers can store any data meaningfull for the state, and
	// made publicly readable with state.GetData()
	// --- timestamps ---
	// data["timeStart"]
	// data["timeEnd"]
	// data["timeElapsed"]
	data StateData

	// handlers["begin"]
	// handlers["exec"]
	// handlers["end"]
	handlers map[string][]StateHandler
}

// inputs can be nil
func NewState(name string) *State {
	outputs := make(StateOutputs)
	data := make(StateData)
	handlers := make(map[string][]StateHandler)

	newState := &State{
		name:     name,
		outputs:  outputs,
		data:     data,
		handlers: handlers,
	}
	newState.addDefaultHandlers()
	return newState
}
func (s *State) GetName() string {
	return s.name
}
func (s *State) GetInputs() StateInputs {
	return s.inputs
}
func (s *State) GetOutputs() StateOutputs {
	return s.outputs
}
func (s *State) GetError() error {
	return s.err
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
// If there is an error in any begin-handler, it will return it and not execute the exec-handlers nor end-handlers
// If there is an error in any exec-handler, then it will still execute the end-handlers and then return the error
func (s *State) activate(smData StateMxnData, inputs StateInputs) (outputs StateOutputs, err error) {
	s.inputs = deepcopy.Copy(inputs).(StateInputs)

	// Executes all begin-handlers
	for _, handler := range s.handlers["begin"] {
		err := handler(s.inputs, s.outputs, s.data, smData)
		if err != nil {
			s.err = err
			return nil, err
		}
	}

	// Executes all exec-handlers
	var execErr error
	for _, handler := range s.handlers["exec"] {
		execErr = handler(s.inputs, s.outputs, s.data, smData)
		if execErr != nil {
			s.err = execErr
			break
		}
	}

	// Executes all end-handlers
	for _, handler := range s.handlers["end"] {
		err := handler(s.inputs, s.outputs, s.data, smData)
		if err != nil {
			s.err = err
			return nil, err
		}
	}
	if execErr != nil {
		return nil, execErr
	}

	return s.outputs, nil
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

func (s *State) deepcopy() StateIfc {
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

func (s *State) addDefaultHandlers() {
	timeStartHandlerBegin := func(inputs StateInputs, outputs StateOutputs, stateData StateData, smData StateMxnData) error {
		stateData["timeStart"] = time.Now()
		return nil
	}
	timeEndHandlerEnd := func(inputs StateInputs, outputs StateOutputs, stateData StateData, smData StateMxnData) error {
		stateData["timeEnd"] = time.Now()
		stateData["timeElapsed"] = stateData["timeEnd"].(time.Time).Sub(stateData["timeStart"].(time.Time))
		return nil
	}
	s.AddHandlerBegin(timeStartHandlerBegin)
	s.AddHandlerEnd(timeEndHandlerEnd)
}
