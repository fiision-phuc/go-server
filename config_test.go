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
	defer os.Remove(Debug)
	CreateConfig(Debug)

	if !util.FileExisted(Debug) {
		t.Errorf("Expected %s file had been created but found nil.", Debug)
	}
}

func Test_LoadConfig(t *testing.T) {
	defer os.Remove(Debug)
	config := LoadConfig(Debug)

	// Validate basic information
	if config.Host != "localhost" {
		t.Errorf(expectedFormat.StringButFoundString, "localhost", config.Host)
	}
	if config.Port != 8080 {
		t.Errorf(expectedFormat.NumberButFoundNumber, 8080, config.Port)
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

func Test_Extensions(t *testing.T) {
	defer os.Remove(Debug)
	config := LoadConfig(Debug)

	if config.GetExtension("key") != nil {
		t.Error(expectedFormat.Nil)
	}

	if _, ok := config.GetExtension("key").(string); ok {
		t.Error(expectedFormat.Nil)
	}

	config.SetExtension("key", "value")
	if config.GetExtension("key") == nil {
		t.Error(expectedFormat.NotNil)
	} else {
		value := config.GetExtension("key")
		if v, ok := value.(string); ok {
			if v != "value" {
				t.Errorf(expectedFormat.StringButFoundString, "value", v)
			}
		} else {
			t.Errorf("Invalid value")
		}
		config.Save(Debug)
	}
}
