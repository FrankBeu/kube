package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_runInDebugMode(t *testing.T) {
	t.Run("should return immediately if not in debug mode", func(t *testing.T) {
		debugMode := false
		timeoutDuration := time.Nanosecond

		assert.Eventually(
			t,
			func() bool { _ = RunInDebugMode(debugMode, timeoutDuration, false); return true },
			10*time.Millisecond,
			5*time.Millisecond,
			"not in debugMode: RunInDebugMode blocked",
		)
	})
	t.Run("should return an TimeoutError if in debugMode and debugReady is not set", func(t *testing.T) {
		debugMode := true
		timeoutDuration := time.Nanosecond

		err := RunInDebugMode(debugMode, timeoutDuration, false)

		errorWanted := DebugingTimeoutError{}
		assert.EqualError(t,
			err,
			errorWanted.Error(),
			"in debugMode: did not receive an TimeoutError",
		)
	})
	t.Run("should return without error, if debugReady is set to true", func(t *testing.T) {
		//// TODO: test ist working
		//// However, debugReady is set to true before the debugLoop has started
		//// channels cannot be used in the dap-interface
		debugMode := true
		timeoutDuration := time.Nanosecond

		debugSetupCompleted := make(chan struct{})
		debugSetupError := make(chan error)

		debugReady := false
		debugReadyPtr := &debugReady

		go func() {
			err := RunInDebugMode(debugMode, timeoutDuration, debugReady)
			if err != nil {
				debugSetupError <- err
			}
			close(debugSetupCompleted)
		}()
		*debugReadyPtr = true

		var err error
		select {
		case <-debugSetupCompleted:
		case err = <-debugSetupError:
			t.Fatal("in debugMode; debugReady set, but: ", err)
		}
	})
}
