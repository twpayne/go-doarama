// Package doarama provides a client to doarama.com's API. See http://www.doarama.com/api/0.2/docs.
package doarama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

const API_URL = "https://api.doarama.com/api/0.2"

// Type Error represents a Doarama server error.
type Error struct {
	HTTPStatusCode int
	HTTPStatus     string
	Status         string
	Message        string
}

// Error returns a string representation of the error.
func (e Error) Error() string {
	return fmt.Sprintf("doarama: %s: %s: %s", e.HTTPStatus, e.Status, e.Message)
}

// Client is an opaque type for a Doarama client.
type Client struct {
	apiName    string
	apiKey     string
	apiURL     string
	client     *http.Client
	userHeader string
	user       string
}

// Type ActivityInfo represents the info associated with an activity.
type ActivityInfo struct {
	TypeId int `json:"activityTypeId"`
}

// Type Activity represents an activity on the server.
type Activity struct {
	c                 *Client
	Id                int
	InfoURL           string
	AuthorURL         string
	altitudeReference string
}

// Type Coords represents a coordinate.
type Coords struct {
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Altitude         float64 `json:"altitude"`
	AltitudeAccuracy float64 `json:"altitudeAccuracy"`
	Speed            float64 `json:"speed"`
	Heading          float64 `json:"heading"`
}

// Type Sample represents a live sample.
type Sample struct {
	Time     int64                  `json:"time"`
	Coords   Coords                 `json:"coords"`
	UserData map[string]interface{} `json:"userData,omitempty"`
}

// Type Visualisation represents a visualisation on the server.
type Visualisation struct {
	c   *Client
	Key string `json:"key"`
}

// Type VisualisationURLOptions represents that options that can be set for a visualisation URL.
type VisualisationURLOptions struct {
	Names         []string
	Avatars       []string
	AvatarBaseURL string
	FixedAspect   bool
	MinimalView   bool
	DZML          string
}

// NewClient creates a new unauthenticated Doarama client.
func NewClient(apiURL, apiName, apiKey string) *Client {
	return &Client{
		apiName: apiName,
		apiKey:  apiKey,
		apiURL:  apiURL,
		client:  &http.Client{},
	}
}

// Anonymous creates a new client based on c using anonymous authentication.
func (c *Client) Anonymous(userId string) *Client {
	anonymous := *c
	anonymous.userHeader = "user-id"
	anonymous.user = userId
	return &anonymous
}

// Delegate creates a new client based on c using delegate authentication.
func (c *Client) Delegate(userKey string) *Client {
	delegate := *c
	delegate.userHeader = "user-key"
	delegate.user = userKey
	return &delegate
}

// newRequest creates an authenticated HTTP request.
func (c *Client) newRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("api-name", c.apiName)
	req.Header.Set("api-key", c.apiKey)
	req.Header.Set(c.userHeader, c.user)
	return req, nil
}

// newRequestJSON creates an authenticated HTTP request with a JSON body.
func (c *Client) newRequestJSON(method, urlStr string, v interface{}) (*http.Request, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest(method, urlStr, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// doRequest performs an HTTP request and unmarshals the JSON response.
func (c *Client) doRequest(req *http.Request, v interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		var r struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}
		json.Unmarshal(body, &r)
		return Error{
			HTTPStatusCode: resp.StatusCode,
			HTTPStatus:     resp.Status,
			Status:         r.Status,
			Message:        r.Message,
		}
	}
	if v != nil {
		if err := json.Unmarshal(body, v); err != nil {
			return err
		}
	}
	return nil
}

// ActivityTypes returns a map of activity types to activity type ids.
func (c *Client) ActivityTypes() (map[string]int, error) {
	req, err := c.newRequest("GET", c.apiURL+"/activityType", nil)
	if err != nil {
		return nil, err
	}
	var activityTypesResponse []struct {
		Id   float64 `json:"id"`
		Name string  `json:"name"`
	}
	if err := c.doRequest(req, &activityTypesResponse); err != nil {
		return nil, err
	}
	result := make(map[string]int)
	for _, at := range activityTypesResponse {
		result[at.Name] = int(at.Id)
	}
	return result, nil
}

// Activity returns the activity with the specified id and altitude reference. altitudeReference is usually "WGS84".
func (c *Client) Activity(id int, altitudeReference string) *Activity {
	return &Activity{
		c:                 c,
		Id:                id,
		altitudeReference: altitudeReference,
	}
}

// CreateActivity creates a new activity.
func (c *Client) CreateActivity(filename string, gpsTrack io.Reader) (*Activity, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("gps_track", filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(fw, gpsTrack); err != nil {
		return nil, err
	}
	w.Close()
	req, err := c.newRequest("POST", c.apiURL+"/activity", &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	data := struct {
		Id        int    `json:"id"`
		InfoURL   string `json:"infoUrl"`
		AuthorURL string `json:"authorUrl"`
	}{}
	if err := c.doRequest(req, &data); err != nil {
		return nil, err
	}
	return &Activity{
		c:         c,
		Id:        data.Id,
		InfoURL:   data.InfoURL,
		AuthorURL: data.AuthorURL,
	}, nil
}

// CreateLiveActivity creates a new live activity.
func (c *Client) CreateLiveActivity(startLatitude, startLongitude float64, startTime int64, altitudeReference string) (*Activity, error) {
	var data = struct {
		StartLatitude  float64 `json:"startLatitude"`
		StartLongitude float64 `json:"startLongitude"`
		StartTime      int64   `json:"startTime"`
	}{
		StartLatitude:  startLatitude,
		StartLongitude: startLongitude,
		StartTime:      startTime,
	}
	req, err := c.newRequestJSON("POST", c.apiURL+"/activity/create", &data)
	if err != nil {
		return nil, err
	}
	a := &Activity{
		c:                 c,
		altitudeReference: altitudeReference,
	}
	if err := c.doRequest(req, a); err != nil {
		return nil, err
	}
	return a, nil
}

// CreateVisualisation creates a new visualiztion.
func (c *Client) CreateVisualisation(activities []*Activity) (*Visualisation, error) {
	data := struct {
		ActivityIds []int `json:"activityIds"`
	}{
		ActivityIds: make([]int, len(activities)),
	}
	for i, a := range activities {
		data.ActivityIds[i] = a.Id
	}
	req, err := c.newRequestJSON("POST", c.apiURL+"/visualisation", &data)
	if err != nil {
		return nil, err
	}
	v := &Visualisation{c: c}
	if err := c.doRequest(req, v); err != nil {
		return nil, err
	}
	return v, nil
}

// Visualisation returns the visualisation with the specified key.
func (c *Client) Visualisation(key string) *Visualisation {
	return &Visualisation{
		c:   c,
		Key: key,
	}
}

// Delete deletes the activity.
func (a *Activity) Delete() error {
	req, err := a.c.newRequest("DELETE", a.URL(), nil)
	if err != nil {
		return err
	}
	if err := a.c.doRequest(req, nil); err != nil {
		return err
	}
	return nil
}

// Record records zero or more samples.
func (a *Activity) Record(samples []*Sample) error {
	data := struct {
		Samples           []*Sample `json:"samples"`
		ActivityId        int       `json:"activityId"`
		AltitudeReference string    `json:"altitudeReference"`
	}{
		Samples:           samples,
		ActivityId:        a.Id,
		AltitudeReference: a.altitudeReference,
	}
	req, err := a.c.newRequestJSON("POST", a.c.apiURL+"/activity/record", &data)
	if err != nil {
		return err
	}
	if err := a.c.doRequest(req, nil); err != nil {
		return err
	}
	return nil
}

// SetInfo sets the info.
func (a *Activity) SetInfo(activityInfo *ActivityInfo) error {
	req, err := a.c.newRequestJSON("POST", a.URL(), activityInfo)
	if err != nil {
		return err
	}
	if err := a.c.doRequest(req, nil); err != nil {
		return err
	}
	return nil
}

// URL returns the URL for the activity.
func (a *Activity) URL() string {
	return a.c.apiURL + "/activity/" + strconv.Itoa(a.Id)
}

// AddActivities adds the activities to the visualisation.
func (v *Visualisation) AddActivities(activities []*Activity) error {
	data := struct {
		VisualisationKey string `json:"visualisationKey"`
		ActivityIds      []int  `json:"activityIds"`
	}{
		VisualisationKey: v.Key,
		ActivityIds:      make([]int, len(activities)),
	}
	for i, activity := range activities {
		data.ActivityIds[i] = activity.Id
	}
	req, err := v.c.newRequestJSON("POST", v.c.apiURL+"/visualisation/addActivities", &data)
	if err != nil {
		return err
	}
	if err := v.c.doRequest(req, nil); err != nil {
		return err
	}
	return nil
}

// Delete deletes the visualisation.
func (v *Visualisation) Delete() error {
	req, err := v.c.newRequest("DELETE", v.URL(nil).String(), nil)
	if err != nil {
		return err
	}
	if err := v.c.doRequest(req, nil); err != nil {
		return err
	}
	return nil
}

// URL returns a URL with the specificed options.
func (v *Visualisation) URL(vo *VisualisationURLOptions) *url.URL {
	u, err := url.Parse(v.c.apiURL + "/visualisation")
	if err != nil {
		panic(err)
	}
	values := u.Query()
	values.Set("k", v.Key)
	if vo != nil {
		if vo.Names != nil {
			for _, n := range vo.Names {
				values.Add("name", n)
			}
		}
		if vo.Avatars != nil {
			for _, a := range vo.Avatars {
				values.Add("avatar", a)
			}
		}
		if vo.AvatarBaseURL != "" {
			values.Set("avatarBaseUrl", vo.AvatarBaseURL)
		}
		if !vo.FixedAspect {
			values.Set("fixedAspect", "false")
		}
		if vo.MinimalView {
			values.Set("minimalView", "true")
		}
		if vo.DZML != "" {
			values.Set("dzml", vo.DZML)
		}
	}
	u.RawQuery = values.Encode()
	return u
}
