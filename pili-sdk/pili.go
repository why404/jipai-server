package pili

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const v1 = "http://api.pili.qiniu.com/v1"

type Stream struct {
	Id            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	IsPrivate     bool   `json:"is_private"`
	StreamKey     string `json:"stream_key,omitempty"`
	StoragePeriod int    `json:"storage_period,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
	PushUrl       string `json:"push_url,omitempty"`
	LiveUrl       struct {
		HLS  string `json:"HLS,omitempty"`
		RTMP string `json:"RTMP,omitempty"`
	} `json:"live_url,omitempty"`

	nonce int
}

func (s *Stream) SignPushUrl() string {
	if s.nonce == 0 {
		s.nonce = int(time.Now().Unix())
	} else {
		s.nonce++
	}
	ret := fmt.Sprintf("%s?nonce=%d", s.PushUrl, s.nonce)
	hash := hmac.New(sha1.New, []byte(s.StreamKey))
	hash.Write([]byte(ret))
	token := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	ret = fmt.Sprintf("%s&token=%s", ret, token)
	return ret
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

func (s *Streams) Create(name, protocol, streamKey string, isPrivate bool, storagePeriod int) (*Stream, error) {
	stream := Stream{
		Name:          name,
		IsPrivate:     isPrivate,
		Protocol:      protocol,
		StreamKey:     streamKey,
		StoragePeriod: storagePeriod,
	}
	body := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(body)
	if err := encoder.Encode(stream); err != nil {
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
