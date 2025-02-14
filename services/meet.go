package services

import (
	"backend-b7/models"
	"backend-b7/pkg/logger"
	"backend-b7/pkg/utils"
	"backend-b7/repositories"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sync"

	"net/http"
	"net/url"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type meetService struct {
	meetRepository repositories.MeetRepositoryInterface
	ClientID       string
	ClientSecret   string
	RedirectURI    string
	ZoomBaseAPI    string
	token          models.TokenResponse
	mu             sync.Mutex
}

type MeetServiceInterface interface {
	CreateMeet(meet *models.ZoomMeet) (res *models.Response, err error)
	GetMeets(filter map[string][]string) (res *models.Response, err error)
	GetMeetById(id string) (res *models.Response, err error)
	UpdateMeet(meet *models.ZoomMeetUpdate) (res *models.Response, err error)
	DeleteMeet(meet *models.ZoomMeetUpdate) (res *models.Response, err error)

	RequestAccessToken(code string) (*models.TokenResponse, error)
}

func NewMeetService(meetRepository repositories.MeetRepositoryInterface, clientID, clientSecret, redirectURI, zoomBaseAPI string) MeetServiceInterface {
	return &meetService{
		meetRepository: meetRepository,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		RedirectURI:    redirectURI,
		ZoomBaseAPI:    zoomBaseAPI,
		token:          models.TokenResponse{},
		mu:             sync.Mutex{},
	}
}

func (u *meetService) CreateMeet(meet *models.ZoomMeet) (res *models.Response, err error) {

	zoomMeet, err := u.createZoomMeeting(meet)
	if err != nil {
		return nil, err
	}

	meet.ID = uuid.New().String()
	meet.Status = models.StatusActive
	meet.MeetingID = int64(zoomMeet["id"].(float64))
	meet.JoinURL = fmt.Sprintf("%v", zoomMeet["join_url"])
	meet.StartTime = meet.StartTime.UTC()
	err = u.meetRepository.CreateMeet(meet)
	if err != nil {
		return nil, err
	}

	return &models.Response{
		Code:    http.StatusCreated,
		Message: "Meet created successfully",
	}, nil
}

func (u *meetService) GetMeets(filter map[string][]string) (res *models.Response, err error) {

	pagination, search := utils.GeneratePaginationFromRequest(filter)
	meets, count, err := u.meetRepository.GetMeets(pagination, search)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return &models.Response{
			Code:    http.StatusNotFound,
			Message: http.StatusText(http.StatusNotFound),
		}, nil
	}

	data := models.ListZoomMeet{
		Page:      pagination.Page,
		Limit:     pagination.Limit,
		Total:     int(count),
		TotalPage: int(math.Ceil(float64(count) / float64(pagination.Limit))),
		ZoomMeets: meets,
	}

	return &models.Response{
		Code:    http.StatusOK,
		Message: "Meet list successfully",
		Data:    data,
	}, nil
}

func (u *meetService) GetMeetById(id string) (res *models.Response, err error) {

	meet, err := u.meetRepository.GetMeetById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &models.Response{
				Code:    http.StatusNotFound,
				Message: "Meet not exist",
			}, nil
		}
		return nil, err
	}

	return &models.Response{
		Code:    http.StatusOK,
		Message: "Meet get successfully",
		Data:    meet,
	}, nil
}

func (u *meetService) UpdateMeet(meet *models.ZoomMeetUpdate) (res *models.Response, err error) {

	if meet.Status == "" {
		meet.Status = models.StatusActive
	}

	_, err = u.updateZoomMeeting(meet)
	if err != nil {
		return nil, err
	}

	meet.StartTime = meet.StartTime.UTC()
	err = u.meetRepository.UpdateMeet(meet)
	if err != nil {
		return nil, err
	}

	return &models.Response{
		Code:    http.StatusOK,
		Message: "Meet updated successfully",
	}, nil
}

func (u *meetService) DeleteMeet(meet *models.ZoomMeetUpdate) (res *models.Response, err error) {
	err = u.meetRepository.DeleteMeet(meet)
	if err != nil {
		return nil, err
	}

	return &models.Response{
		Code:    http.StatusOK,
		Message: "Meet deleted successfully",
	}, nil
}

func (u *meetService) setToken(token models.TokenResponse) {
	u.mu.Lock()
	defer u.mu.Unlock()

	logger.Infof("Set token: %v", token.AccessToken)
	u.token = models.TokenResponse{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		ExpiresIn:    token.ExpiresIn,
		Scope:        token.Scope,
		RefreshToken: token.RefreshToken,
	}
}

func (u *meetService) getToken() models.TokenResponse {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.token
}

func (u *meetService) RequestAccessToken(code string) (*models.TokenResponse, error) {
	// Prepare the request body
	body := url.Values{}
	body.Set("grant_type", "authorization_code")
	body.Set("code", code)
	body.Set("redirect_uri", u.RedirectURI)

	// Create the HTTP request
	req, err := http.NewRequest("POST", "https://zoom.us/oauth/token", bytes.NewBufferString(body.Encode()))
	if err != nil {
		return nil, err
	}

	// Set headers
	authHeader := base64.StdEncoding.EncodeToString([]byte(u.ClientID + ":" + u.ClientSecret))
	req.Header.Set("Authorization", "Basic "+authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var tokenResponse models.TokenResponse
	if err := json.Unmarshal(respBody, &tokenResponse); err != nil {
		return nil, err
	}

	if tokenResponse.AccessToken == "" {
		return nil, fmt.Errorf("no access token in response: %s", respBody)
	}

	u.setToken(tokenResponse)

	return &tokenResponse, nil
}

func (u *meetService) createZoomMeeting(meet *models.ZoomMeet) (resp map[string]interface{}, err error) {

	var uri = fmt.Sprintf("%s/users/%s/meetings", u.ZoomBaseAPI, meet.UserID)

	body := map[string]interface{}{
		"topic":      meet.Topic,
		"type":       2,
		"start_time": meet.StartTime,
		"duration":   meet.Duration,
	}
	reqBody, _ := json.Marshal(body)

	respBody, respCode, err := utils.CallRESTAPIWithToken(uri, "POST", reqBody, u.getToken().AccessToken)
	if err != nil {
		return nil, err
	}

	resp = map[string]interface{}{}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}

	if respCode != 200 {
		return nil, fmt.Errorf("failed to create meeting")
	}

	return
}

func (u *meetService) updateZoomMeeting(meet *models.ZoomMeetUpdate) (resp map[string]interface{}, err error) {

	var uri = fmt.Sprintf("%s/meetings/%v", u.ZoomBaseAPI, meet.MeetingID)

	body := map[string]interface{}{
		"topic":      meet.Topic,
		"type":       2,
		"start_time": meet.StartTime,
		"duration":   meet.Duration,
	}
	reqBody, _ := json.Marshal(body)

	_, respCode, err := utils.CallRESTAPIWithToken(uri, "PATCH", reqBody, u.getToken().AccessToken)
	if err != nil {
		return nil, err
	}

	if respCode != http.StatusNoContent {
		return nil, fmt.Errorf("failed to update meeting")
	}

	return
}
