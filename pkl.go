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
PKL struct and associated functions
http://en.wikipedia.org/wiki/Digital_Cinema_Package#Packing_list_file_or_PKL_Package_key_list
*/

package dcp

import (
	"encoding/xml"
	"io/ioutil"
	"time"
)

// PKL is returned from the parser
type PKL struct {
	ID             string `xml:"Id"`
	AnnotationText string
	IssueDate      time.Time
	Issuer         string
	Creator        string
	Assets         []*PKLAsset `xml:"AssetList>Asset"`
}

// PKLAsset is an asset found inside a PKL
type PKLAsset struct {
	ID             string `xml:"Id"`
	AnnotationText string
	Hash           string
	Size           uint64
	MimeType       string `xml:"Type"`
	// PreventMatch tag prevents clash with the Type tag
	Type AssetType `xml:"PreventMatch"`
}

// ParsePKLFile parses a PKL XML file, whose file path is asFilename
func ParsePKLFile(filename string) (*PKL, error) {
	// load the XML file
	xmlStr, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParsePKL(xmlStr)
}

// ParsePKL parses a PKL XML string
func ParsePKL(xmlBytes []byte) (*PKL, error) {
	var pkl PKL
	err := xml.Unmarshal(xmlBytes, &pkl)
	if err != nil {
		return nil, err
	}
	// Assign correct types to assets
	for _, asset := range pkl.Assets {
		switch asset.MimeType {
		case "application/x-smpte-mxf;asdcpKind=Picture":
			asset.Type = MXFPictureAssetType
		case "application/x-smpte-mxf;asdcpKind=Sound":
			asset.Type = MXFSoundAssetType
		case "text/xml;asdcpKind=CPL":
			asset.Type = CPLAssetType
		}
	}
	return &pkl, nil
}
