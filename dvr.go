package openxbl

import (
	"errors"
	"time"
)

const (
	DVRCaptureTypeClip       = DVRCaptureType("Clip")
	DVRCaptureTypeScreenshot = DVRCaptureType("Screenshot")
	DVRPrivacyBlocked        = DVRPrivacy("Blocked")
	DVRPrivacyEveryone       = DVRPrivacy("Everyone")
	DVRPrivacyPeopleOnMyList = DVRPrivacy("PeopleOnMyList")
)

type (
	DVRCaptureType string
	DVRPrivacy     string
)

type DVRCapture struct {
	ID              string `json:"contentId"`
	ContentLocators []struct {
		Expiration  time.Time `json:"expiration,omitempty"` // Only exists for screenshots
		FileSize    int       `json:"fileSize,omitempty"`
		LocatorType string    `json:"locatorType"`
		URI         string    `json:"uri"`
	} `json:"contentLocators"`
	CreationType     string         `json:"creationType"` // E.g. UserGenerated
	GreatestMomentID string         `json:"greatestMomentId"`
	LocalID          string         `json:"localId"`
	OwnerXUID        int64          `json:"ownerXuid"`
	ResolutionHeight int            `json:"resolutionHeight"`
	ResolutionWidth  int            `json:"resolutionWidth"`
	SandboxID        string         `json:"sandboxId"`
	SharedTo         []interface{}  `json:"sharedTo"`
	TitleData        string         `json:"titleData"`
	TitleID          int            `json:"titleId"`   // Game's ID
	TitleName        string         `json:"titleName"` // Game's name
	UploadDate       time.Time      `json:"uploadDate"`
	UploadLanguage   string         `json:"uploadLanguage"`
	UploadRegion     string         `json:"uploadRegion"`
	UploadTitleID    int            `json:"uploadTitleId"`
	UploadDeviceType string         `json:"uploadDeviceType"`
	UserCaption      string         `json:"userCaption"`
	CommentCount     int            `json:"commentCount"`
	LikeCount        int            `json:"likeCount"`
	ShareCount       int            `json:"shareCount"`
	ViewCount        int            `json:"viewCount"`
	ContentState     string         `json:"contentState"`
	EnforcementState string         `json:"enforcementState"`
	SafetyThreshold  string         `json:"safetyThreshold"`
	Sessions         []interface{}  `json:"sessions"`
	Tournaments      []interface{}  `json:"tournaments"`
	Type             DVRCaptureType `json:"captureType"`
}

// GetDownloadLink iterates over the ContentLocators and looks for a valid Download type. If found, it returns the URI.
func (d *DVRCapture) GetDownloadLink() string {
	for _, contentLocator := range d.ContentLocators {
		if contentLocator.LocatorType == "Download" {
			return contentLocator.URI
		}
	}

	return ""
}

type Clip struct {
	DVRCapture
	ContentSegments []struct {
		ID                int         `json:"segmentId"`
		CreationType      string      `json:"creationType"`
		CreatorChannelID  interface{} `json:"creatorChannelId"`
		CreatorXUID       int64       `json:"creatorXuid"`
		RecordDate        time.Time   `json:"recordDate"`
		DurationInSeconds int         `json:"durationInSeconds"`
		Offset            int         `json:"offset"`
		SecondaryTitleID  interface{} `json:"secondaryTitleId"`
		TitleID           int         `json:"titleId"`
	} `json:"contentSegments"` // Only exists for clips
	DurationInSeconds int `json:"durationInSeconds"` // Only exists for clips
	FrameRate         int `json:"frameRate"`         // Only exists for clips
}

func (c *Client) DeleteDVRClip(id string) error {
	if _, err := c.makeRequest("GET", "dvr/gameclips/delete/"+id, nil, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) GetDVRClips(continuationToken string) ([]*Clip, string, error) {
	response := struct {
		ContinuationToken string  `json:"continuationToken"`
		Clips             []*Clip `json:"values"`
	}{}

	// if a continuation token (their version of pagination) is supplied, pass it to the API
	endpoint := "dvr/gameclips"
	if continuationToken != "" {
		endpoint += "?continuationToken=" + continuationToken
	}

	if _, err := c.makeRequest("GET", endpoint, nil, &response); err != nil {
		return nil, "", err
	}

	if len(response.Clips) == 0 {
		return nil, "", errors.New("failed to find clips")
	}

	for index := range response.Clips {
		response.Clips[index].Type = DVRCaptureTypeClip
	}

	return response.Clips, response.ContinuationToken, nil
}

type Screenshot struct {
	DVRCapture
}

func (c *Client) GetDVRScreenshots(continuationToken string) ([]*Screenshot, string, error) {
	response := struct {
		ContinuationToken string `json:"continuationToken"`
		Screenshots       []*struct {
			DVRCapture
			DateUploaded time.Time `json:"dateUploaded"`
		} `json:"values"`
	}{}

	// if a continuation token (their version of pagination) is supplied, pass it to the API
	endpoint := "dvr/screenshots"
	if continuationToken != "" {
		endpoint += "?continuationToken=" + continuationToken
	}

	if _, err := c.makeRequest("GET", endpoint, nil, &response); err != nil {
		return nil, "", err
	}

	if len(response.Screenshots) == 0 {
		return nil, "", errors.New("failed to find screenshots")
	}

	screenshots := make([]*Screenshot, 0, len(response.Screenshots))

	for index := range response.Screenshots {
		screenshot := Screenshot{DVRCapture: response.Screenshots[index].DVRCapture}
		screenshot.DVRCapture.Type = DVRCaptureTypeScreenshot
		screenshot.UploadDate = response.Screenshots[index].DateUploaded
		screenshots = append(screenshots, &screenshot)
	}

	return screenshots, response.ContinuationToken, nil
}

func (c *Client) SetDVRPrivacy(privacy DVRPrivacy) error {
	switch privacy {
	case DVRPrivacyBlocked, DVRPrivacyEveryone, DVRPrivacyPeopleOnMyList:
	default:
		return errors.New("invalid privacy type")
	}

	request := struct {
		Privacy string `json:"value"`
	}{
		Privacy: string(privacy),
	}

	if _, err := c.makeRequest("POST", "dvr/privacy", request, nil); err != nil {
		return err
	}

	return nil
}
