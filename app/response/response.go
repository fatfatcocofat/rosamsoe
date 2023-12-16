package response

const (
	SERVER_ERROR_MSG = "An error has occurred while processing the request"
	BAD_GATEWAY_MSG  = "A communication error occurred with the data source"
	NOT_FOUND_MSG    = "An error occurred or the resource could not be found"
)

type (
	BadRequest struct {
		Success bool        `json:"success"`
		Message string      `json:"message,omitempty"`
		Errors  interface{} `json:"errors,omitempty"`
	}

	Unauthorized struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
	}

	Forbidden struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
	}

	ServerError struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
	}

	BadGateway struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
	}

	NotFound struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
	}

	Success struct {
		Success bool        `json:"success"`
		Message string      `json:"message,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}

	TokenData struct {
		Token     string `json:"token"`
		ExpiresIn int64  `json:"expires_in"`
	}

	TokenResponse struct {
		Success bool      `json:"success"`
		Data    TokenData `json:"data"`
	}
)
