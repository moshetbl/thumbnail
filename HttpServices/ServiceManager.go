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

package HttpServices

// services manager implementation
import (
	"sync"
	"errors"
	"log"
	"io/ioutil"
	"github.com/go-yaml/yaml"
	"net/http"
	"os"
)

// global pointer to service manager
var gServiceManager *serviceManager = nil

// service manager structure / needed data
type serviceManager struct {
	sessionId int // unique session id counter
	mutex sync.Mutex // session id guard
	config ServiceManagerConfig // configuration
	servicesRegistration map[string] registerService // registration function map
}

// create new manager. only if not exists
func newManager() error {
	if gServiceManager != nil {
		return errors.New("Manager already exists")
	}

	gServiceManager = &serviceManager{}
	gServiceManager.sessionId = 0
	gServiceManager.servicesRegistration = make(map[string] registerService)
	gServiceManager.config.Services = make(map[string]CommonServiceConfig)
	return nil
}

// create and init manager and run
func Init(filePath string) error {
	if err := newManager(); err != nil {
		return err
	}

	if err := gServiceManager.Init(filePath); err != nil {
		return err
	}

	log.Println("Service Manager is going up!!!")
	if err := gServiceManager.Start(); err != nil {
		return err
	}

	return nil
}


// init manager
func (p *serviceManager) Init(filePath string) error {

	p.fillRegistration()

	if err := p.loadConfiguration(filePath); err != nil {
		return err
	}

	if err := p.registerServices(); err != nil {
		return err
	}

	registerImageFormats() // register supported image formats for services

	return nil
}

// start services
func (p *serviceManager) Start() error {
	port := os.Getenv("PORT")
	
	if port == "" {
		log.Fatal("$PORT not set")
        }
	
	log.Printf("**** Listen on Port:%s *****\n", port)
	
	return http.ListenAndServe(":" + port, nil)
}

func (p *serviceManager) getSessionId() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	sessionId := p.sessionId
	p.sessionId += 1
	return sessionId
}

// fill all available supported services
func (p *serviceManager) fillRegistration() {
	p.servicesRegistration["thumbnail"] = registerThumbnail
	// TBD add same lines for each service
}

// register services in configuration
func (p *serviceManager) registerServices() error {

	// search if
	for serviceKey, serviceConfig := range p.config.Services {
		if serviceReg, ok := p.servicesRegistration[serviceKey]; ok {
			if err := serviceReg(&serviceConfig); err!= nil { // register single service
				log.Printf("Service:%s Registration Error %s", serviceKey, err.Error())
				return err
			}

			log.Printf("Service:%s was registered", serviceKey)
		}
	}
	return nil
}

// Load configuration file
func (p *serviceManager) loadConfiguration(filePath string) error {
	// load configuration file
	b, err := ioutil.ReadFile(filePath) // just pass the file name
	if err != nil {
		log.Println(err)
		return err
	}

	// unmarshal yaml, to configuration struct
	if err = yaml.Unmarshal(b, &p.config); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

