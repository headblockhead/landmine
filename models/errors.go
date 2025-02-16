package models

import "net/http"

type Error struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

func NewError(message string, status int) Error {
	return Error{
		Error:  message,
		Status: status,
	}
}

type BadRequest struct {
	Error         string            `json:"error"`
	InvalidFields map[string]string `json:"invalidFields"`
	Status        int               `json:"status"`
}

func NewBadRequest(message string, invalidFields map[string]string) BadRequest {
	return BadRequest{
		Error:         message,
		InvalidFields: invalidFields,
		Status:        http.StatusBadRequest,
	}
}
