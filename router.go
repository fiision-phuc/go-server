package server

import (
	"bytes"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/phuc0302/go-server/util"
)

// Router describes a router component implementation.
type Router struct {
	groups []string
	routes []*Route
}

// GroupRoute generates path's prefix for following URLs.
//
// @param
// - prefixURI {string} (the prefix for url)
// - handler {HandleGroupFunc} (the callback func)
func (r *Router) GroupRoute(prefixURI string, handler HandleGroupFunc) {
	r.groups = append(r.groups, prefixURI)
	handler()
	r.groups = r.groups[:len(r.groups)-1]
}

// BindRoute binds a patternURL with handler.
//
// @param
// - method {string} (HTTP request method)
// - patternURL {string} (the URL matching pattern)
// - handler {HandleContextFunc} (the callback func)
func (r *Router) BindRoute(method string, patternURL string, handler HandleContextFunc) {
	patternURL = r.mergeGroup(patternURL)
	logrus.Infof("%-6s -> %s", strings.ToUpper(method), patternURL)

	// Define regex pattern
	regexPattern := util.ConvertPath(patternURL)

	// Look for existing one before create new
	for _, route := range r.routes {
		if route.regex.String() == regexPattern {
			route.BindHandler(method, handler)
			return
		}
	}
	newRoute := DefaultRoute(regexPattern)
	newRoute.BindHandler(method, handler)

	// Append to current list
	r.routes = append(r.routes, newRoute)
}

// MatchRoute matches a route with a pathURL.
//
// @param
// - method {string} (HTTP request method)
// - pathURL {string} (request's path that will be matched)
//
// @return
// - route {Route} (a route that lead to request's handler, might be null if it is not yet defined)
// - pathParams {map[string]string} (a path params, might be null if there is no route)
func (r *Router) MatchRoute(method string, pathURL string) (*Route, map[string]string) {
	for _, route := range r.routes {
		if ok, pathParams := route.Match(method, pathURL); ok {
			return route, pathParams
		}
	}
	return nil, nil
}

// mergeGroup merges multiple prefixURIs into single prefixURI.
//
// @param
// - patternURL {string} (the URL matching pattern)
//
// @return
// - patternURL {string} (the URL matching pattern)
func (r *Router) mergeGroup(patternURL string) string {
	if len(r.groups) > 0 {
		var buffer bytes.Buffer
		for _, prefixURI := range r.groups {
			buffer.WriteString(prefixURI)
		}

		if len(patternURL) > 0 {
			buffer.WriteString(patternURL)
		}
		patternURL = buffer.String()
	}
	return httprouter.CleanPath(patternURL)
}
