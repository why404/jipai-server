package pili

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	HTTPCode int    `json:"-"`
	Message  string `json:"message"`
	ErrCode  int    `json:"error"`
	Details  map[string]struct {
		Message string `json:"message"`
		ErrCode int    `json:"error"`
	} `json:"details"`
}

func (e *Error) SetStatus(code int) {
	e.HTTPCode = code
}

func (e Error) Status() int {
	return e.HTTPCode
}

func (e Error) Error() string {
	return fmt.Sprintf("(%d/%d)%s", e.HTTPCode, e.ErrCode, e.Message)
}

func handleResp(resp *http.Response, reply interface{}) error {
	decoder := json.NewDecoder(resp.Body)
	if resp.StatusCode != http.StatusOK {
		e := new(Error)
		if err := decoder.Decode(e); err != nil {
			return err
		}
		e.SetStatus(resp.StatusCode)
		return e
	}

	if reply == nil {
		return nil
	}

	if err := decoder.Decode(reply); err != nil {
		return err
	}
	return nil
}
