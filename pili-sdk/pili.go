package pili

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

const v1 = "http://api.pili.qiniu.com/v1"

type StreamUrl map[string]string

const StreamDefaultLive = "[original]"

type Stream struct {
	ID           string               `bson:"_id" json:"id"`
	CreatedAt    time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time            `bson:"updated_at" json:"updated_at"`
	Application  string               `bson:"application" json:"application"`
	IsPrivate    bool                 `bson:"is_private" json:"is_private"`
	StreamKey    string               `bson:"stream_key" json:"stream_key"`
	PushNonce    int64                `bson:"push_nonce" json:"-"`
	PushUrl      []StreamUrl          `bson:"push_url" json:"push_url"`
	LiveUrl      map[string]StreamUrl `bson:"live_url" json:"live_url"`
	LivestreamId string               `bson:"livestream_id" json:"-"`
}

type Streams struct {
	client *http.Client
	mac    Mac
}

func New(mac Mac) *Streams {
	return &Streams{
		client: http.DefaultClient,
		mac:    mac,
	}
}

func (s *Streams) Create(app, streamKey string, isPrivate bool) (*Stream, error) {
	args := map[string]interface{}{
		"application": app,
		"is_private":  isPrivate,
	}
	if streamKey != "" {
		args["stream_key"] = streamKey
	}
	body := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(body)
	if err := encoder.Encode(args); err != nil {
		return nil, err
	}
	req, err := s.mac.newRequest("POST", "/streams", body)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var stream Stream
	if err := handleResp(resp, &stream); err != nil {
		return nil, err
	}
	return &stream, nil
}

func (s *Streams) Get(id string) (*Stream, error) {
	req, err := s.mac.newRequest("GET", "/streams/"+id, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ret Stream
	if err := handleResp(resp, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (s *Streams) Delete(id string) (*Stream, error) {
	req, err := s.mac.newRequest("DELETE", "/streams/"+id, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ret Stream
	if err := handleResp(resp, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}
