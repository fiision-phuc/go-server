package util

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// ReadArgs reads arguments from console input.
//
// @return
// - tlsMode {bool} (enable SSL mode or not)
// - sandboxMode {bool} (enable sandbox mode or not)
// - shouldContinue {bool} (indicate that application should continue or not)
func ReadArgs() (sandboxMode bool, tlsMode bool, shouldContinue bool) {
	shouldContinue = true
	sandboxMode = true
	tlsMode = false

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
					"--port":        "Port's number that server will listen on.",
					"--configPath":  "path to server's configuration file.",
					"--sslPath":     "path to server's X.509 certificate & private key.",
				}

				var buffer bytes.Buffer
				for k, v := range info {
					buffer.WriteString(fmt.Sprintf("%-15s%s\n", k, v))
				}
				fmt.Println(buffer.String())
				os.Exit(0)

			case "--sandboxMode":
				if i+1 < l {
					if flag, err := strconv.ParseBool(args[i+1]); err == nil {
						sandboxMode = flag
					}
				}

			case "--tlsMode":
				if i+1 < l {
					if flag, err := strconv.ParseBool(args[i+1]); err == nil {
						tlsMode = flag
					}
				}

			case "--port":
				if i+1 < l {
					if number, err := strconv.Atoi(args[i+1]); err == nil {
						SetEnv(Port, fmt.Sprintf("%d", number))
					}
				}

			case "--configPath":
				if i+1 < l {
					dirPath := formatPath(args[i+1])
					if DirExisted(dirPath) {
						SetEnv(ConfigPath, dirPath)
					}
				}

			case "--sslPath":
				if i+1 < l {
					dirPath := formatPath(args[i+1])
					if DirExisted(dirPath) {
						SetEnv(SSLPath, dirPath)
					}
				}

			default:
				break
			}
		}
	}
	return
}

// formatPath formats input path and remove trail slash if there is any.
//
// @param
// - input {string} (input directory path)
//
// @return
// - output {string} (path to directory as expected)
func formatPath(input string) (output string) {
	output = httprouter.CleanPath(input)

	if output[(len(output)-1):] == "/" {
		output = output[:(len(output) - 1)]
	}
	return
}
