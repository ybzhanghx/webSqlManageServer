package controllers

import "errors"

const (
	successGetMsg  = "success to get"
	successSaveMsg = "success to save"
)

var (
	FailGetErr = errors.New("failed to get")
)

type CommonReturn struct {
	Code int
	Msg  string
}

func (c *CommonReturn) SetData(code int, msg string) {
	c.Code, c.Msg = code, msg
}
