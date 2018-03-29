/*
Copyright 2018 Moshe Tubul

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/moshetbl/go/Common"
	"github.com/moshetbl/go/HttpHandlers"
	"log"
	"os"
	"errors"
)

var configFile = "/home/moshe/go/projects/src/github.com/moshetbl/go/Config/config.yaml"

func getConfigFile() (string, error) {

	if len(os.Args) == 1 || len(os.Args) > 2 {
		log.Fatal("Wrong number of arguments")
		return "", errors.New("Wrong number of arguments")
	}

	fileName := os.Args[1]

	if _, err := os.Stat(fileName); err != nil {
		log.Fatal("Configuration file state error")
		return "", errors.New("Configuration file state error")
	}

	return fileName, nil
}

func main(){

	// get file name
	path, err := getConfigFile();

	if err != nil {
		log.Fatal(err)
		return
	}

	// load configuration file:
	conf, err := Common.LoadConfiguration(path)

	if err != nil{
		log.Printf("Failed to load configuration file:%s", configFile)
	}

	// start service
	HttpHandlers.LoadHandlersService(conf)
}
