package controllers

type CommonReturn struct {
	Code int
	Msg  string
}

func (c *CommonReturn) SetData(code int, msg string) {
	c.Code, c.Msg = code, msg
}
