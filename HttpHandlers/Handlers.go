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

package HttpHandlers

import (
	"github.com/moshetbl/go/Common"
	"log"
	"net/http"
	"fmt"
)

// create small error json
func errorStringToJson(str string) string{
	return fmt.Sprintf("{\"error\": \"%s\"}", str)
}
// register all handlers
func registerHandlers(config *Common.ServiceConfig) {
	http.HandleFunc(config.Path, thumbnailHandler)
	// TBD register services
}

// load service
func LoadHandlersService(config *Common.ServiceConfig) {
	registerHandlers(config)

	log.Fatal(http.ListenAndServe(":" + config.Port, nil))
}
