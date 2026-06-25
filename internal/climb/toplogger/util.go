package toplogger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Topvennie/beta-log/internal/database/model"
)

func parseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

var noResp = &struct{}{}

func (c *Client) request(ctx context.Context, setting model.Setting, method, url string, body io.Reader, target any) error {
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", baseURL, url), body)
	if err != nil {
		return fmt.Errorf("new http request %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+setting.ClimbToploggerAuthToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do http request %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return fmt.Errorf("wrong status code %s", resp.Status)
	}

	if target != noResp {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("decode body to json %w", err)
		}
	}

	return nil
}
