package server

import (
	"testing"

	"github.com/phuc0302/go-server/expected_format"
	"github.com/phuc0302/go-server/util"
)

func Test_bindHandlers(t *testing.T) {
	// [Test 1] Create new route empty string
	route := DefaultRoute("")
	route.BindHandler(Get, func(request *RequestContext) {})
	if route.handlers == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if route.handlers[Get] == nil {
			t.Error(expectedFormat.NotNil)
		}
	}

	// [Test 2] Create new route with non empty string
	path := "/example/{userID}"
	regexPattern := util.ConvertPath(path)

	route = DefaultRoute(regexPattern)
	route.BindHandler(Get, func(request *RequestContext) {})
	if route.regex.String() != "^/example/(?P<userID>[^/#?]+)/?$" {
		t.Errorf(expectedFormat.StringButFoundString, "^/example/(?P<userID>[^/#?]+)/?$", route.regex.String())
	}
	matched, params := route.Match(Get, "/example/1")
	if !matched {
		t.Errorf(expectedFormat.BoolButFoundBool, true, matched)
	}
	if len(params) != 1 {
		t.Errorf(expectedFormat.NumberButFoundNumber, 1, len(params))
	} else {
		if params["userID"] != "1" {
			t.Errorf(expectedFormat.StringButFoundString, "1", params["userID"])
		}
	}
}
func Test_bindHandlersWithPanic(t *testing.T) {
	route := DefaultRoute("")

	defer func() {
		if r := recover(); r != nil {
			/* Expected panic */
		}
	}()
	route.BindHandler(Post, nil)
	t.Errorf(expectedFormat.Panic)
}

func Test_invokeHandlers(t *testing.T) {
	route := DefaultRoute("/example/{userID}/profile/{profileID}")
	defer func() {
		if r := recover(); r != nil {
			/* Expected panic */
		}
	}()
	route.BindHandler(Get, func(request *RequestContext) {
		panic("Test if func had been invoked or not.")
	})
	route.InvokeHandlers(nil)
	t.Errorf(expectedFormat.Panic)
}

func Test_match_InvalidHTTPMethod(t *testing.T) {
	route := DefaultRoute("/example/{userID}/profile/{profileID}")
	route.BindHandler(Get, func(request *RequestContext) {})

	matched, pathParams := route.Match(Post, "/example/1")
	if matched {
		t.Errorf(expectedFormat.BoolButFoundBool, false, matched)
	}
	if pathParams != nil {
		t.Error(expectedFormat.Nil)
	}
}

func Test_match_InvalidHTTPMethodAndInvalidPath(t *testing.T) {
	route := DefaultRoute("/example/{userID}/profile/{profileID}")
	route.BindHandler(Get, func(request *RequestContext) {})

	matched, pathParams := route.Match(Get, "/example/1/profile")
	if matched {
		t.Errorf(expectedFormat.BoolButFoundBool, false, matched)
	}
	if pathParams != nil {
		t.Error(expectedFormat.Nil)
	}
}

func Test_match_ValidHTTPMethodAndValidPath(t *testing.T) {
	path := "/example/{userID}/profile/{profileID}"
	regexPattern := util.ConvertPath(path)

	route := DefaultRoute(regexPattern)
	route.BindHandler(Get, func(request *RequestContext) {})

	matched, pathParams := route.Match(Get, "/example/1/profile/1")
	if !matched {
		t.Errorf(expectedFormat.BoolButFoundBool, true, matched)
	}
	if pathParams == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if pathParams["userID"] != "1" {
			t.Errorf(expectedFormat.StringButFoundString, "1", pathParams["userID"])
		}
		if pathParams["profileID"] != "1" {
			t.Errorf(expectedFormat.StringButFoundString, "1", pathParams["profileID"])
		}
	}
}
