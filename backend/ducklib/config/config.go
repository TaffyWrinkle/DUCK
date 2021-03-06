// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package config

//TODO rename either package or struct to sth better

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"path/filepath"

	"fmt"

	"github.com/Microsoft/DUCK/backend/ducklib/structs"
)

var cfgWebDir string
var cfgJwtKey string
var cfgRulebaseDir string

func init() {
	flag.StringVar(&cfgWebDir, "webdir", "", "The root directory for serving web content")
	flag.StringVar(&cfgJwtKey, "jwtkey", "", "The secret used to sign the JWT")
	flag.StringVar(&cfgRulebaseDir, "rulebasedir", "", "The Directory to the Rulebases")
	flag.Parse()
}

//Configuration contains configuration values like the JWT Key or the RulebaseDir.
// these values can be set via configuration.json file, environment variables or flags
type Configuration struct {
	DBConfig    *structs.DBConf `json:"database,omitempty"`
	JwtKey      []byte          `json:"jwtkey,omitempty"`
	WebDir      string          `json:"webdir,omitempty"`
	RulebaseDir string          `json:"rulebasedir,omitempty"`
}

//NewConfiguration is the Constructor for a new structs.Configuration struct.
//it uses information from a cofiguration file, command flags, environment Variables and its own defaults to
//decide initial values for the configuration
func NewConfiguration(confpath string) Configuration {

	c := Configuration{}

	key, err := randomJWT()
	if err != nil {
		panic(fmt.Sprintf("Could not generate random JWT: %s", err))
	}

	//setting defaults
	c.JwtKey = key
	c.WebDir = "src/github.com/Microsoft/DUCK/frontend/dist"
	c.RulebaseDir = "src/github.com/Microsoft/DUCK/RuleBases"

	//overwrite defaults with information from config file
	if err := c.getFileConfig(confpath); err != nil {
		log.Printf("Could not load configuration file: %s", err)

	}
	//overwrite with information from environment
	c.getEnv()
	//overwrite with information from flags
	c.getFlags()

	c.setAbsPaths()
	return c
}

//If there is an env var GOPATH and RulebaseDir or WebDir are relative paths they are assumed to be ralative to the gopath
func (c *Configuration) setAbsPaths() {
	goPath := os.Getenv("GOPATH")
	if goPath != "" {
		log.Println("Found GOPATH, will use gopath for relative paths.")
		if !filepath.IsAbs(c.RulebaseDir) {
			c.RulebaseDir = filepath.Join(goPath, c.RulebaseDir)
		}

		if !filepath.IsAbs(c.WebDir) {
			c.WebDir = filepath.Join(goPath, c.WebDir)
		}
	}
}

//getFlags populates the configuration struct with data from the flags this program is called with
func (c *Configuration) getFlags() {

	if cfgJwtKey != "" {
		c.JwtKey = []byte(cfgJwtKey)
	}
	if cfgWebDir != "" {
		c.WebDir = cfgWebDir
	}
	if cfgRulebaseDir != "" {
		c.RulebaseDir = cfgRulebaseDir
	}

}

//getFileConfig reads the config from a JSON formatted config file and
// sets the read values as configuration values
func (c *Configuration) getFileConfig(confpath string) error {
	dat, err := ioutil.ReadFile(confpath)
	if err != nil {
		return err

	}
	err = json.Unmarshal(dat, &c)
	return err

}

func (c *Configuration) getEnv() {
	//Get Environment Variables

	env := os.Getenv("DUCK_JWTKEY")
	if env != "" {
		c.JwtKey = []byte(env)
	}

	env = os.Getenv("DUCK_WEBDIR")
	if env != "" {
		c.WebDir = env
	}

	env = os.Getenv("DUCK_RULEBASEDIR")
	if env != "" {
		c.RulebaseDir = env
	}
	//has to be not empty and also something like a boolean to be set
	env = os.Getenv("DUCK_DATABASE.LOCATION")
	if env != "" {
		c.DBConfig.Location = env
	}
	env = os.Getenv("DUCK_DATABASE.PORT")
	if env != "" {
		if p, err := strconv.Atoi(env); err == nil {
			c.DBConfig.Port = p
		} else {
			log.Printf("Could not read value for PORT: %s", err)
		}
	}
	env = os.Getenv("DUCK_DATABASE.USERNAME")
	if env != "" {
		c.DBConfig.Username = env
	}
	env = os.Getenv("DUCK_DATABASE.PASSWORD")
	if env != "" {
		c.DBConfig.Password = env
	}
	env = os.Getenv("DUCK_DATABASE.NAME")
	if env != "" {
		c.DBConfig.Name = env
	}
}

func randomJWT() ([]byte, error) {
	c := 28
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
