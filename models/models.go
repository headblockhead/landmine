package models

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

type Sort struct {
	Field     string
	Direction SortDirection
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

type Record struct {
	ID           string         `json:"id"`
	CreatedTime  time.Time      `json:"createdTime"`
	Fields       map[string]any `json:"fields"`
	CommentCount int            `json:"commentCount,omitempty"`
}

// https://airtable.com/developers/web/api/list-records
type ListRecordsRequest struct {
	// Path parameters
	BaseID        string
	TableIDOrName string

	// Query parameters
	TimeZone              string     `json:"timeZone"`
	UserLocale            string     `json:"userLocale"`
	PageSize              int        `json:"pageSize"`
	MaxRecords            int        `json:"maxRecords"`
	Offset                string     `json:"offset"`
	View                  string     `json:"view"`
	Sort                  []Sort     `json:"sort"`
	FilterByFormula       string     `json:"filterByFormula"`
	CellFormat            CellFormat `json:"cellFormat"`
	Fields                []string   `json:"fields"`
	ReturnFieldsByFieldId bool       `json:"returnFieldsByFieldId"`
	RecordMetadata        []string   `json:"recordMetadata"`
}

func SortMapFromQuery(q url.Values) []Sort {
	// sort[0][field]=Field Name
	// sort[0][direction]=desc

	fields := make(map[int]string)
	direction := make(map[int]SortDirection)
	for k, v := range q {
		if !strings.HasPrefix(k, "sort") {
			continue
		}
		split := strings.SplitN(k, "[", 3)
		if len(split) != 3 {
			continue
		}
		fieldKey, err := strconv.Atoi(strings.TrimSuffix(split[1], "]"))
		if err != nil {
			continue
		}
		fieldOrDirection := strings.TrimSuffix(split[2], "]")
		switch fieldOrDirection {
		case "field":
			fields[fieldKey] = v[0]
		case "direction":
			direction[fieldKey] = NewSortDirection(v[0])
		}
	}
	var i int
	result := make([]Sort, len(fields))
	for k, v := range fields {
		dir, ok := direction[k]
		if !ok {
			dir = SortAscending
		}
		result[i] = Sort{Field: v, Direction: dir}
		i++
	}
	return result
}

func NewListRecordsRequest(r *http.Request) (req ListRecordsRequest, validationFailures map[string]string) {
	validationFailures = make(map[string]string)

	req.BaseID = r.PathValue("baseID")
	req.TableIDOrName = r.PathValue("tableIDOrName")

	q := r.URL.Query()

	req.TimeZone = q.Get("timeZone")
	req.UserLocale = q.Get("userLocale")
	req.PageSize = intOrDefault(q.Get("pageSize"), 100)
	req.MaxRecords = intOrDefault(q.Get("maxRecords"), -1)
	req.Offset = q.Get("offset")
	req.View = q.Get("view")
	req.Sort = SortMapFromQuery(q)
	req.FilterByFormula = q.Get("filterByFormula")
	req.CellFormat = NewCellFormat(q.Get("cellFormat"))
	req.Fields = q["fields"]
	if q.Get("returnFieldsByFieldId") != "" {
		var ok bool
		if req.ReturnFieldsByFieldId, ok = parseBool(q.Get("returnFieldsByFieldId")); !ok {
			validationFailures["returnFieldsByFieldId"] = "invalid boolean"
		}
	}
	req.RecordMetadata = q["recordMetadata"]

	return
}

type ListRecordsResponse struct {
	Offset  string   `json:"offset,omitempty"`
	Records []Record `json:"records"`
}

func intOrDefault(s string, def int) int {
	if s == "" {
		return def
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}

func parseBool(s string) (v, ok bool) {
	if s != "true" && s != "false" {
		return false, false
	}
	return s == "true", true
}

type CreateRecordsRequest struct {
	// Path parameters
	BaseID        string
	TableIDOrName string

	// Body parameters

	// If one record is created:
	Fields map[string]any `json:"fields"`

	// If multiple records are created:
	Records []map[string]any `json:"records"`

	ReturnFieldsByFieldId bool `json:"returnFieldsByFieldId"`
	TypeCast              bool `json:"typecast"`
}

type CreateRecordsResponse struct {
	// If one record is created:
	ID          string         `json:"id"`
	CreatedTime time.Time      `json:"createdTime"`
	Fields      map[string]any `json:"fields"`

	// If multiple records are created:
	Records []Record `json:"records"`
}

func NewCreateRecordsRequest(r *http.Request) (req CreateRecordsRequest, validationFailures map[string]string) {
	validationFailures = make(map[string]string)

	req.BaseID = r.PathValue("baseID")
	if req.BaseID == "" {
		validationFailures["baseID"] = "missing or invalid"
	}
	req.TableIDOrName = r.PathValue("tableIDOrName")
	if req.TableIDOrName == "" {
		validationFailures["tableIDOrName"] = "missing or invalid"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		validationFailures["body"] = "invalid json"
	}

	return
}
