package comtradeapi

import (
	"proj/internal/app/server"
)

type Response struct {
	server.Response
	ElapsedTime string `json:"elapsedTime"`
	Error       string `json:"error"`
}

func (r *Response) verify() bool {
	if r.Count == 250000 {
		return false
	}
	return true
}

func (r *Response) Append(resp *Response) {
	r.Count += resp.Count
	r.Data = append(r.Data, resp.Data...)
	r.ElapsedTime += resp.ElapsedTime + ","
}
