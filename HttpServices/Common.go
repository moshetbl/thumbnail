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

// common tools for services

import (
	"fmt"
	"strings"

	"errors"
	"net/http"
	"os"
	"io"
	"image"
	"image/jpeg"
)

// service registration function type
type registerService func(*CommonServiceConfig) error

// common services configuration
type CommonServiceConfig struct {
	Path string `yaml:"path"`
}

// service manager configuration
type ServiceManagerConfig struct {
	Port string `yaml:"port"`
	TempPath string `yaml:"tmppath"`
	Services map[string]CommonServiceConfig `yaml:"services"`
}

// convert error to json
func errorStringToJson(str string) string{
	return fmt.Sprintf("{\"error\": \"%s\"}", str)
}

// extract and validate file name from url
func extractFileNameFromUrl(url string) (string, error) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	tokens = strings.Split(fileName, ".")
	fileType := tokens[len(tokens)-1]
	if isImageFileTypeValid(fileType) == false {
		return "", errors.New("File type not supported or not valid")
	}
	return fileName, nil
}

// is file type valid/supported
func isImageFileTypeValid(fileType string) bool {
	for _, fType := range []string{"jpeg", "jpg"} {
		if fileType == fType {
			return  true
		}
	}
	return false
}

func registerImageFormats() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("jpeg", "jpg", jpeg.Decode, jpeg.DecodeConfig)
}

// download file and save it
func downloadFile(url string, filePath string) error {

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// save to local file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}