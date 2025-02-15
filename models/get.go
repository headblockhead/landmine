package models

import (
	"net/http"
	"time"
)

type SortDirection string

const (
	SortAscending  SortDirection = "asc"
	SortDescending SortDirection = "desc"
)

func NewSortDirection(s string) SortDirection {
	switch s {
	case "asc":
		return SortAscending
	case "desc":
		return SortDescending
	default:
		return SortAscending
	}
}

type CellFormat string

const (
	CellFormatJSON   CellFormat = "json"
	CellFormatString CellFormat = "string"
)

func NewCellFormat(s string) CellFormat {
	switch s {
	case "json":
		return CellFormatJSON
	case "string":
		return CellFormatString
	default:
		return CellFormatJSON
	}
}

// https://airtable.com/developers/web/api/list-records
type ListGetRequest struct {
	// Path parameters
	BaseID        string
	TableIDOrName string

	// Query parameters
	TimeZone              string        `json:"timeZone"`
	UserLocale            string        `json:"userLocale"`
	PageSize              int           `json:"pageSize"`
	MaxRecords            int           `json:"maxRecords"`
	Offset                string        `json:"offset"`
	View                  string        `json:"view"`
	Sort                  SortDirection `json:"sort"`
	FilterByFormula       string        `json:"filterByFormula"`
	CellFormat            CellFormat    `json:"cellFormat"`
	Fields                []string      `json:"fields"`
	ReturnFieldsByFieldId bool          `json:"returnFieldsByFieldId"`
	RecordMetadata        []string      `json:"recordMetadata"`
}

type Record struct {
	ID           string         `json:"id"`
	CreatedTime  time.Time      `json:"createdTime"`
	Fields       map[string]any `json:"fields"`
	CommentCount int            `json:"commentCount,omitempty"`
}

type ListGetResponse struct {
	Offset  string   `json:"offset,omitempty"`
	Records []Record `json:"records"`
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
