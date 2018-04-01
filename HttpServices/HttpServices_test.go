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
	"testing"
	"net/url"
)

func TestLoadConfiguration(t *testing.T) {
	if err := newManager(); err != nil {
		t.Error("Cannot create manager")
	}

	// valid case
	if err:= gServiceManager.loadConfiguration("testFiles/config_valid.yaml"); err != nil {
		t.Error("Configuration file not opened")
	}

	// file does not exists
	if err:= gServiceManager.loadConfiguration("testFiles/config11.yaml"); err == nil {
		t.Error("Configuration file not opened")
	}

	// file is not valid
	if err:= gServiceManager.loadConfiguration("testFiles/config_notvalid.yaml"); err == nil {
		t.Error("Node valid configuration file, should not be opened")
	}
}

func TestServiceManagerInit(t *testing.T) {
	gServiceManager = nil
	if err := newManager(); err != nil {
		t.Error("Cannot create manager")
	}

	// double creation test
	if err := newManager(); err == nil {
		t.Error("Manager should not be created")
	}

	if err := gServiceManager.Init("testFiles/config_valid.yaml"); err != nil {
		t.Error("Manager Init error")
	}

	if err := gServiceManager.Init("testFiles/config_xxxx.yaml"); err == nil {
		t.Error("Manager should be inited")
	}

	if err := gServiceManager.Init("testFiles/config_notvalid.yaml"); err == nil {
		t.Error("Manager should be inited")
	}

	if err := gServiceManager.Init(""); err == nil {
		t.Error("Manager should be inited")
	}
}

func TestFillThumbnailParams(t *testing.T) {
	values := make(url.Values)
	values["url"] = append(values["url"],"http://www.example.com/image.jpg")
	values["width"] = append(values["width"],"100")
	values["height"] = append(values["height"],"200")

	// valid case
	params, err := fillThumbnailParams(values)
	if err != nil {
		t.Error("value are valid, should be parsed")
	}

	if params.url != "http://www.example.com/image.jpg" || params.height != 200 || params.width != 100 {
		t.Error("values are not as expected")
	}

	// not valid
	values["width"][0] = "www"
	params, err = fillThumbnailParams(values)
	if err == nil {
		t.Error("Type is not numeric, function should return error")
	}

	values["width"][0] = "100"
	values["height"][0] = "www"
	params, err = fillThumbnailParams(values)
	if err == nil {
		t.Error("Type is not numeric, function should return error")
	}

	values["width"][0] = "100"
	values["height"][0] = "100"
	values["url"][0] = ""
	params, err = fillThumbnailParams(values)
	if err == nil {
		t.Error("value empty, function should return error")
	}

	values["url"][0] = "http://www.example.com/image.jpg"
	delete(values,"width")
	params, err = fillThumbnailParams(values)
	if err == nil {
		t.Error("value empty, function should return error")
	}

	delete(values,"height")
	params, err = fillThumbnailParams(values)
	if err == nil {
		t.Error("value empty, function should return error")
	}

	delete(values,"url")
	params, err = fillThumbnailParams(values)
	if err == nil {
		t.Error("value empty, function should return error")
	}
}

func TestCommonExtractFileNameFromUrl(t *testing.T) {

	// valid cases
	res, err := extractFileNameFromUrl("http://www.example.com/image.jpg")
	if err != nil {
		t.Error("Valid url file name should be parsed")
	}
	if res != "image.jpg" {
		t.Error("image string not as expected")
	}

	res, err = extractFileNameFromUrl("http://www.example.com/_image.jpg")
	if err != nil {
		t.Error("Valid url file name should be parsed")
	}
	if res != "_image.jpg" {
		t.Error("image string not as expected")
	}

	res, err = extractFileNameFromUrl("http://www.example.com/image.jpeg")
	if err != nil {
		t.Error("Valid url file name should be parsed")
	}
	if res != "image.jpeg" {
		t.Error("image string not as expected")
	}

	res, err = extractFileNameFromUrl("http://www.example.com/_.jpeg")
	if err != nil {
		t.Error("Valid url file name should be parsed")
	}
	if res != "_.jpeg" {
		t.Error("image string not as expected")
	}

	res, err = extractFileNameFromUrl("http://www.example.com/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.jpeg")
	if err != nil {
		t.Error("Valid url file name should be parsed")
	}
	if res != "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.jpeg" {
		t.Error("image string not as expected")
	}

	res, err = extractFileNameFromUrl("http://www.example.com/file.file.jpeg")
	if err != nil {
		t.Error("Valid url file name should be parsed")
	}
	if res != "file.file.jpeg" {
		t.Error("image string not as expected")
	}

	// none valid cases
	res, err = extractFileNameFromUrl("http://www.example.com/jpeg")
	if err == nil {
		t.Error("file name not valid should not be parsed")
	}

	res, err = extractFileNameFromUrl("http://www.example.com/")
	if err == nil {
		t.Error("file name not valid should not be parsed")
	}

	res, err = extractFileNameFromUrl("")
	if err == nil {
		t.Error("url not valid should not be parsed")
	}

	res, err = extractFileNameFromUrl("http://www.example.com/image.bmp")
	if err == nil {
		t.Error("none valid file type should not be parsed")
	}
}

func TestCommonIsImageFileTypeValid(t *testing.T) {
	// valid cases
	if res := isImageFileTypeValid("jpeg"); res != true {
		t.Error("valid type should be parsed")
	}

	if res := isImageFileTypeValid("jpg"); res != true {
		t.Error("valid type should be parsed")
	}

	if res := isImageFileTypeValid("JPEG"); res != true {
		t.Error("valid type should be parsed")
	}

	if res := isImageFileTypeValid("JPG"); res != true {
		t.Error("valid type should be parsed")
	}

	// none valid cases
	if res := isImageFileTypeValid("jpeg_"); res == true {
		t.Error("not valid type should not be parsed")
	}

	if res := isImageFileTypeValid("jpg_"); res == true {
		t.Error("not valid type should not be parsed")
	}

	if res := isImageFileTypeValid("JPEG_"); res == true {
		t.Error("not valid type should not be parsed")
	}

	if res := isImageFileTypeValid("JPG_"); res == true {
		t.Error("not valid type should not be parsed")
	}

	if res := isImageFileTypeValid("_jpeg"); res == true {
		t.Error("not valid type should not be parsed")
	}

	if res := isImageFileTypeValid("_jpg"); res == true {
		t.Error("not valid type should not be parsed")
	}

	if res := isImageFileTypeValid("_JPEG"); res == true {
		t.Error("not valid type should not be parsed")
	}

	if res := isImageFileTypeValid("_JPG"); res == true {
		t.Error("not valid type should not be parsed")
	}

	if res := isImageFileTypeValid(""); res == true {
		t.Error("not valid type should not be parsed")
	}

	if res := isImageFileTypeValid("bmp"); res == true {
		t.Error("not valid type should not be parsed")
	}

	if res := isImageFileTypeValid("BMP"); res == true {
		t.Error("not valid type should not be parsed")
	}
}



