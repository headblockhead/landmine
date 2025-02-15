package airtable

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/headblockhead/landmine/models"
)

type Client struct {
	Client *http.Client
	PAT    string
}

func New(client *http.Client, PAT string) *Client {
	return &Client{
		Client: client,
		PAT:    PAT,
	}
}

func (c *Client) List(ctx context.Context, request models.ListGetRequest) (response models.ListGetResponse, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.airtable.com/v0/%s/%s", url.PathEscape(request.BaseID), url.PathEscape(request.TableIDOrName)), nil)
	if err != nil {
		return models.ListGetResponse{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.PAT))
	rsp, err := c.Client.Do(req)
	if err != nil {
		return models.ListGetResponse{}, err
	}
	defer rsp.Body.Close()

	err = json.NewDecoder(rsp.Body).Decode(&response)
	if err != nil {
		return models.ListGetResponse{}, err
	}

	return response, nil
}
