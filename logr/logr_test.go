package logr

import (
	"bytes"
	"testing"
)

func TestReInit(t *testing.T) {
	// Set up test cases
	testCases := []struct {
		name                string
		isVerboseRequested  bool
		expectedTraceOutput string
	}{
		{
			name:                "verbose is requested",
			isVerboseRequested:  true,
			expectedTraceOutput: "trace message\n",
		},
		{
			name:                "verbose is not requested",
			isVerboseRequested:  false,
			expectedTraceOutput: "",
		},
	}

	// SetTarget output to buffer for testing
	var buf bytes.Buffer
	SetTarget(&buf)

	for _, tc := range testCases {

		// Run the test case
		t.Run(tc.name, func(t *testing.T) {

			SetVerbosity(tc.isVerboseRequested)

			// Test Trace logger output
			buf.Reset()
			Trace.Println("trace message")
			actualTraceOutput := buf.String()

			if actualTraceOutput != tc.expectedTraceOutput {
				t.Errorf("Trace output was %q, expected %q", actualTraceOutput, tc.expectedTraceOutput)
			}
		})
	}
}
