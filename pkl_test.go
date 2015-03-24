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
var testPKLXML = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<PackingList xmlns="http://www.digicine.com/PROTO-ASDCP-PKL-20040311#">
  <Id>urn:uuid:4d9e98c3-c923-4910-ae0e-9f5951c9cc5f</Id>
  <AnnotationText>Bewegte Bilder - Tricks17 - Test Film - Full Content - 5.1 - JPEG2000</AnnotationText>
  <IssueDate>2012-09-28T03:40:08+00:00</IssueDate>
  <Issuer>FreeDCP.net</Issuer>
  <Creator>OpenDCP 0.0.26</Creator>
  <AssetList>
    <Asset>
      <Id>urn:uuid:db95199c-0e2f-4ac4-9e54-b97919dcdf07</Id>
      <AnnotationText>bewegte_bilder-tricks17-test_film-full_content-51-j2k_video.mxf</AnnotationText>
      <Hash>R53oZs3TlqpWkRa1AhcduI8glak=</Hash>
      <Size>3906847916</Size>
      <Type>application/x-smpte-mxf;asdcpKind=Picture</Type>
    </Asset>
    <Asset>
      <Id>urn:uuid:5fbb3067-4166-4a19-9ba2-0a2b4c5cd397</Id>
      <AnnotationText>bewegte_bilder-tricks17-test_film-full_content-51-j2k_audio.mxf</AnnotationText>
      <Hash>6M60C7/OjyhH+tjBlaHrPreXWIU=</Hash>
      <Size>879052345</Size>
      <Type>application/x-smpte-mxf;asdcpKind=Sound</Type>
    </Asset>
    <Asset>
      <Id>urn:uuid:d65572db-2e09-4745-817d-a2881222e2db</Id>
      <Hash>IWCQsMJUrIEighw/7ViCPeP7nw4=</Hash>
      <Size>1550</Size>
      <Type>text/xml;asdcpKind=CPL</Type>
    </Asset>
  </AssetList>
</PackingList>`)

func parsePKL(t *testing.T) *PKL {
	pkl, err := ParsePKL(testPKLXML)
	if err != nil {
		t.Errorf("%s", err)
	}
	return pkl
}

func TestPKL(t *testing.T) {
	pkl := parsePKL(t)
	// test ID
	expectedID := "urn:uuid:4d9e98c3-c923-4910-ae0e-9f5951c9cc5f"
	if pkl.ID != expectedID {
		t.Errorf("ID is incorrect: %s != %s",
			pkl.ID, expectedID)
	}
	// test AnnotationText
	expectedText := "Bewegte Bilder - Tricks17 - Test Film - Full Content - 5.1 - JPEG2000"
	if pkl.AnnotationText != expectedText {
		t.Errorf("AnnotationText is incorrect: %s != %s",
			pkl.AnnotationText, expectedText)
	}
	// test IssueDate
	issueDate, err := time.Parse(time.RFC3339Nano, "2012-09-28T03:40:08+00:00")
	if err != nil {
		t.Errorf("%s", err)
	}
	if !pkl.IssueDate.Equal(issueDate) {
		t.Errorf("IssueDate is incorrect: %s != %s",
			pkl.IssueDate, issueDate)
	}
	// test Issuer
	expectedIssuer := "FreeDCP.net"
	if pkl.Issuer != expectedIssuer {
		t.Errorf("Issuer is incorrect: %s != %s",
			pkl.Issuer, expectedIssuer)
	}
	// test Creator
	expectedCreator := "OpenDCP 0.0.26"
	if pkl.Creator != expectedCreator {
		t.Errorf("Creator is incorrect: %s != %s",
			pkl.Creator, expectedCreator)
	}
	// test asset count
	expectedCount := 3
	if len(pkl.Assets) != expectedCount {
		t.Errorf("Asset count is incorrect: %d != %d",
			len(pkl.Assets), expectedCount)
	}
	// test picture asset id
	expectedID = "urn:uuid:db95199c-0e2f-4ac4-9e54-b97919dcdf07"
	if pkl.Assets[0].ID != expectedID {
		t.Errorf("Picture asset id is incorrect: %s != %s",
			pkl.Assets[0].ID, expectedID)
	}
	// test picture annotation text
	expectedText = "bewegte_bilder-tricks17-test_film-full_content-51-j2k_video.mxf"
	if pkl.Assets[0].AnnotationText != expectedText {
		t.Errorf("Picture asset annotation text is incorrect: %s != %s",
			pkl.Assets[0].AnnotationText, expectedText)
	}
	// test picture hash
	expectedHash := "R53oZs3TlqpWkRa1AhcduI8glak="
	if pkl.Assets[0].Hash != expectedHash {
		t.Errorf("Picture asset hash is incorrect: %s != %s",
			pkl.Assets[0].Hash, expectedHash)
	}
	// test picture size
	expectedSize := uint64(3906847916)
	if pkl.Assets[0].Size != expectedSize {
		t.Errorf("Picture asset size is incorrect: %d != %d",
			pkl.Assets[0].Size, expectedSize)
	}
	// test picture mime type
	expectedType := "application/x-smpte-mxf;asdcpKind=Picture"
	if pkl.Assets[0].MimeType != expectedType {
		t.Errorf("Picture asset mime type is incorrect: %d != %s",
			pkl.Assets[0].Type, expectedType)
	}
	// test picture type
	if pkl.Assets[0].Type != MXFPictureAssetType {
		t.Errorf("Picture asset type is incorrect: %d != %d",
			pkl.Assets[0].Type, MXFPictureAssetType)
	}
	// test sound id
	expectedID = "urn:uuid:5fbb3067-4166-4a19-9ba2-0a2b4c5cd397"
	if pkl.Assets[1].ID != expectedID {
		t.Errorf("Sound asset id is incorrect: %s != %s",
			pkl.Assets[1].ID, expectedID)
	}
	// test sound annotation text
	expectedText = "bewegte_bilder-tricks17-test_film-full_content-51-j2k_audio.mxf"
	if pkl.Assets[1].AnnotationText != expectedText {
		t.Errorf("Sound asset annotation text is incorrect: %s != %s",
			pkl.Assets[1].AnnotationText, expectedText)
	}
	// test sound hash
	expectedHash = "6M60C7/OjyhH+tjBlaHrPreXWIU="
	if pkl.Assets[1].Hash != expectedHash {
		t.Errorf("Sound asset hash is incorrect: %s != %s",
			pkl.Assets[1].Hash, expectedHash)
	}
	// test sound size
	expectedSize = uint64(879052345)
	if pkl.Assets[1].Size != expectedSize {
		t.Errorf("Sound asset size is incorrect: %d != %d",
			pkl.Assets[1].Size, expectedSize)
	}
	// test sound mime type
	expectedType = "application/x-smpte-mxf;asdcpKind=Sound"
	if pkl.Assets[1].MimeType != expectedType {
		t.Errorf("Sound asset mime type is incorrect: %s != %s",
			pkl.Assets[1].MimeType, expectedType)
	}
	// test sound type
	if pkl.Assets[1].Type != MXFSoundAssetType {
		t.Errorf("Sound asset type is incorrect: %d != %d",
			pkl.Assets[1].Type, MXFSoundAssetType)
	}
	// test cpl asset id
	expectedID = "urn:uuid:d65572db-2e09-4745-817d-a2881222e2db"
	if pkl.Assets[2].ID != expectedID {
		t.Errorf("CPL asset id is incorrect: %s != %s",
			pkl.Assets[2].ID, expectedID)
	}
	// test cpl hash
	expectedHash = "IWCQsMJUrIEighw/7ViCPeP7nw4="
	if pkl.Assets[2].Hash != expectedHash {
		t.Errorf("CPL asset hash is incorrect: %s != %s",
			pkl.Assets[2].Hash, expectedHash)
	}
	// test cpl size
	expectedSize = uint64(1550)
	if pkl.Assets[2].Size != expectedSize {
		t.Errorf("CPL asset size is incorrect: %d != %d",
			pkl.Assets[2].Size, expectedSize)
	}
	// test cpl mime type
	expectedType = "text/xml;asdcpKind=CPL"
	if pkl.Assets[2].MimeType != expectedType {
		t.Errorf("CPL asset mime type is incorrect: %s != %s",
			pkl.Assets[2].MimeType, expectedType)
	}
	// test cpl type
	if pkl.Assets[2].Type != CPLAssetType {
		t.Errorf("CPL asset type is incorrect: %d != %d",
			pkl.Assets[2].Type, CPLAssetType)
	}
}
