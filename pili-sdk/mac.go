package pili

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
)

type Mac struct {
	AccessKey string
	SecretKey string
}

func (m *Mac) newRequest(method, urlStr string, body *bytes.Buffer) (*http.Request, error) {
	fmt.Println("req:", method, urlStr, body.String())
	var req *http.Request
	var err error
	if body == nil {
		req, err = http.NewRequest(method, v1+urlStr, nil)
	} else {
		req, err = http.NewRequest(method, v1+urlStr, body)
	}
	if err != nil {
		return nil, err
	}
	u := *req.URL
	u.Scheme = ""
	u.Host = ""
	hash := hmac.New(sha1.New, []byte(m.SecretKey))
	hash.Write([]byte(u.String()))
	hash.Write([]byte("\n"))
	if body != nil {
		hash.Write(body.Bytes())
	}
	token := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	req.Header.Set("Authorization", fmt.Sprintf("pili %s:%s", m.AccessKey, token))
	return req, nil
}
