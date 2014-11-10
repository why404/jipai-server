package main

import (
	"fmt"
)

type Error struct {
	HTTPCode int    `json:"-"`
	Message  string `json:"message"`
	ErrCode  int    `json:"error"`
}

func (e Error) Status() int {
	return e.HTTPCode
}

func (e Error) Error() string {
	return fmt.Sprintf("(%d/%d)%s", e.HTTPCode, e.ErrCode, e.Message)
}
