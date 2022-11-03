package output

//go:generate gomodifytags --file $GOFILE --struct HttpGeneralBody -add-tags json -w -transform snakecase
type HttpGeneralBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//go:generate gomodifytags --file $GOFILE --struct HttpUnprocessableContentError -add-tags json -w -transform snakecase
type HttpUnprocessableContentError struct {
	Attribute string `json:"attribute"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

//go:generate gomodifytags --file $GOFILE --struct HttpUnprocessableContent -add-tags json -w -transform snakecase
type HttpUnprocessableContent struct {
	Code    int                              `json:"code"`
	Message string                           `json:"message"`
	Errors  *[]HttpUnprocessableContentError `json:"errors"`
}

func NewHttpUnauthorized() HttpGeneralBody {
	return HttpGeneralBody{
		Code:    401,
		Message: "Unauthorized",
	}
}

func NewHttpNotFound() HttpGeneralBody {
	return HttpGeneralBody{
		Code:    404,
		Message: "Not Found",
	}
}

func NewHttpConflict() HttpGeneralBody {
	return HttpGeneralBody{
		Code:    409,
		Message: "Conflict",
	}
}

func NewHttpUnprocessableContent(errors *[]HttpUnprocessableContentError) HttpUnprocessableContent {
	return HttpUnprocessableContent{
		Code:    422,
		Message: "Unprocessable Content",
		Errors:  errors,
	}
}

func NewHttpInternalServerError() string {
	return "{\"code\":500,\"message\":\"Internal Server Error\"}"
}
