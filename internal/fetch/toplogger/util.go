package toplogger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Topvennie/beta-log/internal/database/model"
)

func (c *Client) request(ctx context.Context, token, method, url string, body io.Reader) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", baseURL, url), body)
	if err != nil {
		return nil, 0, fmt.Errorf("new http request %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("do http request %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("read body %w", err)
	}

	return respBody, resp.StatusCode, nil
}

func (c *Client) resetSetting(ctx context.Context, setting model.Setting) error {
	setting.ClimbToploggerAuthToken = ""
	setting.ClimbToploggerRefreshToken = ""
	setting.ClimbToploggerExpiration = time.Time{}

	return c.setting.ToploggerUpdate(ctx, setting)
}

func getError(data []byte) error {
	type errorResponse struct {
		Errors []cError `json:"errors"`
	}

	var result []errorResponse

	if err := json.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("unmarshal data %w", err)
	}

	if len(result) == 0 || len(result[0].Errors) == 0 {
		return nil
	}

	var errs []error
	for _, e := range result[0].Errors {
		switch e.Extension.OriginalError.StatusCode {
		case 401, 403:
			errs = append(errs, ErrUnauthorized)
		default:
			errs = append(errs, errors.New(e.Message))
		}
	}
	return errors.Join(errs...)
}
