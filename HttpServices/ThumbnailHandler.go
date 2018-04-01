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
	"net/url"
	"log"
	"errors"
	"strconv"
	"image"
	"github.com/disintegration/imaging"
	"image/color"
	"os"
	"io"
)

// implements thumnail service handler


// thumbnail service parameters
type thumbnailParameters struct {
	width int // width of the new image
	height int // height of the new image
	url string // url of the image, used for downloading the image
	tumbnailTmpPath string // full path of the image
	fileName string // only the file name
	tmpPath string // temporary path in which the files are saved
	sessionId int // current session id
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
	params.fileName = fileName

	return &params, nil
}

func thumbnailImageResize(params *thumbnailParameters) error{
	// decode image
	srcImg, err := imaging.Open(params.tumbnailTmpPath)
	if err != nil {
		log.Println("Decode Error file: ", params.tumbnailTmpPath)
		return errors.New("Decode Error file: " + params.tumbnailTmpPath)
	}

	// calculate size of original image
	b:= srcImg.Bounds()
	origHeight := b.Max.Y
	origWidth := b.Max.X
	origRatio := float64(origWidth) / float64(origHeight)
	dstRatio := float64(params.width) / float64(params.height)
	var dstWidth, dstHeight int

	// in case aspect ratio is the same (rounded values)
	if int(origRatio * 3.0) == int(dstRatio * 3.0) {
		if origWidth < params.width {
			dstHeight = origHeight
			dstWidth = origWidth
		} else {
			dstHeight = params.height
			dstWidth = params.width
		}
	} else { // aspect ratio is different
		// pad left/right
		if dstRatio > origRatio {
			dstHeight = params.height
			dstWidth = int(float64(params.height) * origRatio)

		} else {
			// pad top/bottom
			dstWidth = params.width
			dstHeight = int(float64(params.width) / origRatio)
		}
	}

	// create background black image
	dstFinalImg := imaging.New(params.width, params.height, color.NRGBA{0, 0, 0, 0})
	resizedImg := imaging.Resize(srcImg, dstWidth, dstHeight, imaging.Lanczos)

	// merge images
	dstFinalImg = imaging.Paste(dstFinalImg, resizedImg, image.Pt((params.width - dstWidth)/2 , (params.height - dstHeight)/2))

	// save image back to file
	err = imaging.Save(dstFinalImg, params.tumbnailTmpPath)
	if err != nil {
		log.Fatalf("Failed to save image: %v", err)
		return err
	}

	return nil
}

// upload resized file as a response
func thumbnailUploadFile(params *thumbnailParameters, w http.ResponseWriter) error {

	fp, err := os.Open(params.tumbnailTmpPath)
	defer fp.Close() //Close after function return
	if err != nil {
		return errors.New("File Not Found")
	}

	// fill header

	//detect Content-Type of the file
	fileHeader := make([]byte, 512)
	fp.Read(fileHeader)
	fileContentType := http.DetectContentType(fileHeader)

	//get the file size
	fileStat, _ := fp.Stat()                     //Get info from file
	fileSize := strconv.FormatInt(fileStat.Size(), 10) //Get file size as a string

	//send the headers
	w.Header().Set("Content-Disposition", "attachment; filename=" + params.fileName)
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Set("Content-Length", fileSize)

	//send the file
	fp.Seek(0, 0) // return to the begining of the file
	if _, err := io.Copy(w, fp); err != nil{ //'Copy' the file to the client
		return errors.New("File Copy Error")
	}

	return nil
}

// thumbnail service handler
func thumbnailHandler(w http.ResponseWriter, r *http.Request) {
	// load client attributes, and internal information
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
	defer os.Remove(params.tumbnailTmpPath) // dont forget to delete file at the end of the session

	// resize image
	if err := thumbnailImageResize(params); err != nil {
		http.Error(w, errorStringToJson(err.Error()), http.StatusInternalServerError)
		return
	}

	// upload image to browser
	if err := thumbnailUploadFile(params, w); err != nil {
		http.Error(w, errorStringToJson(err.Error()), http.StatusInternalServerError)
	}
}
