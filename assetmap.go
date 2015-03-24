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
Assetmap struct and associated functions
http://en.wikipedia.org/wiki/Digital_Cinema_Package#Asset_map_file
*/

package dcp

import (
	"encoding/xml"
	"io/ioutil"
	"regexp"
	"time"
)

// AssetMap is the struct produced by the parser
type AssetMap struct {
	Format      Format
	ID          string
	Creator     string
	VolumeCount uint8
	Issuer      string
	IssueDate   time.Time
	Assets      []*AMAsset
}

// AMAsset is an asset map asset
type AMAsset struct {
	ID     string
	Type   AssetType
	Chunks []*Chunk
}

// Chunk is a single file and a component of an asset
type Chunk struct {
	Path string
	Size uint64 `xml:"Length"`
}

// Size is the summed size of all the assets referenced by the asset map
func (am AssetMap) Size() uint64 {
	var totalSize uint64
	for _, asset := range am.Assets {
		totalSize += asset.Size()
	}
	return totalSize
}

// Paths returns all asset file paths
func (am AssetMap) Paths() []string {
	var paths []string
	for _, asset := range am.Assets {
		paths = append(paths, asset.Paths()...)
	}
	return paths
}

// Paths returns all file paths of an asset
func (a AMAsset) Paths() []string {
	var paths []string
	for _, chunk := range a.Chunks {
		paths = append(paths, chunk.Path)
	}
	return paths
}

// Size returns the total size of an asset
func (a AMAsset) Size() uint64 {
	var size uint64
	for _, chunk := range a.Chunks {
		size += chunk.Size
	}
	return size
}

// ParseAssetMapFile parses an asset map xml file, whose file path is asFilename
func ParseAssetMapFile(amFilename string) (*AssetMap, error) {
	xmlStr, err := ioutil.ReadFile(amFilename)
	if err != nil {
		return nil, err
	}
	return ParseAssetMap(xmlStr)
}

// ParseAssetMap parses an asset map XML string
func ParseAssetMap(xmlStr []byte) (*AssetMap, error) {
	var amXML assetMapXML
	if err := xml.Unmarshal(xmlStr, &amXML); err != nil {
		return nil, err
	}
	return makeAssetMap(amXML)
}

/*
Structure used to unmarshall the XML; used internally only - AssetMap struct
is passed back by ParseAssetMap() & ParseAssetMapFile()
*/
type assetMapXML struct {
	Xmlns       string `xml:"xmlns,attr"`
	ID          string `xml:"Id"`
	Creator     string
	VolumeCount uint8
	IssueDate   time.Time
	Issuer      string
	Assets      []*assetXML `xml:"AssetList>Asset"`
}

type assetXML struct {
	ID          string `xml:"Id"`
	PackingList string
	Chunks      []*Chunk `xml:"ChunkList>Chunk"`
}

// Creates an AssetMap from a raw assetMapXML
func makeAssetMap(amXML assetMapXML) (*AssetMap, error) {
	assetMap := &AssetMap{
		ID:          amXML.ID,
		Creator:     amXML.Creator,
		VolumeCount: amXML.VolumeCount,
		IssueDate:   amXML.IssueDate,
		Issuer:      amXML.Issuer}
	// Set the type
	if amXML.Xmlns == "http://www.digicine.com/PROTO-ASDCP-AM-20040311#" {
		assetMap.Format = INTEROP
	} else if amXML.Xmlns == "http://www.smpte-ra.org/schemas/429-9/2007/AM" {
		assetMap.Format = SMPTE
	}
	// Convert the xml assets to Asset
	for _, assetXML := range amXML.Assets {
		assetMap.Assets = append(assetMap.Assets,
			&AMAsset{assetXML.ID, guessAssetType(assetXML), assetXML.Chunks})
	}
	return assetMap, nil
}

// Regular expressions for guessing file type from the file name
var cplRegExp = regexp.MustCompile(`(cpl|CPL)(.xml|.XML)$`)
var pklRegExp = regexp.MustCompile(`(pkl|PKL)(.xml|.XML)$`)
var mxfRegExp = regexp.MustCompile(`(.mxf|.MXF)$`)

// Attempt to guess the asset type; it may be wrong
func guessAssetType(aXML *assetXML) AssetType {
	assetType := UnknownAssetType
	if aXML.PackingList == "true" {
		assetType = PKLAssetType
	} else {
		// Look at the each chunk's path and guess the asset's type
		for _, chunk := range aXML.Chunks {
			if cplRegExp.MatchString(chunk.Path) {
				assetType = CPLAssetType
			} else if pklRegExp.MatchString(chunk.Path) {
				assetType = PKLAssetType
			} else if mxfRegExp.MatchString(chunk.Path) {
				assetType = MXFAssetType
			}
		}
	}
	return assetType
}
