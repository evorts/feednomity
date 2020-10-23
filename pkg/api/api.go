package api

type (
	ErrorResponse struct {
		Code    string            `json:"code"`
		Message string            `json:"message"`
		Reasons map[string]string `json:"reasons"`
		Details []interface{}     `json:"details,omitempty"`
	}
	Response struct {
		RequestID string                 `json:"request_id"`
		Status    int                    `json:"status"`
		Content   map[string]interface{} `json:"content,omitempty"`
		Error     *ErrorResponse         `json:"error,omitempty"`
	}
)
