package gomemory

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// BaseURL will be used as base API url
var BaseURL = "https://api.mymemory.translated.net/"

// EmailDomains used for setting email suffix when "rand" parameter is used
var EmailDomains = []string{
	"gmail.com",
	"outlook.com",
	"yahoo.com",
	"pm.me",
}

// OParameters represents high-level parameters for Object translate
type OParameters struct {
	// API
	Value interface{}
	Src   string
	Dest  string
	Key   string
	Email string

	// Internal
	Timeout time.Duration // Default is 10 seconds
}

// Parameters represents translate parameters
type Parameters struct {
	// API
	Text     string
	Src      string
	Dest     string
	MimeType string
	Key      string
	Email    string // Will be used as "de" parameter to extend limits, use "gen" for random value

	// Internal
	Timeout time.Duration // Default is 10 seconds
}

// Response represents API response
type Response struct {
	Data ResponseData `json:"responseData"`

	// Internal
	Splitted bool
}

// ResponseData represents API response data
type ResponseData struct {
	Text       string  `json:"translatedText"`
	MatchLevel float64 `json:"match"`
}

// Translate is a main low-level translation function
func translate(p Parameters) (Response, error) {
	// Validate
	if p.Src == "" {
		return Response{}, errors.New("Src parameter is required")
	}
	if p.Dest == "" {
		return Response{}, errors.New("Dest paramter is required")
	}
	// Init url
	turl, err := url.Parse(BaseURL)
	if err != nil {
		return Response{}, err
	}
	// Set path
	turl.Path = "/get"
	// Set query
	q := turl.Query()
	q.Set("q", p.Text)
	q.Set("langpair", p.Src+"|"+p.Dest)
	if p.Key != "" {
		q.Set("key", p.Key)
	}
	if p.Email != "" {
		q.Set("de", p.Email)
	}
	turl.RawQuery = q.Encode()
	// Make request
	if p.Timeout == 0 {
		p.Timeout = 10 * time.Second
	}
	client := &http.Client{
		Timeout: p.Timeout,
	}
	rawresp, err := client.Get(turl.String())
	if err != nil {
		return Response{}, err
	}
	if rawresp.StatusCode != 200 {
		return Response{}, errors.New("Status code: " + strconv.Itoa(rawresp.StatusCode))
	}
	// Decode
	var resp Response
	err = json.NewDecoder(rawresp.Body).Decode(&resp)
	if err != nil {
		return Response{}, err
	}
	// Return
	return resp, nil
}

// Translate is a main high-level translation function
func Translate(p Parameters) (Response, error) {
	// Generate email, if needed
	if p.Email == "gen" {
		val := rand.Intn(1000000)
		domain := EmailDomains[rand.Intn(len(EmailDomains))]
		p.Email = strconv.Itoa(val) + "@" + domain
	}
	// Handle small text as-is
	if len(p.Text) < 400 {
		return translate(p)
	}
	// Handle big text with splitting
	texts := strings.Split(p.Text, ".") // Split texts
	textsres := []string{}              // Result texts store
	for _, text := range texts {        // Translate one by one
		if text == "" { // Pass empty text (case when dot in the end)
			continue
		}
		_p := p
		_p.Text = text
		resp, err := translate(_p)
		if err != nil {
			return Response{}, err
		}
		textsres = append(textsres, resp.Data.Text)
	}
	// Build result text
	res := strings.Join(textsres, ".")
	// Postprocessing
	res = strings.ReplaceAll(res, "..", ".")
	// Construct half-filled response and return
	return Response{
		Data: ResponseData{
			Text:       res,
			MatchLevel: 1,
		},
	}, nil
}

func TranslateObject(p OParameters) (interface{}, error) {
	// Build list of texts from object
	texts := []string{}
	switch p.Value.(type) {
	case string:
		texts = append(texts, p.Value.(string))
	case []string:
		texts = append(texts, p.Value.([]string)...)
	}
	// Translate
	textsres := []string{}
	for _, text := range texts {
		r, err := Translate(Parameters{
			Text:    text,
			Src:     p.Src,
			Dest:    p.Dest,
			Key:     p.Key,
			Email:   p.Email,
			Timeout: p.Timeout,
		})
		if err != nil {
			return nil, err
		}
		textsres = append(textsres, r.Data.Text)
	}
	// Built object back
	var v interface{}
	switch p.Value.(type) {
	case string:
		v = textsres[0]
	case []string:
		v = textsres
	}
	return v, nil
}
