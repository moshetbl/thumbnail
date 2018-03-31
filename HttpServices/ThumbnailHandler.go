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

import (
	"net/http"
	"fmt"
	"net/url"
	"log"
	"errors"
	"strconv"
)

// implements thumnail service handler


// thumbnail service parameters
type thumbnailParameters struct {
	width int
	height int
	url string
	tumbnailTmpPath string
	tmpPath string
	sessionId int
}

// registration function
func registerThumbnail(config *CommonServiceConfig) error {
	http.HandleFunc(config.Path, thumbnailHandler)
	return nil
}

// extract parameters from URL
func fillThumbnailParams(values url.Values) (*thumbnailParameters, error){
	var err error

	params := thumbnailParameters{}

	value := values.Get("url")

	if value == "" {
		log.Print("url parameter not exists")
		return nil, errors.New("url not found")
	}

	params.url = value

	value = values.Get("width")

	if value == "" {
		log.Print("width parameter not exists")
		return nil, errors.New("width not found")
	}

	params.width, err = strconv.Atoi(value)

	if err != nil {
		log.Print("width Not valid")
		return nil, errors.New("width Not valid")
	}

	value = values.Get("height")

	if value == "" {
		log.Print("height parameter not exists")
		return nil, errors.New("height not found")
	}

	params.height, err = strconv.Atoi(value)

	if err != nil {
		log.Print("height Not valid")
		return nil, errors.New("height Not valid")
	}

	//load needed params from handler's service
	params.sessionId = gServiceManager.getSessionId()
	params.tmpPath = gServiceManager.config.TempPath

	// create file information
	fileName, err := extractFileNameFromUrl(params.url)
	if err != nil {
		return nil, err
	}

	params.tumbnailTmpPath = params.tmpPath + "/" + strconv.Itoa(params.sessionId) + fileName

	return &params, nil
}

func thumbnailImageResize(params *thumbnailParameters) error{
	return nil
}

// thumbnail service handler
func thumbnailHandler(w http.ResponseWriter, r *http.Request) {
	params, err :=fillThumbnailParams(r.URL.Query())

	if err != nil {
		http.Error(w, errorStringToJson(err.Error()), http.StatusMethodNotAllowed)
		return
	}

	// download image
	if err := downloadFile(params.url, params.tumbnailTmpPath); err != nil {
		http.Error(w, errorStringToJson(err.Error()), http.StatusNotFound)
		return
	}

	// TBD resize image
	if err := thumbnailImageResize(params); err != nil {
		http.Error(w, errorStringToJson(err.Error()), http.StatusInternalServerError)
		return
	}

	// TBD upload image to browser

	fmt.Fprintf(w, "TEST OK")
}
