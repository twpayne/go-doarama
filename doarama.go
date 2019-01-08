// Package doarama provides a client to doarama.com's API. See
// http://www.doarama.com/api/0.2/docs.
package doarama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const Version = "0.2"

// DefaultAPIURL is the default Doarama API endpoint.
const DefaultAPIURL = "https://api.doarama.com/api/0.2"

// An Error represents a Doarama server error.
type Error struct {
	HTTPStatusCode int
	HTTPStatus     string
	Status         string
	Message        string
}

// Error returns a string representation of the error.
func (e Error) Error() string {
	return fmt.Sprintf("doarama: %d %s: %s: %s", e.HTTPStatusCode, e.HTTPStatus, e.Status, e.Message)
}

// An ErrAmbiguousActivityType is returned when an activity type is ambiguous.
type ErrAmbiguousActivityType struct {
	S       string
	Matches []ActivityType
}

// Error implements error.
func (e *ErrAmbiguousActivityType) Error() string {
	candidates := make([]string, len(e.Matches))
	for i, m := range e.Matches {
		candidates[i] = fmt.Sprintf("%q", m.Name)
	}
	return fmt.Sprintf("ambiguous activity type %q, candidates are %s", e.S, strings.Join(candidates, ", "))
}

// An ErrUnknownActivityType is returned when an activity type is unknown.
type ErrUnknownActivityType struct {
	S string
}

func (e *ErrUnknownActivityType) Error() string {
	return fmt.Sprintf("unknown activity type %q", e.S)
}

// A Client is an opaque type for a Doarama client.
type Client struct {
	apiName    string
	apiKey     string
	apiURL     string
	httpClient *http.Client
	userAgent  string
	userHeader string
	user       string
}

// An ActivityInfo represents the info associated with an activity.
type ActivityInfo struct {
	TypeID        int    `json:"activityTypeId,omitempty"`
	UserName      string `json:"userName,omitempty"`
	UserAvatarURL string `json:"userAvatarUrl,omitempty"`
}

// An ActivityType is an activity type.
type ActivityType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// An ActivityTypes is an array of ActivityTypes.
type ActivityTypes []ActivityType

// Find returns the activity type that most closely matches s.
func (ats ActivityTypes) Find(s string) (ActivityType, error) {
	if id, err := strconv.Atoi(s); err == nil {
		for _, at := range ats {
			if at.ID == id {
				return at, nil
			}
		}
		return ActivityType{}, &ErrUnknownActivityType{S: s}
	}
	var matches []ActivityType
	for _, at := range ats {
		if strings.Contains(strings.ToLower(at.Name), strings.ToLower(s)) {
			matches = append(matches, at)
		}
	}
	switch len(matches) {
	case 0:
		return ActivityType{}, &ErrUnknownActivityType{S: s}
	case 1:
		return matches[0], nil
	default:
		return ActivityType{}, &ErrAmbiguousActivityType{S: s, Matches: matches}
	}
}

// FindByName returns the activity type with the given name.
func (ats ActivityTypes) FindByName(name string) (ActivityType, bool) {
	for _, at := range ats {
		if at.Name == name {
			return at, true
		}
	}
	return ActivityType{}, false
}

// FindByID returns the activity type with the given id.
func (ats ActivityTypes) FindByID(id int) (ActivityType, bool) {
	for _, at := range ats {
		if at.ID == id {
			return at, true
		}
	}
	return ActivityType{}, false
}

// An Activity represents an activity on the server.
type Activity struct {
	Client *Client
	ID     int
}

// A Coords represents a coordinate.
type Coords struct {
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Altitude         float64 `json:"altitude"`
	AltitudeAccuracy float64 `json:"altitudeAccuracy"`
	Speed            float64 `json:"speed"`
	Heading          float64 `json:"heading"`
}

// A Timestamp represents at Doarama timestamp. Doarama timestamps are in
// milliseconds since the epoch.
type Timestamp int64

// A Sample represents a live sample.
type Sample struct {
	Time     Timestamp              `json:"time"`
	Coords   Coords                 `json:"coords"`
	UserData map[string]interface{} `json:"userData,omitempty"`
}

// A Visualisation represents a visualisation on the server.
type Visualisation struct {
	Client *Client
	Key    string `json:"key"`
}

// A VisualisationURLOptions represents the options that can be set for a
// visualisation URL.
type VisualisationURLOptions struct {
	Names         []string
	Avatars       []string
	AvatarBaseURL string
	FixedAspect   bool
	MinimalView   bool
	DZML          string
}

// An ClientOption sets an option on a client.
type ClientOption func(*Client)

// NewClient creates a new unauthenticated Doarama client.
func NewClient(options ...ClientOption) *Client {
	c := &Client{
		apiURL:     DefaultAPIURL,
		httpClient: &http.Client{},
		userAgent:  defaultUserAgent(),
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// newRequest creates an authenticated HTTP request.
func (c *Client) newRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if c.apiName != "" {
		req.Header.Set("api-name", c.apiName)
	}
	if c.apiKey != "" {
		req.Header.Set("api-key", c.apiKey)
	}
	if c.userHeader != "" && c.user != "" {
		req.Header.Set(c.userHeader, c.user)
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
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
func (c *Client) doRequest(ctx context.Context, req *http.Request, v interface{}) error {
	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		var r struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(body, &r); err != nil {
			return err
		}
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

// ActivityTypes returns an array of activity types.
func (c *Client) ActivityTypes(ctx context.Context) (ActivityTypes, error) {
	req, err := c.newRequest("GET", c.apiURL+"/activityType", nil)
	if err != nil {
		return nil, err
	}
	var activityTypes ActivityTypes
	if err := c.doRequest(ctx, req, &activityTypes); err != nil {
		return nil, err
	}
	return activityTypes, nil
}

// Activity returns the activity with the specified id.
func (c *Client) Activity(id int) *Activity {
	return &Activity{
		Client: c,
		ID:     id,
	}
}

// Close releases any associated resources.
func (c *Client) Close() error {
	return nil
}

// CreateActivity creates a new activity.
func (c *Client) CreateActivity(ctx context.Context, filename string, gpsTrack io.Reader) (*Activity, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("gps_track", filename)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, gpsTrack); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	req, err := c.newRequest("POST", c.apiURL+"/activity", &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	data := struct {
		ID int `json:"id"`
	}{}
	if err := c.doRequest(ctx, req, &data); err != nil {
		return nil, err
	}
	return &Activity{
		Client: c,
		ID:     data.ID,
	}, nil
}

// CreateActivityWithInfo creates a new doarama.Activity with the specified
// doarama.ActivityInfo.
func (c *Client) CreateActivityWithInfo(ctx context.Context, filename string, gpsTrack io.Reader, activityInfo *ActivityInfo) (*Activity, error) {
	activity, err := c.CreateActivity(ctx, filename, gpsTrack)
	if err != nil {
		return activity, err
	}
	err = activity.SetInfo(ctx, activityInfo)
	return activity, err
}

// CreateLiveActivity creates a new live activity.
func (c *Client) CreateLiveActivity(ctx context.Context, startLatitude, startLongitude float64, startTime Timestamp) (*Activity, error) {
	var data = struct {
		StartLatitude  float64   `json:"startLatitude"`
		StartLongitude float64   `json:"startLongitude"`
		StartTime      Timestamp `json:"startTime"`
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
		Client: c,
	}
	if err := c.doRequest(ctx, req, a); err != nil {
		return nil, err
	}
	return a, nil
}

// CreateVisualisation creates a new visualiztion.
func (c *Client) CreateVisualisation(ctx context.Context, activities []*Activity) (*Visualisation, error) {
	data := struct {
		ActivityIds []int `json:"activityIds"`
	}{
		ActivityIds: make([]int, len(activities)),
	}
	for i, a := range activities {
		data.ActivityIds[i] = a.ID
	}
	req, err := c.newRequestJSON("POST", c.apiURL+"/visualisation", &data)
	if err != nil {
		return nil, err
	}
	v := &Visualisation{Client: c}
	if err := c.doRequest(ctx, req, v); err != nil {
		return nil, err
	}
	return v, nil
}

// Visualisation returns the visualisation with the specified key.
func (c *Client) Visualisation(key string) *Visualisation {
	return &Visualisation{
		Client: c,
		Key:    key,
	}
}

// Anonymous sets anonymous authentication.
func Anonymous(userID string) ClientOption {
	return func(c *Client) {
		c.userHeader = "user-id"
		c.user = userID
	}
}

// APIName sets the API name.
func APIName(apiName string) ClientOption {
	return func(c *Client) {
		c.apiName = apiName
	}
}

// APIKey sets the API key.
func APIKey(apiKey string) ClientOption {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

// APIURL sets the API URL.
func APIURL(apiURL string) ClientOption {
	return func(c *Client) {
		c.apiURL = apiURL
	}
}

// HTTPClient sets the http.Client used for requests.
func HTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// Delegate sets delegate authentication.
func Delegate(userKey string) ClientOption {
	return func(c *Client) {
		c.userHeader = "user-key"
		c.user = userKey
	}
}

// UserAgent sets the user agent.
func UserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// Delete deletes the activity.
func (a *Activity) Delete(ctx context.Context) error {
	req, err := a.Client.newRequest("DELETE", a.URL(), nil)
	if err != nil {
		return err
	}
	if err := a.Client.doRequest(ctx, req, nil); err != nil {
		return err
	}
	return nil
}

// Record records zero or more samples. altitudeReference should normally be
// "WGS84".
func (a *Activity) Record(ctx context.Context, samples []*Sample, altitudeReference string) error {
	data := struct {
		Samples           []*Sample `json:"samples"`
		ActivityID        int       `json:"activityId"`
		AltitudeReference string    `json:"altitudeReference"`
	}{
		Samples:           samples,
		ActivityID:        a.ID,
		AltitudeReference: altitudeReference,
	}
	req, err := a.Client.newRequestJSON("POST", a.Client.apiURL+"/activity/record", &data)
	if err != nil {
		return err
	}
	if err := a.Client.doRequest(ctx, req, nil); err != nil {
		return err
	}
	return nil
}

// SetInfo sets the info.
func (a *Activity) SetInfo(ctx context.Context, activityInfo *ActivityInfo) error {
	req, err := a.Client.newRequestJSON("POST", a.URL(), activityInfo)
	if err != nil {
		return err
	}
	if err := a.Client.doRequest(ctx, req, nil); err != nil {
		return err
	}
	return nil
}

// URL returns the URL for the activity.
func (a *Activity) URL() string {
	return a.Client.apiURL + "/activity/" + strconv.Itoa(a.ID)
}

// AddActivities adds the activities to the visualisation.
func (v *Visualisation) AddActivities(ctx context.Context, activities []*Activity) error {
	data := struct {
		VisualisationKey string `json:"visualisationKey"`
		ActivityIds      []int  `json:"activityIds"`
	}{
		VisualisationKey: v.Key,
		ActivityIds:      make([]int, len(activities)),
	}
	for i, activity := range activities {
		data.ActivityIds[i] = activity.ID
	}
	req, err := v.Client.newRequestJSON("POST", v.Client.apiURL+"/visualisation/addActivities", &data)
	if err != nil {
		return err
	}
	if err := v.Client.doRequest(ctx, req, nil); err != nil {
		return err
	}
	return nil
}

// URL returns a URL with the specificed options.
func (v *Visualisation) URL(vo *VisualisationURLOptions) *url.URL {
	u, err := url.Parse(v.Client.apiURL + "/visualisation")
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

// NewTimestamp creates a Timestamp from a time.Time.
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp(t.UnixNano() / 1000000)
}

// Time returns ts as a time.Time.
func (ts Timestamp) Time() time.Time {
	return time.Unix(int64(ts)/1000, int64(ts)%1000*1000000).UTC()
}

func defaultUserAgent() string {
	return fmt.Sprintf("go-doarama/%v (https://github.com/twpayne/go-doarama) %s (%s/%s)", Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
