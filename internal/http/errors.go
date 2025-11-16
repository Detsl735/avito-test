package http

func errorResponse(code, message string) ErrorResponse {
	var er ErrorResponse
	er.Error.Code = code
	er.Error.Message = message
	return er
}

func errorBadRequest(msg string) ErrorResponse {
	return errorResponse("BAD_REQUEST", msg)
}
