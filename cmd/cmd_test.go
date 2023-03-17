package cmd

import (
	"testing"
)

func TestAppCmdWithNoParams(t *testing.T) {
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to execute root cmd with no params which should show usage %s", err)
	} else {
		t.Logf("Successfully executed root cmd without params which should show usage")
	}
}
