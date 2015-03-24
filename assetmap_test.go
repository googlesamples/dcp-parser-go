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
	"testing"
	"time"
)

// XML doc taken from DCPs available at http://www.freedcp.net
var testAssetMapXML = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<AssetMap xmlns="http://www.digicine.com/PROTO-ASDCP-AM-20040311#">
  <Id>urn:uuid:88ef5d99-e2aa-483e-9697-943e18b77cea</Id>
  <Creator>OpenDCP 0.0.26</Creator>
  <VolumeCount>1</VolumeCount>
  <IssueDate>2012-09-28T03:40:08+00:00</IssueDate>
  <Issuer>FreeDCP.net</Issuer>
  <AssetList>
    <Asset>
      <Id>urn:uuid:4d9e98c3-c923-4910-ae0e-9f5951c9cc5f</Id>
      <PackingList>true</PackingList>
      <ChunkList>
        <Chunk>
          <Path>4d9e98c3-c923-4910-ae0e-9f5951c9cc5f_pkl.xml</Path>
          <VolumeIndex>1</VolumeIndex>
          <Offset>0</Offset>
          <Length>1288</Length>
        </Chunk>
      </ChunkList>
    </Asset>
    <Asset>
      <Id>urn:uuid:d65572db-2e09-4745-817d-a2881222e2db</Id>
      <ChunkList>
        <Chunk>
          <Path>d65572db-2e09-4745-817d-a2881222e2db_cpl.xml</Path>
          <VolumeIndex>1</VolumeIndex>
          <Offset>0</Offset>
          <Length>1550</Length>
        </Chunk>
      </ChunkList>
    </Asset>
    <Asset>
      <Id>urn:uuid:db95199c-0e2f-4ac4-9e54-b97919dcdf07</Id>
      <ChunkList>
        <Chunk>
          <Path>bewegte_bilder-tricks17-test_film-full_content-51-j2k_video.mxf</Path>
          <VolumeIndex>1</VolumeIndex>
          <Offset>0</Offset>
          <Length>3906847916</Length>
        </Chunk>
      </ChunkList>
    </Asset>
    <Asset>
      <Id>urn:uuid:5fbb3067-4166-4a19-9ba2-0a2b4c5cd397</Id>
      <ChunkList>
        <Chunk>
          <Path>bewegte_bilder-tricks17-test_film-full_content-51-j2k_audio.mxf</Path>
          <VolumeIndex>1</VolumeIndex>
          <Offset>0</Offset>
          <Length>879052345</Length>
      </Chunk>
      </ChunkList>
    </Asset>
  </AssetList>
</AssetMap>`)

func parseAM(t *testing.T) *AssetMap {
	assetmap, err := ParseAssetMap(testAssetMapXML)
	if err != nil {
		t.Errorf("%s", err)
	}
	return assetmap
}

func TestAssetMap(t *testing.T) {
	assetmap := parseAM(t)
	// test Format
	if assetmap.Format != INTEROP {
		t.Errorf("Type is incorrect: %d != %d", assetmap.Format, INTEROP)
	}
	// test ID
	expectedID := "urn:uuid:88ef5d99-e2aa-483e-9697-943e18b77cea"
	if assetmap.ID != expectedID {
		t.Errorf("ID is incorrect: %s != %s",
			assetmap.ID, expectedID)
	}
	// Test Creator
	creatorID := "OpenDCP 0.0.26"
	if assetmap.Creator != creatorID {
		t.Errorf("Creator is incorrect: %s != %s",
			assetmap.Creator, creatorID)
	}
	// test VolumeCount
	volumeCount := uint8(1)
	if assetmap.VolumeCount != volumeCount {
		t.Errorf("VolumeCount is incorrect: %v != %v",
			assetmap.VolumeCount, volumeCount)
	}
	// test Issuer
	issuer := "FreeDCP.net"
	if assetmap.Issuer != issuer {
		t.Errorf("Issuer is incorrect: %s != %s",
			assetmap.Issuer, issuer)
	}
	// test IssueDate
	issueDate, err := time.Parse(time.RFC3339Nano, "2012-09-28T03:40:08+00:00")
	if err != nil {
		t.Errorf("%s", err)
	}
	if !assetmap.IssueDate.Equal(issueDate) {
		t.Errorf("IssueDate is incorrect: %s != %s",
			assetmap.IssueDate, issueDate)
	}
}

func TestAMNrAssets(t *testing.T) {
	assetmap := parseAM(t)
	if len(assetmap.Assets) != 4 {
		t.Errorf("Number of assets is incorrect: %d != %d", len(assetmap.Assets), 4)
	}
}

func TestAMSize(t *testing.T) {
	assetmap := parseAM(t)
	var amSize uint64 = 1288 + 1550 + 3906847916 + 879052345
	if assetmap.Size() != amSize {
		t.Errorf("Size of assetmap is incorrect: %d != %d",
			assetmap.Size(), amSize)
	}
}

func TestAMPaths(t *testing.T) {
	assetmap := parseAM(t)
	var paths = []string{"4d9e98c3-c923-4910-ae0e-9f5951c9cc5f_pkl.xml",
		"d65572db-2e09-4745-817d-a2881222e2db_cpl.xml",
		"bewegte_bilder-tricks17-test_film-full_content-51-j2k_video.mxf",
		"bewegte_bilder-tricks17-test_film-full_content-51-j2k_audio.mxf"}
	for i, path := range assetmap.Paths() {
		if path != paths[i] {
			t.Errorf("Path mismatch: %s != %s", path, paths[i])
		}
	}
}

func TestAMAssetSize(t *testing.T) {
	assetmap := parseAM(t)
	var amAssetSize uint64 = 1288
	if assetmap.Assets[0].Size() != amAssetSize {
		t.Errorf("Size of asset is incorrect: %d != %d",
			assetmap.Assets[0].Size(), amAssetSize)
	}
	amAssetSize = 1550
	if assetmap.Assets[1].Size() != amAssetSize {
		t.Errorf("Size of asset is incorrect: %d != %d",
			assetmap.Assets[1].Size(), amAssetSize)
	}
}

func TestAMAssetPaths(t *testing.T) {
	assetmap := parseAM(t)
	var paths = []string{"4d9e98c3-c923-4910-ae0e-9f5951c9cc5f_pkl.xml",
		"d65572db-2e09-4745-817d-a2881222e2db_cpl.xml",
		"bewegte_bilder-tricks17-test_film-full_content-51-j2k_video.mxf",
		"bewegte_bilder-tricks17-test_film-full_content-51-j2k_audio.mxf"}
	for i, asset := range assetmap.Assets {
		if asset.Chunks[0].Path != paths[i] {
			t.Errorf("Path mismatch: %s != %s", asset.Chunks[0].Path, paths[i])
		}
	}
}

func TestAMPKLAsset(t *testing.T) {
	assetmap := parseAM(t)
	asset := assetmap.Assets[0]
	assetID := "urn:uuid:4d9e98c3-c923-4910-ae0e-9f5951c9cc5f"
	if asset.ID != assetID {
		t.Errorf("Asset Id is incorrect: %s != %s", asset.ID, assetID)
	}
	if asset.Type != PKLAssetType {
		t.Errorf("Asset type is incorrect: %d != %d", asset.Type, PKLAssetType)
	}
	assetPath := "4d9e98c3-c923-4910-ae0e-9f5951c9cc5f_pkl.xml"
	if asset.Chunks[0].Path != assetPath {
		t.Errorf("Asset path is incorrect: %s != %s", asset.Chunks[0].Path, assetPath)
	}
	var assetSize uint64 = 1288
	if asset.Chunks[0].Size != assetSize {
		t.Errorf("Asset size is incorrect: %d != %d", asset.Chunks[0].Size, assetSize)
	}
}

func TestAMCPLAsset(t *testing.T) {
	assetmap := parseAM(t)
	asset := assetmap.Assets[1]
	assetID := "urn:uuid:d65572db-2e09-4745-817d-a2881222e2db"
	if asset.ID != assetID {
		t.Errorf("Asset Id is incorrect: %s != %s", asset.ID, assetID)
	}
	if asset.Type != CPLAssetType {
		t.Errorf("Asset type is incorrect: %d != %d", asset.Type, CPLAssetType)
	}
	assetPath := "d65572db-2e09-4745-817d-a2881222e2db_cpl.xml"
	if asset.Chunks[0].Path != assetPath {
		t.Errorf("Asset path is incorrect: %s != %s", asset.Chunks[0].Path, assetPath)
	}
	var assetSize uint64 = 1550
	if asset.Chunks[0].Size != assetSize {
		t.Errorf("Asset size is incorrect: %d != %d", asset.Chunks[0].Size, assetSize)
	}
}

func TestAMPictureAsset(t *testing.T) {
	assetmap := parseAM(t)
	asset := assetmap.Assets[2]
	assetID := "urn:uuid:db95199c-0e2f-4ac4-9e54-b97919dcdf07"
	if asset.ID != assetID {
		t.Errorf("Asset Id is incorrect: %s != %s", asset.ID, assetID)
	}
	if asset.Type != MXFAssetType {
		t.Errorf("Asset type is incorrect: %d != %d", asset.Type, MXFAssetType)
	}
	assetPath := "bewegte_bilder-tricks17-test_film-full_content-51-j2k_video.mxf"
	if asset.Chunks[0].Path != assetPath {
		t.Errorf("Asset path is incorrect: %s != %s", asset.Chunks[0].Path, assetPath)
	}
	var assetSize uint64 = 3906847916
	if asset.Chunks[0].Size != assetSize {
		t.Errorf("Asset size is incorrect: %d != %d", asset.Chunks[0].Size, assetSize)
	}
}

func TestAMSoundAsset(t *testing.T) {
	assetmap := parseAM(t)
	asset := assetmap.Assets[3]
	assetID := "urn:uuid:5fbb3067-4166-4a19-9ba2-0a2b4c5cd397"
	if asset.ID != assetID {
		t.Errorf("Asset Id is incorrect: %s != %s", asset.ID, assetID)
	}
	if asset.Type != MXFAssetType {
		t.Errorf("Asset type is incorrect: %d != %d", asset.Type, MXFAssetType)
	}
	assetPath := "bewegte_bilder-tricks17-test_film-full_content-51-j2k_audio.mxf"
	if asset.Chunks[0].Path != assetPath {
		t.Errorf("Asset path is incorrect: %s != %s", asset.Chunks[0].Path, assetPath)
	}
	var assetSize uint64 = 879052345
	if asset.Chunks[0].Size != assetSize {
		t.Errorf("Asset size is incorrect: %d != %d", asset.Chunks[0].Size, assetSize)
	}
}
