//
//  Copyright 2015  Google Inc. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

/*
CPL struct and associated functions
http://en.wikipedia.org/wiki/Digital_Cinema_Package#Composition_playlist_file
*/

package dcp

import (
	"encoding/xml"
	"io/ioutil"
	"time"
)

// ContentKind is the type of content referenced by a CPL
type ContentKind int

const (
	unkownCPLKind ContentKind = iota
	testCPLKind
	featureCPLKind
	advertisementCPLKind
)

// CPL struct is returned by the parser
type CPL struct {
	Format           Format
	ID               string
	AnnotationText   string
	Creator          string
	ContentTitleText string
	IssueDate        time.Time
	ContentKind      ContentKind
	Reels            []*Reel
}

// Asset is a CPL asset
type Asset struct {
	ID                string `xml:"Id"`
	AnnotationText    string
	EditRate          string
	IntrinsicDuration uint64
	EntryPoint        uint64
	Duration          uint64
}

// Picture is a specific form of a CPL asset
type Picture struct {
	Asset
	FrameRate         string
	ScreenAspectRatio string
}

// Sound is a specific form of a CPL asset
type Sound struct {
	Asset
	Language string
}

// Subtitle is a specific form of a CPL asset
type Subtitle struct {
	Asset
	Language string
}

// Reel is a reel from a CPL
type Reel struct {
	ID       string    `xml:"Id"`
	Picture  *Picture  `xml:"AssetList>MainPicture"`
	Sound    *Sound    `xml:"AssetList>MainSound"`
	Subtitle *Subtitle `xml:"AssetList>MainSubtitle"`
}

// Pictures returns all the picture assets in a CPL
func (cpl CPL) Pictures() []*Picture {
	pictures := make([]*Picture, 0, len(cpl.Reels))
	for _, reel := range cpl.Reels {
		if reel.Picture != nil {
			pictures = append(pictures, reel.Picture)
		}
	}
	return pictures
}

// Sounds returns all the sound assets in a CPL
func (cpl CPL) Sounds() []*Sound {
	sounds := make([]*Sound, 0, len(cpl.Reels))
	for _, reel := range cpl.Reels {
		if reel.Sound != nil {
			sounds = append(sounds, reel.Sound)
		}
	}
	return sounds
}

// Subtitles returns all the subtitle assets in a CPL
func (cpl CPL) Subtitles() []*Subtitle {
	subtitles := make([]*Subtitle, 0, len(cpl.Reels))
	for _, reel := range cpl.Reels {
		if reel.Subtitle != nil {
			subtitles = append(subtitles, reel.Subtitle)
		}
	}
	return subtitles
}

// ParseCPLFile parses a CPL XML file, whose file path is asFilename
func ParseCPLFile(filename string) (*CPL, error) {
	// load the XML file
	xmlStr, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseCPL(xmlStr)
}

// ParseCPL parses a CPL XML string
func ParseCPL(xmlStr []byte) (*CPL, error) {
	var cplXML cplXML
	err := xml.Unmarshal(xmlStr, &cplXML)
	if err != nil {
		return nil, err
	}
	cpl, err := makeCPL(&cplXML)
	return cpl, err
}

/*
cplXML is used to unmarshall XML; used internally only - CPL struct
is passed back by ParseCPL() & ParseCPLFile()
*/
type cplXML struct {
	Xmlns            string `xml:"xmlns,attr"`
	ID               string `xml:"Id"`
	AnnotationText   string
	IssueDate        time.Time
	Creator          string
	ContentTitleText string
	ContentKind      string
	Reels            []*Reel `xml:"ReelList>Reel"`
}

// makeCPL creates a CPL from a raw cplXML
func makeCPL(cplXML *cplXML) (*CPL, error) {
	cpl := CPL{
		ID:               cplXML.ID,
		AnnotationText:   cplXML.AnnotationText,
		IssueDate:        cplXML.IssueDate,
		Creator:          cplXML.Creator,
		ContentTitleText: cplXML.ContentTitleText}
	if cplXML.Xmlns == "http://www.digicine.com/PROTO-ASDCP-CPL-20040511#" {
		cpl.Format = INTEROP
	} else if cplXML.Xmlns == "http://www.smpte-ra.org/schemas/429-7/2006/CPL" {
		cpl.Format = SMPTE
	}
	// Set the content kind
	switch cplXML.ContentKind {
	case "test":
		cpl.ContentKind = testCPLKind
	case "feature":
		cpl.ContentKind = featureCPLKind
	case "advertisement":
		cpl.ContentKind = advertisementCPLKind
	}
	cpl.Reels = cplXML.Reels
	return &cpl, nil
}
