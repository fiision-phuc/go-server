package server

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/phuc0302/go-server/expected_format"
	"github.com/phuc0302/go-server/util"
)

func Test_CreateConfig(t *testing.T) {
	defer os.Remove(debug)
	CreateConfig(debug)

	if !util.FileExisted(debug) {
		t.Errorf("Expected %s file had been created but found nil.", debug)
	}
}

func Test_LoadConfig(t *testing.T) {
	defer os.Remove(debug)
	config := LoadConfig(debug)

	// Validate basic information
	if config.Host != "localhost" {
		t.Errorf(expectedFormat.StringButFoundString, "localhost", config.Host)
	}
	if config.Port != 8080 {
		t.Errorf(expectedFormat.NumberButFoundNumber, 8080, config.Port)
	}
	if config.TLSPort != 8443 {
		t.Errorf(expectedFormat.NumberButFoundNumber, 8443, config.TLSPort)
	}
	if config.HeaderSize != 5120 {
		t.Errorf(expectedFormat.NumberButFoundNumber, 5120, config.HeaderSize)
	}
	if config.ReadTimeout != 15*time.Second {
		t.Errorf(expectedFormat.NumberButFoundNumber, 15*time.Second, config.ReadTimeout)
	}
	if config.WriteTimeout != 15*time.Second {
		t.Errorf(expectedFormat.NumberButFoundNumber, 15*time.Second, config.WriteTimeout)
	}

	// Validate allow methods
	allowMethods := []string{Copy, Delete, Get, Head, Link, Options, Patch, Post, Purge, Put, Unlink}
	if !reflect.DeepEqual(allowMethods, config.AllowMethods) {
		t.Errorf(expectedFormat.StringButFoundString, allowMethods, config.AllowMethods)
	}
	if methodsValidation == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		if !methodsValidation.MatchString(Copy) {
			t.Errorf(expectedFormat.BoolButFoundBool, true, methodsValidation.MatchString(Copy))
		}
	}

	// Validate redirect paths
	if redirectPaths == nil || len(redirectPaths) != 1 {
		t.Error(expectedFormat.NotNil)
	}
	if redirectPaths[401] != "/login" {
		t.Errorf(expectedFormat.StringButFoundString, "/login", redirectPaths[401])
	}

	// Validate static folders
	staticFolders := map[string]string{
		"/assets":    "assets",
		"/resources": "resources",
	}
	if !reflect.DeepEqual(staticFolders, config.StaticFolders) {
		t.Errorf(expectedFormat.StringButFoundString, staticFolders, config.StaticFolders)
	}
}
