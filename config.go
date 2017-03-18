package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/johntdyer/slackrus"
	"github.com/phuc0302/go-server/util"
)

// Config file's name.
const (
	Debug   = "server.debug.cfg"
	Release = "server.release.cfg"
)

// HTTP Methods.
const (
	Copy    = "copy"
	Delete  = "delete"
	Get     = "get"
	Head    = "head"
	Link    = "link"
	Options = "options"
	Patch   = "patch"
	Post    = "post"
	Purge   = "purge"
	Put     = "put"
	Unlink  = "unlink"
)

// Config describes a configuration object that will be used during application life time.
type Config struct {
	// Server
	Host string `json:"host"`
	Port int    `json:"port"`

	// Header
	HeaderSize    int           `json:"header_size"`    // In KB
	MultipartSize int64         `json:"multipart_size"` // In MB
	ReadTimeout   time.Duration `json:"timeout_read"`   // In seconds
	WriteTimeout  time.Duration `json:"timeout_write"`  // In seconds

	// Log
	LogLevel     string `json:"log_level"`
	SlackURL     string `json:"slack_url"`
	SlackIcon    string `json:"slack_icon"`
	SlackUser    string `json:"slack_user"`
	SlackChannel string `json:"slack_channel"`

	// HTTP Method
	AllowMethods  []string          `json:"allow_methods"`
	RedirectPaths map[string]string `json:"redirect_paths"`
	StaticFolders map[string]string `json:"static_folders"`

	// Extensions
	Extensions map[string]interface{} `json:"extensions,omitempty"`

	// File's path
	configPath string
}

// CreateConfig generates a default configuration file.
//
// @param
// - configFile {string} (a file's path that will be used to generate configuration file)
func CreateConfig(configFile string) {
	if configPath := util.GetEnv(util.ConfigPath); len(configPath) > 0 && !strings.HasPrefix(configFile, configPath) {
		configFile = fmt.Sprintf("%s/%s", configPath, configFile)
	}

	// Create default config
	config := &Config{
		Host:          "localhost",
		Port:          8080,
		HeaderSize:    (5 << 10),
		MultipartSize: (1 << 20),
		ReadTimeout:   (15 * time.Second),
		WriteTimeout:  (15 * time.Second),
		LogLevel:      "debug",
		SlackURL:      "",
		SlackIcon:     ":ghost:",
		SlackUser:     "Server",
		SlackChannel:  "#channel",
		AllowMethods:  []string{Copy, Delete, Get, Head, Link, Options, Patch, Post, Purge, Put, Unlink},

		RedirectPaths: map[string]string{
			"401": "/login",
		},
		StaticFolders: map[string]string{
			"/assets":    "assets",
			"/resources": "resources",
		},

		configPath: configFile,
	}

	// Create new file
	config.Save()
}

// LoadConfig will load pre-generated configuration file into memory for later used.
//
// @param
// - configFile {string} (a file's path that will be used to load pre-generated configuration file)
//
// @return
// - config {Config} (an instance of server's configuration)
func LoadConfig(configFile string) *Config {
	// Append file's path if necessary
	if configPath := util.GetEnv(util.ConfigPath); len(configPath) > 0 && !strings.HasPrefix(configFile, configPath) {
		configFile = fmt.Sprintf("%s/%s", configPath, configFile)
	}

	// Check if config file is available or not
	if !util.FileExisted(configFile) {
		CreateConfig(configFile)
	}

	file, _ := os.Open(configFile)
	defer file.Close()

	// Load config file
	var config Config
	bytes, _ := ioutil.ReadAll(file)

	if err := json.Unmarshal(bytes, &config); err != nil {
		fmt.Println("Could not load config file at: ", configFile)
		os.Exit(1)
	}
	config.configPath = configFile

	// Convert duration to seconds
	config.HeaderSize <<= 10
	config.MultipartSize <<= 20
	config.ReadTimeout *= time.Second
	config.WriteTimeout *= time.Second

	// Define redirectPaths
	redirectPaths = make(map[int]string, len(config.RedirectPaths))
	for s, path := range config.RedirectPaths {
		if status, err := strconv.Atoi(s); err == nil {
			redirectPaths[status] = path
		}
	}

	// Setup logger
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		level = logrus.DebugLevel
	}
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(level)

	// Setup slack notification if neccessary
	if len(config.SlackURL) > 0 {
		logrus.AddHook(&slackrus.SlackrusHook{
			HookURL:        config.SlackURL,
			Channel:        config.SlackChannel,
			Username:       config.SlackUser,
			IconEmoji:      config.SlackIcon,
			AcceptedLevels: slackrus.LevelThreshold(level),
		})
	}

	// Return config's instance
	return &config
}

// Save will create new configuration file, override if necessary.
func (c *Config) Save() {
	if util.FileExisted(c.configPath) {
		os.Remove(c.configPath)
	}

	// Revert changed
	c.HeaderSize >>= 10
	c.MultipartSize >>= 20
	c.ReadTimeout /= time.Second
	c.WriteTimeout /= time.Second

	// Create new file
	file, _ := os.Create(c.configPath)
	defer file.Close()

	configJSON, _ := json.MarshalIndent(c, "", "  ")
	file.Write(configJSON)

	// Revert changed
	c.HeaderSize <<= 10
	c.MultipartSize <<= 20
	c.ReadTimeout *= time.Second
	c.WriteTimeout *= time.Second
}

// GetExtension returns extension data that had been associated with input key.
//
// @param
// - key {string} (an input key to retrieve associated extension data)
//
// @return
// - value {interface} (a generic data that had been associated with input key or null)
func (c *Config) GetExtension(key string) interface{} {
	if c.Extensions == nil {
		return nil
	}
	return c.Extensions[key]
}

// SetExtension extends server's default configuration.
//
// @param
// - key {string} (an input key that will be associated with extension value)
// - value {interface} (a generic data that will be associated with input key)
func (c *Config) SetExtension(key string, value interface{}) {
	/* Condition validation: validate input */
	if len(key) == 0 || value == nil {
		return
	}

	if c.Extensions == nil {
		c.Extensions = make(map[string]interface{})
	}
	c.Extensions[key] = value
}
