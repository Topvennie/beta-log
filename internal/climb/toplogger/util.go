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

func parseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

func (c *Client) request(ctx context.Context, token, method, url string, body io.Reader) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", baseURL, url), body)
	if err != nil {
		return nil, 0, fmt.Errorf("new http request %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
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
	setting.ClimbTopLoggerExpiration = time.Time{}

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

	switch result[0].Errors[0].Extension.OriginalError.StatusCode {
	case 401:
		return ErrUnauthorized
	case 403:
		return ErrUnauthorized
	default:
		return errors.New(result[0].Errors[0].Message)
	}
}
