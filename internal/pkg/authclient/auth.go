package authclient

import (
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
)

type HTTPResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	// If want to request other auth endpoint use json.RawMessage
	Data UserInfo `json:"data,omitempty"`
}

type UserInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var c *fasthttp.HostClient

func Init(host string) {
	c = &fasthttp.HostClient{
		Addr: host,
	}
}

func ValidateToken(token string) (string, bool) {

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI("http://" + c.Addr + "/get_user_info")
	req.Header.Set(fasthttp.HeaderAuthorization, token)
	req.Header.SetHost(c.Addr)
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := c.Do(req, resp)
	if err != nil {
		return "", false
	}
	if resp.StatusCode() != http.StatusOK {
		return "", false
	}

	var authResp HTTPResponse
	if err = json.Unmarshal(resp.Body(), &authResp); err != nil {
		return "", false
	}
	if !authResp.Success {
		return "", false
	}

	return authResp.Data.ID, true
}
