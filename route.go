package server

import "regexp"

// Route describes a route component.
type Route struct {
	regex    *regexp.Regexp
	handlers map[string]HandleContextFunc
}

// DefaultRoute creates new route component.
//
// @param
// - regexPattern: a raw path that had been converted to regex pattern
func DefaultRoute(regexPattern string) *Route {
	route := &Route{
		regex:    regexp.MustCompile(regexPattern),
		handlers: make(map[string]HandleContextFunc),
	}
	return route
}

// BindHandler binds handler with specific http method.
//
// @param
// - method: the HTTP method
// - handler: the callback func to handle context request
func (r *Route) BindHandler(method string, handler HandleContextFunc) {
	/* Condition validation: only accept function */
	if handler == nil {
		panic("Request handlers must not be nil.")
	}
	r.handlers[method] = handler
}

// InvokeHandlers invokes handlers.
//
// @param
// - c: the request context
func (r *Route) InvokeHandlers(c *RequestContext) {
	handler := r.handlers[c.Method]
	handler(c)
}

// Match matchs request path against route's regex pattern.
//
// @param
// - method: the HTTP method
// - path: path from url request
func (r *Route) Match(method string, path string) (bool, map[string]string) {
	if matches := r.regex.FindStringSubmatch(path); len(matches) > 0 && matches[0] == path {
		if handler := r.handlers[method]; handler != nil {

			// Find path params if there is any
			var params map[string]string
			if names := r.regex.SubexpNames(); len(names) > 1 {

				params = make(map[string]string)
				for i, name := range names {
					if len(name) > 0 {
						params[name] = matches[i]
					}
				}
			}

			// Return result
			return true, params
		}
	}
	return false, nil
}
