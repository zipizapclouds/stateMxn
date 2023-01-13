package stateMxn

import (
	"fmt"
	"strconv"
	"time"

	"github.com/trislu/plantuml"
)

type plantUmlGenOpts struct {
	stripHeaderFooter bool
}

// opts can be nil
func plantUmlGen(smx StateMxnIfc, opts *plantUmlGenOpts) (text string, diagramUrl string) {
	if opts == nil {
		opts = &plantUmlGenOpts{}
	}

	// define text
	{
		/*
			https://pkg.go.dev/github.com/trislu/plantuml#section-readme

			@startuml
			scale 3/4

			[*] -> State1
			note as N1
				asasd
				na aasd asdasd
			end note
			State1 --> State2 : Succeeded
			State1 --> [*] : Aborted
			State2 --> State3 : Succeeded
			State2 --> [*] : Aborted
			state State3 {
				state "Accumulate Enough Data\nLong State Name" as long1
				long1 : Just a test
				[*] --> long1
				long1 --> long1 : New Data
				long1 --> ProcessData : Enough Data
			}
			State3 --> State3 : Failed
			State3 --> [*] : Succeeded / Save Result
			State3 --> [*] : Aborted

			@enduml
		*/

		// define header and footer
		var header, footer string
		{
			header = `
@startuml
scale 5/8
skinparam sequenceMessageAlign left
skinparam sequenceReferenceAlign left

`
			footer = "\n\n@enduml\n"

			if opts.stripHeaderFooter {
				header = ""
				footer = ""
			}
		}
		// define body
		var body string
		{
			type specialKeysType map[string]func(k string, v interface{}, mapName string) string
			mapStringInterfaceFormatter := func(m map[string]interface{}, mapName string, specialKeys specialKeysType, eol string) (str string) {
				str = ""
				for k, v := range m {
					// if k is in specialCases, use the specialCases[k] function
					if f, ok := specialKeys[k]; ok {
						str += f(k, v, mapName)
						continue
					}

					// switch branches to handle different types of the v variable
					// There is one branch matching all the basic types
					// and one branch for all the other types
					switch v.(type) {
					case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128, bool, time.Duration:
						// basic types and some other specific types
						// will have its value showned with %v
						str += mapName + "[" + k + "]: " + fmt.Sprintf("%v", v) + eol
					default:
						// other types will have the type showned with %T
						str += mapName + "[" + k + "]: " + fmt.Sprintf("(%T)", v) + eol
					}
				}
				return str
			}
			inputFormatter := func(o StateInputs) string {
				return mapStringInterfaceFormatter(o, "inputs", make(specialKeysType), `\n`)
			}
			outputFormatter := func(o StateOutputs) string {
				return mapStringInterfaceFormatter(o, "outputs", make(specialKeysType), `\n`)
			}
			sdataFormatter := func(d StateData) string {
				str := mapStringInterfaceFormatter(
					d,
					"state.data",
					specialKeysType(map[string]func(k string, v interface{}, mapName string) string{
						"enclosedSmx": func(k string, v interface{}, mapName string) string {
							return mapName + "[" + k + "]: " + fmt.Sprintf("%s (%T)", v.(StateMxnIfc).GetName(), v) + `\n`
						},
						"timeEnd":   func(k string, v interface{}, mapName string) string { return "" },
						"timeStart": func(k string, v interface{}, mapName string) string { return "" },
					}),
					`\n`,
				)
				return str
			}
			smxdataFormatter := func(d StateMxnData) string {
				str := mapStringInterfaceFormatter(
					d,
					"smx.data",
					specialKeysType(map[string]func(k string, v interface{}, mapName string) string{}),
					"\n",
				)
				return str
			}

			initialStateName := smx.GetName() + "_0" + smx.GetHistoryOfStates()[0].GetName()
			initialStateName = replace2alphanum(initialStateName)
			initialStateInputs := smx.GetHistoryOfStates()[0].GetInputs()
			body = "[*] --> " + initialStateName
			{
				if iinputsTxt := inputFormatter(initialStateInputs); len(iinputsTxt) > 0 {
					body += " : " + iinputsTxt
				}
				body += "\n"
			}
			// add smx.data as floating note, if its not empty
			{
				smxData := smx.GetData()
				if len(smxData) > 0 {
					body += "note as " + smx.GetName() + "\n" + identLinesInString("  ", smxdataFormatter(smxData)) + "\nend note\n"
				}
			}
			for i := 1; i <= len(smx.GetHistoryOfStates()); i++ {
				prevStateName := smx.GetName() + "_" + strconv.Itoa(i-1) + "" + smx.GetHistoryOfStates()[i-1].GetName()
				prevStateName = replace2alphanum(prevStateName)
				prevStateData := smx.GetHistoryOfStates()[i-1].GetData()
				prevStateOutputs := smx.GetHistoryOfStates()[i-1].GetOutputs()
				prevStateOutputsStr := outputFormatter(prevStateOutputs)
				var prevStateErr string
				{
					prevStateErr = ""
					if smx.GetHistoryOfStates()[i-1].GetError() != nil {
						prevStateErr = `\nERROR ` + smx.GetHistoryOfStates()[i-1].GetError().Error()
					}
				}
				var nextStateName string
				{
					if i == len(smx.GetHistoryOfStates()) {
						nextStateName = "[*]"
					} else {
						nextStateName = smx.GetName() + "_" + strconv.Itoa(i) + smx.GetHistoryOfStates()[i].GetName()
						nextStateName = replace2alphanum(nextStateName)
					}
				}
				// prevStateName --> nextStateName : prevStateOutputsStr + prevStateErr \n
				{
					body += prevStateName + " --> " + nextStateName
					if len(prevStateOutputsStr+prevStateErr) > 0 {
						body += " : " + prevStateOutputsStr + prevStateErr
					}
					body += "\n"
				}
				// prevStateName : prevStateData \n
				{
					body += prevStateName
					if len(sdataFormatter(prevStateData)) > 0 {
						body += " : " + sdataFormatter(prevStateData)
					}
					body += "\n"
				}
				if eSmx, ok := prevStateData["enclosedSmx"]; ok {
					eSmx := eSmx.(StateMxnIfc)
					eSmxText, _ := plantUmlGen(eSmx, &plantUmlGenOpts{stripHeaderFooter: true})
					body += "state " + prevStateName + " ##[bold]green {\n" + identLinesInString("    ", eSmxText) + "\n}\n"
				}

			}
		}
		text = header + body + footer
	}

	// define diagramUrl
	diagramUrl = `http://www.plantuml.com/plantuml/img/` + plantuml.Encode(text)

	return text, diagramUrl
}

func plantUmlGen4TransitionsMap(transitionsMap map[string][]string) (text string, diagramUrl string) {
	var header, footer string
	{
		header = `
@startuml
scale 5/8
skinparam sequenceMessageAlign left
skinparam sequenceReferenceAlign left
`
		footer = "\n@enduml\n"
	}

	var body string
	{
		body = ""
		for fromState, toStates := range transitionsMap {
			for _, toState := range toStates {
				body += fromState + " -[dotted]-> " + toState + "\n"
			}
		}
	}
	text = header + body + footer

	// define diagramUrl
	diagramUrl = `http://www.plantuml.com/plantuml/img/` + plantuml.Encode(text)

	return text, diagramUrl
}
