package api

type (
	ResponseError struct {
		Code    string            `json:"code"`
		Message string            `json:"message"`
		Reasons map[string]string `json:"reasons"`
		Details []interface{}     `json:"details,omitempty"`
	}
	Response struct {
		Status    int                    `json:"status"`
		Content   map[string]interface{} `json:"content,omitempty"`
		Error     *ResponseError         `json:"error,omitempty"`
	}
)
