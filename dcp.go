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

package dcp

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// DCP represents all the components found within a DCP
type DCP struct {
	RootDir  string // root directory of the DCP
	AssetMap *AssetMap
	CPLs     []*CPL
	PKLs     []*PKL

	assetMapFile string
}

// String produces a human-readable representation of a DCP
func (dcp *DCP) String() string {
	var dcpStr string
	switch dcp.Format() {
	case INTEROP:
		dcpStr += "Type: " + "Interop\n"
	case SMPTE:
		dcpStr += "Type: " + "SMPTE\n"
	}
	dcpStr += "AssetMap: " + dcp.AssetMap.ID + "\n"
	for _, cpl := range dcp.CPLs {
		dcpStr += "CPL: " + cpl.AnnotationText + "\n"
	}
	for _, pkl := range dcp.PKLs {
		dcpStr += "PKL: " + pkl.AnnotationText + "\n"
	}
	return dcpStr
}

// Format returns the format of the DCP (INTEROP OR SMPTE)
func (dcp *DCP) Format() Format {
	return dcp.AssetMap.Format
}

// Files returns the physical files that comprise the DCP
func (dcp *DCP) Files() []string {
	return append([]string{dcp.assetMapFile}, dcp.AssetMap.Paths()...)
}

// Generate builds a new DCP from a root directory path containing an assetmap
func (dcp *DCP) Generate(dir string) error {
	amFileName, err := findAssetMap(dir)
	if err != nil {
		return err
	}
	am, err := ParseAssetMapFile(amFileName)
	if err != nil {
		return err
	}
	dcp.RootDir = dir
	dcp.assetMapFile = filepath.Base(amFileName)
	for _, asset := range am.Assets {
		for _, chunk := range asset.Chunks {
			// Check the file size
			assetPath := filepath.Join(dir, chunk.Path)
			err = checkFileSize(assetPath, chunk.Size)
			if err != nil {
				return err
			}
			// Determine the asset type
			aType := assetType(assetPath)
			// Parse CPLs and PKLs
			if aType == CPLAssetType {
				cpl, err := ParseCPLFile(assetPath)
				if err != nil {
					return err
				}
				dcp.CPLs = append(dcp.CPLs, cpl)
			}
			if aType == PKLAssetType {
				pkl, err := ParsePKLFile(assetPath)
				if err != nil {
					return err
				}
				dcp.PKLs = append(dcp.PKLs, pkl)
			}
		}
	}
	dcp.AssetMap = am
	return nil
}

/*
findAssetMap looks in a directory for an assetmap
and if found returns its absolute path
*/
func findAssetMap(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, f := range files {
		if regexp.MustCompile(`^(assetmap|ASSETMAP)(.xml|.XML)*$`).
			MatchString(f.Name()) {
			return filepath.Join(dir, f.Name()), nil
		}
	}
	return "", errors.New("Unable to find an assetmap file")
}

// Check if a file exists and that it's the correct size
func checkFileSize(filename string, size uint64) error {
	info, err := os.Stat(filename)
	if err != nil {
		return err
	}
	fileSize := info.Size()
	if uint64(fileSize) != size {
		return errors.New("File size for" + filename + " is incorrect")
	}
	return nil
}

// mxfHeader is the starting bytes of an audio or picture MXF file
var mxfHeader = []byte{6, 14, 43, 52, 2, 5, 1, 1, 13, 1, 2, 1, 1, 2,
	4, 0, 131, 0, 0, 120, 0, 1, 0, 2, 0, 0, 0, 1}

// Determine the asset type from the file
func assetType(filepath string) AssetType {
	data, err := readHeader(filepath)
	if err != nil {
		return UnknownAssetType
	}
	switch {
	case bytes.Contains(data, []byte("PackingList")):
		return PKLAssetType
	case bytes.Contains(data, []byte("CompositionPlaylist")):
		return CPLAssetType
	case bytes.HasPrefix(data, mxfHeader):
		return MXFAssetType
	default:
		return UnknownAssetType
	}
}

// readHeader reads the first 100 bytes from a file;
// returns an error if the file can't be read
func readHeader(filepath string) ([]byte, error) {
	data := make([]byte, 100)
	file, err := os.Open(filepath)
	if err == nil {
		_, err = file.Read(data)
	}
	return data, err
}
