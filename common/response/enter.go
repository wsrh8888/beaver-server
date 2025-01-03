package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
}

// Response allows for a custom message to be passed in case of success
func Response(r *http.Request, w http.ResponseWriter, resp interface{}, err error, successMsg ...string) {
	msg := ""
	if len(successMsg) > 0 {
		msg = successMsg[0]
	}
	if err == nil {
		r := &Body{
			Code:   0,
			Msg:    msg,
			Result: resp,
		}
		httpx.WriteJson(w, http.StatusOK, r)
		return
	}
	httpx.WriteJson(w, http.StatusOK, &Body{
		Code:   1,
		Msg:    err.Error(),
		Result: nil,
	})
}
