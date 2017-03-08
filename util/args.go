package util

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

// ReadArgs reads arguments from console input.
func ReadArgs() (sandboxMode bool, tlsMode bool, keyFile string, certFile string, shouldContinue bool) {
	shouldContinue = true
	sandboxMode = true
	tlsMode = false
	keyFile = ""
	certFile = ""

	args := os.Args[1:]
	l := len(args)

	if l > 0 {
		for i := 0; i < l; i += 2 {
			switch args[i] {

			case "--help":
				shouldContinue = false
				info := map[string]string{
					"--sandboxMode": "[true|false]",
					"--tlsMode":     "[true|false]",
					"--keyFile":     "path to server's private key file.",
					"--certFile":    "path to server's X.509 certificate.",
				}

				var buffer bytes.Buffer
				for k, v := range info {
					buffer.WriteString(fmt.Sprintf("%-15s%s\n", k, v))
				}
				fmt.Println(buffer.String())

			case "--sandboxMode":
				if flag, err := strconv.ParseBool(args[i+1]); err == nil {
					sandboxMode = flag
				}

			case "--tlsMode":
				if flag, err := strconv.ParseBool(args[i+1]); err == nil {
					tlsMode = flag
				}

			case "--certFile":
				if FileExisted(args[i+1]) {
					certFile = args[i+1]
				}

			case "--keyFile":
				if FileExisted(args[i+1]) {
					keyFile = args[i+1]
				}

			default:
				break
			}
		}
	}
	return
}
