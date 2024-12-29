package v1

import "encoding/json"

// Response represents the response of subject access review request.
type Response struct {
	Allowed bool   `json:"allowed"`          // 通过
	Denied  bool   `json:"denied,omitempty"` // 拒绝
	Reason  string `json:"reason,omitempty"` // 原因
	Error   string `json:"error,omitempty"`  // 错误信息
}

// ToString marshal Response struct to a json string.
func (rsp *Response) ToString() string {
	data, _ := json.Marshal(rsp)

	return string(data)
}
