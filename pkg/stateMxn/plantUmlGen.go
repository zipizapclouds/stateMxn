package stateMxn

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/trislu/plantuml"
)

type plantUmlGenOpts struct {
	stripHeaderFooter bool
}

// opts can be nil
func plantUmlGen(smx *StateMxnGeneric, opts *plantUmlGenOpts) (text string, diagramUrl string) {
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
			inputFormatter := func(o StateInputs) string {
				return fmt.Sprintf("%v", o)
			}
			outputFormatter := func(o StateOutputs) string {
				return fmt.Sprintf("%v", o)
			}
			dataFormatter := func(d StateData) string {
				// strWithEol := spew.Sdump(d)
				// strWithEolEscaped := regexp.MustCompile("\n").ReplaceAllString(strWithEol, `\n`)
				// return strWithEolEscaped
				// >> Seems PlantUml cannot process all that

				str := ""
				for k, v := range d {
					if k == "timeEnd" || k == "timeStart" {
						continue
					}
					str += "state.data[" + k + "]: " + fmt.Sprintf("%+v", v) + `\n`
				}
				return str
			}
			replace2alphanum := func(s string) string {
				return regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(s, "_")
			}

			initialStateName := smx.GetName() + "_0_" + smx.historyOfStates[0].GetName()
			initialStateName = replace2alphanum(initialStateName)
			initialStateInputs := smx.historyOfStates[0].GetInputs()
			body = "[*] --> " + initialStateName + " : " + inputFormatter(initialStateInputs) + "\n"
			for i := 1; i <= len(smx.historyOfStates); i++ {
				prevStateName := smx.GetName() + "_" + strconv.Itoa(i-1) + "_" + smx.historyOfStates[i-1].GetName()
				prevStateName = replace2alphanum(prevStateName)
				prevStateData := smx.historyOfStates[i-1].GetData()
				prevStateOutputs := smx.historyOfStates[i-1].GetOutputs()
				prevStateOutputsStr := outputFormatter(prevStateOutputs)
				var prevStateErr string
				{
					prevStateErr = ""
					if smx.historyOfStates[i-1].GetError() != nil {
						prevStateErr = `\nERROR ` + smx.historyOfStates[i-1].GetError().Error()
					}
				}
				var nextStateName string
				{
					if i == len(smx.historyOfStates) {
						nextStateName = "[*]"
					} else {
						nextStateName = smx.GetName() + "_" + strconv.Itoa(i) + "_" + smx.historyOfStates[i].GetName()
						nextStateName = replace2alphanum(nextStateName)
					}
				}
				body += prevStateName + " --> " + nextStateName + " : " + prevStateOutputsStr + prevStateErr + "\n"
				body += prevStateName + " : " + dataFormatter(prevStateData) + "\n"
				if eSmx, ok := prevStateData["enclosedSmx"]; ok {
					eSmx := eSmx.(*StateMxnGeneric)
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
