package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/phuc0302/go-server/util"
)

// recovery recovers server from panic state.
func recovery(w http.ResponseWriter, r *http.Request) {
	if err := recover(); err != nil {
		var status *util.Status
		if httpError, ok := err.(*util.Status); ok {
			status = httpError
		} else {
			status = util.Status500()
		}

		// Return error
		if redirectURL := redirectPaths[status.Code]; len(redirectURL) > 0 {
			http.Redirect(w, r, redirectURL, status.Code)
		} else {
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(status.Code)

			cause, _ := json.Marshal(status)
			w.Write(cause)
		}

		// Slack log
		go func() {
			// Generate error report
			var buffer bytes.Buffer
			buffer.WriteString(fmt.Sprintf("[%s][%d] %s\n", time.Now().UTC().Format(time.RFC822), status.Code, status.Description))
			buffer.WriteString(fmt.Sprintf("%s %s %s\n\n", r.Proto, r.Method, r.URL.Path))
			buffer.WriteString(fmt.Sprintf("%s: %s\n", "address", r.RemoteAddr))
			buffer.WriteString(fmt.Sprintf("%s: %s\n\n", "user-agent", r.UserAgent()))

			// Write header
			buffer.WriteString(fmt.Sprintf("%s: %s\n", "referer", r.Referer()))
			buffer.WriteString(fmt.Sprintf("%s:\n", "header"))

			for header, value := range r.Header {
				header = strings.ToLower(header)

				if header == "user-agent" || header == "referer" {
					continue
				}
				buffer.WriteString(fmt.Sprintf("- %s: %s\n", header, value))
			}

			//			// Write Path Params
			//			if c.PathParams != nil && len(c.PathParams) > 0 {
			//				buffer.WriteString("\n")
			//				idx = 0
			//				for key, value := range c.PathParams {
			//					if idx == 0 {
			//						buffer.WriteString(fmt.Sprintf("%-12s: %s = %s\n", "Path Params", key, value))
			//					} else {
			//						buffer.WriteString(fmt.Sprintf("%-12s: %s = %s\n", "", key, value))
			//					}
			//					idx++
			//				}
			//			}

			//			// Write Query Params
			//			if c.QueryParams != nil && len(c.QueryParams) > 0 {
			//				buffer.WriteString("\n")
			//				idx = 0
			//				for key, value := range c.QueryParams {
			//					if idx == 0 {
			//						buffer.WriteString(fmt.Sprintf("%-12s: %s = %s\n", "Query Params", key, value))
			//					} else {
			//						buffer.WriteString(fmt.Sprintf("%-12s: %s = %s\n", "", key, value))
			//					}
			//					idx++
			//				}
			//			}

			//			// Write stack trace
			//			buffer.WriteString("\nStack Trace:\n")
			//			callStack(3, &buffer)

			// Log error
			logrus.Warningln(buffer.String())
		}()
	}
}
