package lib

import (
	"fmt"
	"time"
)

type DebugingTimeoutError struct{}

func (t *DebugingTimeoutError) Error() string {
	return "reached DebuggingTimeout"
}

// RunInDebugMode enables attaching delve to a currently running pulumi-process
// 1. !!! set BreakPoint inside RunInDebugMode !!! otherwise you will get an 'LSP :: There is no stopped thread?"
// 2. `p{,S,PROD} up -c debugMode=true`
// 3. start delve (attach to executable) with pulumi-main's pID
// 4. inside your debug session set debugReady=true
//
// !!! -c debug=true is sticky:
// subsequent calls will be executed with debugMode==true until explicitly set debugMode=false
func RunInDebugMode(debugMode bool, timeoutDuration time.Duration, debugReady bool) error {
	//// attaching does neither work with channels nor with struct-fields or with pointers
	//// a blocking loop is needed so that a breakpoint can be catched
	//// TODO: after switch routines in delve: one loop, one to control debugReady
	if debugMode {
		t := time.Now()
		timeoutStamp := t.Add(timeoutDuration)

		fmt.Println("Use delve to set `debugReady=true`")
		for !debugReady {
			time.Sleep(time.Second)

			if t = time.Now(); t.After(timeoutStamp) {
				return &DebugingTimeoutError{}
			}
		}
	}
	return nil
}
