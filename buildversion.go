// **********************************************************************
//    Copyright (c) 2017 Henry Seurer
//
//   Permission is hereby granted, free of charge, to any person
//    obtaining a copy of this software and associated documentation
//    files (the "Software"), to deal in the Software without
//    restriction, including without limitation the rights to use,
//    copy, modify, merge, publish, distribute, sublicense, and/or sell
//    copies of the Software, and to permit persons to whom the
//    Software is furnished to do so, subject to the following
//    conditions:
//
//   The above copyright notice and this permission notice shall be
//   included in all copies or substantial portions of the Software.
//
//    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
//    EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
//    OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//    NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
//    HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
//    WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
//    FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
//    OTHER DEALINGS IN THE SOFTWARE.
//
// **********************************************************************
//
// This program is used to generate version data we use for our services.
// It accepts the following parameters:
//		image - (required) Docker image for this application.
//		imageid - (required) Docker tagged image id.
//		output - file to output date to, the default is build_version.json
//
// This will generate a build_version.json file we use for our /version calls on our services:
//	{
//		"version": "5752d37d5bae4eacae6cd0bde82c3a24",
//		"build_time": "2017-08-13_22:21:51",
//		"image": "hello",
//		"image_id": "1234",
//		"versions": {
//			"debian": "52f04f4cf4eb4dc091ae5c2efceb7798",
//			"ghost": "c2692b51fc5f4e7bb7792b0764bdff48",
//			"nginx": "6b407ff8b85c4926ba56a780a18045ee",
//			"node": "f6eb68798add47b3b763fe44fa6b23fb",
//			"proxy": "f0be97025af54421bf8074289061b104"
//		}
//	}
//
// This allows to to trace what build/version we are using for any given image via the URL.
//
package main

import (
	"github.com/google/uuid"
	"log"
	"flag"
	"os"
	"io/ioutil"
	"encoding/json"
	"strings"
	"time"
	"github.com/henryse/go-strftime"
	"fmt"
)

func main() {
	var image, imageID string
	fileName := "build_version.json"

	boolPtr := flag.Bool("version", false, "a bool")
	flag.StringVar(&image, "image", "", "Image name for this build")
	flag.StringVar(&image, "image", "", "Image name for this build")
	flag.StringVar(&imageID, "imageid", "", "Image ID for a given docker file.")
	flag.StringVar(&fileName, "output", "", "File to write build data to")
	flag.Parse()

	if *boolPtr == true{
		fmt.Println("BuildVerion Version 1.0")
	} else {
		if len(image) == 0 || len(imageID) == 0 {
			log.Fatal("ERROR: You need to pass at lest an imageID and image name")
		}

		generate(fileName, image, imageID)
	}
}

type BuildVersion struct {
	Version   string                    `json:"version"`
	BuildTime string                    `json:"build_time"`
	Image     string                    `json:"image"`
	ImageID   string                    `json:"image_id"`
	Versions  map[string]interface{}    `json:"versions"`
}

func generate(fileName string, image string, imageID string) {

	buildVersion := new(BuildVersion)

	if fileExists(fileName) {
		dataInput, err := ioutil.ReadFile(fileName)

		if err == nil {
			json.Unmarshal(dataInput, buildVersion)
		}
	}

	update(buildVersion, image, imageID)
	write(buildVersion, fileName)
}

func getVersion() string {
	key := uuid.Must(uuid.NewRandom())
	return removeCharacter(key.String(), '-')
}

func write(buildVersion *BuildVersion, fileName string) {
	dataOutput, err := json.Marshal(buildVersion)

	if err != nil {
		log.Fatal("Marshaling failed.")
	} else {
		ioutil.WriteFile(fileName, dataOutput, 0644)
	}
}

func update(buildVersion *BuildVersion, image string, imageID string) {
	buildVersion.BuildTime = getBuildDateTime()
	buildVersion.ImageID = imageID
	buildVersion.Image = image
	buildVersion.Version = getVersion()

	if buildVersion.Versions == nil {
		buildVersion.Versions = make(map[string]interface{})
		buildVersion.Versions[image] = buildVersion.Version
	} else {
		buildVersion.Versions[image] = buildVersion.Version
	}
}

func removeCharacter(s string, c rune) string {
	return strings.Map(
		func(r rune) rune {
			if r != c {
				return r
			}
			return -1
		},
		s,
	)
}

func getBuildDateTime() string {
	t := time.Now()
	utc, _ := time.LoadLocation("UTC")
	return strftime.Format("%Y-%m-%d_%H:%M:%S", t.In(utc))
}

func fileExists(output string) bool {
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		return true
	}

	return false
}