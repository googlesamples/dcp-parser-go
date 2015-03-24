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
var testCPLXML = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<CompositionPlaylist xmlns="http://www.digicine.com/PROTO-ASDCP-CPL-20040511#">
  <Id>urn:uuid:d65572db-2e09-4745-817d-a2881222e2db</Id>
  <AnnotationText>Bewegte Bilder - Tricks17 - Test Film - Full Content - 5.1 - JPEG2000</AnnotationText>
  <IssueDate>2012-09-28T03:40:08+00:00</IssueDate>
  <Creator>OpenDCP 0.0.26</Creator>
  <ContentTitleText>Bewegte Bilder - Tricks17 - Test Film - Full Content - 5.1 - JPEG2000</ContentTitleText>
  <ContentKind>test</ContentKind>
  <RatingList/>
  <ReelList>
    <Reel>
      <Id>urn:uuid:da41abd1-a7df-4383-a79f-438dde1d4fd3</Id>
      <AssetList>
        <MainPicture>
          <Id>urn:uuid:db95199c-0e2f-4ac4-9e54-b97919dcdf07</Id>
          <AnnotationText>bewegte_bilder-tricks17-test_film-full_content-51-j2k_video.mxf</AnnotationText>
          <EditRate>24 1</EditRate>
          <IntrinsicDuration>23400</IntrinsicDuration>
          <EntryPoint>0</EntryPoint>
          <Duration>23400</Duration>
          <FrameRate>24 1</FrameRate>
          <ScreenAspectRatio>1.90</ScreenAspectRatio>
        </MainPicture>
        <MainSound>
          <Id>urn:uuid:5fbb3067-4166-4a19-9ba2-0a2b4c5cd397</Id>
          <AnnotationText>bewegte_bilder-tricks17-test_film-full_content-51-j2k_audio.mxf</AnnotationText>
          <EditRate>24 1</EditRate>
          <IntrinsicDuration>24404</IntrinsicDuration>
          <EntryPoint>0</EntryPoint>
          <Duration>23400</Duration>
        </MainSound>
      </AssetList>
    </Reel>
  </ReelList>
</CompositionPlaylist>`)

func parseCPL(t *testing.T) *CPL {
	cpl, err := ParseCPL(testCPLXML)
	if err != nil {
		t.Errorf("%s", err)
	}
	return cpl
}

func TestCPL(t *testing.T) {
	cpl := parseCPL(t)
	// test Format
	if cpl.Format != INTEROP {
		t.Errorf("Format is incorrect: %d != %d",
			cpl.Format, INTEROP)
	}
	// test ID
	expectedID := "urn:uuid:d65572db-2e09-4745-817d-a2881222e2db"
	if cpl.ID != expectedID {
		t.Errorf("ID is incorrect: %s != %s",
			cpl.ID, expectedID)
	}
	// test AnnotationText
	expectedText := "Bewegte Bilder - Tricks17 - Test Film - Full Content - 5.1 - JPEG2000"
	if cpl.AnnotationText != expectedText {
		t.Errorf("AnnotationText is incorrect: %s != %s",
			cpl.AnnotationText, expectedText)
	}
	// test IssueDate
	issueDate, err := time.Parse(time.RFC3339Nano, "2012-09-28T03:40:08+00:00")
	if err != nil {
		t.Errorf("%s", err)
	}
	if !cpl.IssueDate.Equal(issueDate) {
		t.Errorf("IssueDate is incorrect: %s != %s",
			cpl.IssueDate, issueDate)
	}
	// test Creator
	expectedCreator := "OpenDCP 0.0.26"
	if cpl.Creator != expectedCreator {
		t.Errorf("Creator is incorrect: %s != %s",
			cpl.Creator, expectedCreator)
	}
	// test ContentTitleText
	expectedTitle := "Bewegte Bilder - Tricks17 - Test Film - Full Content - 5.1 - JPEG2000"
	if cpl.ContentTitleText != expectedTitle {
		t.Errorf("ContentTitleText is incorrect: %s != %s",
			cpl.ContentTitleText, expectedTitle)
	}
	// test ContentKind
	expectedKind := testCPLKind
	if cpl.ContentKind != expectedKind {
		t.Errorf("ContentKind is incorrect: %d != %d",
			cpl.ContentKind, expectedKind)
	}
	// test reel count
	expectedReelCount := 1
	if len(cpl.Reels) != expectedReelCount {
		t.Errorf("Reel count is incorrect: %d != %d", len(cpl.Reels), expectedReelCount)
	}
	// test reel picture id
	expectedPicID := "urn:uuid:db95199c-0e2f-4ac4-9e54-b97919dcdf07"
	if cpl.Reels[0].Picture.ID != expectedPicID {
		t.Errorf("Picture asset id is incorrect: %s != %s", cpl.Reels[0].Picture.ID, expectedPicID)
	}
	// test picture annotation text
	expectedText = "bewegte_bilder-tricks17-test_film-full_content-51-j2k_video.mxf"
	if cpl.Reels[0].Picture.AnnotationText != expectedText {
		t.Errorf("Picture annotation text is incorrect: %s != %s",
			cpl.Reels[0].Picture.AnnotationText, expectedText)
	}
	// test reel sound id
	expectedSoundID := "urn:uuid:5fbb3067-4166-4a19-9ba2-0a2b4c5cd397"
	if cpl.Reels[0].Sound.ID != expectedSoundID {
		t.Errorf("Sound asset id is incorrect: %s != %s", cpl.Reels[0].Sound.ID, expectedSoundID)
	}
	// test sound annotation text
	expectedText = "bewegte_bilder-tricks17-test_film-full_content-51-j2k_audio.mxf"
	if cpl.Reels[0].Sound.AnnotationText != expectedText {
		t.Errorf("Sound annotation text is incorrect: %s != %s",
			cpl.Reels[0].Sound.AnnotationText, expectedText)
	}
	// test reel subtitle id
	if cpl.Reels[0].Subtitle != nil {
		t.Errorf("Subtitle asset id should be nil")
	}
}

func TestCPLPictureCount(t *testing.T) {
	cpl := parseCPL(t)
	expectedPictureCount := 1
	if len(cpl.Pictures()) != expectedPictureCount {
		t.Errorf("Picture count is incorrect: %d != %d", len(cpl.Pictures()), expectedPictureCount)
	}
}

func TestCPLSoundCount(t *testing.T) {
	cpl := parseCPL(t)
	expectedSoundCount := 1
	if len(cpl.Sounds()) != expectedSoundCount {
		t.Errorf("Sound count is incorrect: %d != %d", len(cpl.Sounds()), expectedSoundCount)
	}
}

func TestCPLSubtitleCount(t *testing.T) {
	cpl := parseCPL(t)
	expectedSubtitleCount := 0
	if len(cpl.Subtitles()) != expectedSubtitleCount {
		t.Errorf("Subtitles count is incorrect: %d != %d", len(cpl.Subtitles()), expectedSubtitleCount)
	}
}
