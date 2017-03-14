package util

import (
	"os"
	"testing"

	"github.com/phuc0302/go-server/expected_format"
)

func Test_ReadArgs(t *testing.T) {
	os.Args = os.Args[1:]
	os.Args = append(os.Args, "--port")
	os.Args = append(os.Args, "8080")
	os.Args = append(os.Args, "--configPath")
	os.Args = append(os.Args, "Users/phuc/Workspaces///")
	os.Args = append(os.Args, "--sandboxMode")

	if _, _, shouldContinue := ReadArgs(); !shouldContinue {
		t.Error("Expected to be able to handle old input.")
	} else {
		if GetEnv(Port) != "8080" {
			t.Errorf(expectedFormat.StringButFoundString, "8080", GetEnv(Port))
		}

		if GetEnv(ConfigPath) != "/Users/phuc/Workspaces" {
			t.Errorf(expectedFormat.StringButFoundString, "/Users/phuc/Workspaces", GetEnv(ConfigPath))
		}
	}
}
