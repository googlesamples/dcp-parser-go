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

// Format of the two DCP standards
type Format int

// Format types
const (
	UNKNOWN Format = iota
	INTEROP
	SMPTE
)

// AssetType of the different file components of a DCP
type AssetType int

// Asset types
const (
	UnknownAssetType AssetType = iota
	CPLAssetType
	PKLAssetType
	MXFAssetType
	MXFPictureAssetType
	MXFSoundAssetType
)

// IsMxf checks if an asset is an MXF file
func IsMxf(assetType AssetType) bool {
	return assetType >= MXFAssetType && assetType <= MXFSoundAssetType
}
