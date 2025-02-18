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

func doRequest[TResponse any](ctx context.Context, log *slog.Logger, client *http.Client, method string, url string, pat string, body any) (response TResponse, err error) {
	log.Debug("sending request", slog.String("url", url))

	var br io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return response, err
		}
		br = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, br)
	if err != nil {
		return response, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pat))
	if br != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rsp, err := client.Do(req)
	if err != nil {
		return response, err
	}

	err = parseResponse(log, rsp, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func parseResponse(log *slog.Logger, response *http.Response, v any) error {
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, response.Body)
	if err != nil {
		return err
	}

	log.Debug("HTTP response", slog.Int("status", response.StatusCode), slog.String("body", buf.String()))

	return json.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(v)
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

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, rsp.Body)
	if err != nil {
		return models.ListRecordsResponse{}, err
	}

	c.Log.Debug("HTTP response", slog.Int("status", rsp.StatusCode), slog.String("body", buf.String()))

	err = json.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(&response)
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

	return doRequest[models.CreateRecordsResponse](ctx, c.Log, c.Client, http.MethodPost, u.String(), c.PAT, request)
}

func (c *Client) DeleteMultiple(ctx context.Context, request models.DeleteRecordsRequest) (response models.DeleteRecordsResponse, err error) {

	u, err := url.Parse(fmt.Sprintf("https://api.airtable.com/v0/%s/%s", url.PathEscape(request.BaseID), url.PathEscape(request.TableIDOrName)))
	if err != nil {
		return models.DeleteRecordsResponse{}, err
	}

	return doRequest[models.DeleteRecordsResponse](ctx, c.Log, c.Client, http.MethodDelete, u.String(), c.PAT, nil)
}
