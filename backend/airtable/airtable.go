package airtable

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/headblockhead/landmine/models"
)

type Client struct {
	Log    *slog.Logger
	Client *http.Client
	PAT    string
}

func New(log *slog.Logger, client *http.Client, PAT string) *Client {
	return &Client{
		Log:    log,
		Client: client,
		PAT:    PAT,
	}
}

func QueryFromSortMap(sort []models.Sort) (q url.Values) {
	q = make(url.Values)
	for i, v := range sort {
		q.Set(fmt.Sprintf("sort[%d][field]", i), string(v.Field))
		q.Set(fmt.Sprintf("sort[%d][direction]", i), string(v.Direction))
		i++
	}
	return q
}

func (c *Client) List(ctx context.Context, request models.ListRecordsRequest) (response models.ListRecordsResponse, err error) {
	u, err := url.Parse(fmt.Sprintf("https://api.airtable.com/v0/%s/%s", url.PathEscape(request.BaseID), url.PathEscape(request.TableIDOrName)))
	if err != nil {
		return models.ListRecordsResponse{}, err
	}

	q := u.Query()
	if request.TimeZone != "" {
		q.Set("timeZone", request.TimeZone)
	}
	if request.UserLocale != "" {
		q.Set("userLocale", request.UserLocale)
	}
	q.Set("pageSize", fmt.Sprintf("%d", request.PageSize))
	if request.MaxRecords != -1 {
		q.Set("maxRecords", fmt.Sprintf("%d", request.MaxRecords))
	}
	if request.Offset != "" {
		q.Set("offset", request.Offset)
	}
	if request.View != "" {
		q.Set("view", request.View)
	}
	sortValues := QueryFromSortMap(request.Sort)
	for key, value := range sortValues {
		q[key] = value
	}
	if request.FilterByFormula != "" {
		q.Set("filterByFormula", request.FilterByFormula)
	}
	q.Set("cellFormat", string(request.CellFormat))
	for _, field := range request.Fields {
		q.Add("fields", field)
	}
	if request.ReturnFieldsByFieldId {
		q.Set("returnFieldsByFieldId", "true")
	}
	for _, recordMetadata := range request.RecordMetadata {
		q.Add("recordMetadata", recordMetadata)
	}
	u.RawQuery = q.Encode()

	c.Log.Debug("sending request", slog.String("url", u.String()))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return models.ListRecordsResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.PAT))
	rsp, err := c.Client.Do(req)
	if err != nil {
		return models.ListRecordsResponse{}, err
	}
	defer rsp.Body.Close()

	err = json.NewDecoder(rsp.Body).Decode(&response)
	if err != nil {
		return models.ListRecordsResponse{}, err
	}

	return response, nil
}

func (c *Client) Create(ctx context.Context, request models.CreateRecordsRequest) (response models.CreateRecordsResponse, err error) {
	u, err := url.Parse(fmt.Sprintf("https://api.airtable.com/v0/%s/%s", url.PathEscape(request.BaseID), url.PathEscape(request.TableIDOrName)))
	if err != nil {
		return models.CreateRecordsResponse{}, err
	}

	b, err := json.Marshal(request)
	if err != nil {
		return models.CreateRecordsResponse{}, err
	}

	c.Log.Debug("sending request", slog.String("url", u.String()), slog.String("body", string(b)))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(b))
	if err != nil {
		return models.CreateRecordsResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.PAT))
	req.Header.Set("Content-Type", "application/json")
	rsp, err := c.Client.Do(req)
	if err != nil {
		return models.CreateRecordsResponse{}, err
	}
	defer rsp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, rsp.Body)
	if err != nil {
		return models.CreateRecordsResponse{}, err
	}

	c.Log.Debug("HTTP response", slog.Int("status", rsp.StatusCode), slog.String("body", buf.String()))

	err = json.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(&response)
	if err != nil {
		return models.CreateRecordsResponse{}, err
	}

	return response, nil
}
