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

package Common

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
)

// services configuration structure
type ServiceConfig struct {
	Port string `yaml:"port"`
	Path string `yaml:"path"`

}

// Load service configuration file
func LoadConfiguration(filePath string) (*ServiceConfig, error) {
	// load configuration file
	b, err := ioutil.ReadFile(filePath) // just pass the file name
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// unmarshal yaml, to configuration struct
	conf := ServiceConfig{}
	if err = yaml.Unmarshal(b, &conf); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &conf, nil
}
