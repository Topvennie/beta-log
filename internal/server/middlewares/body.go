package middlewares

import (
	"encoding/json"
	"fmt"
	"mime/multipart"

	"github.com/Topvennie/beta-log/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

const maxStringLength = 200

func BodyCapture(c fiber.Ctx) error {
	body := c.Body()
	if len(body) == 0 {
		return c.Next()
	}

	contentType := string(c.Request().Header.ContentType())
	data := map[string]interface{}{}

	switch contentType {
	case "multipart/form-data":
		form, err := c.MultipartForm()
		if form == nil || err != nil {
			break
		}

		// Add normal values
		for k, v := range form.Value {
			data[k] = utils.SliceMap(v, truncate)
		}

		// Add file values
		for k, files := range form.File {
			data[k] = utils.SliceMap(files, func(f *multipart.FileHeader) string { return fmt.Sprintf("%s (%d MB)", f.Filename, f.Size/1000000) })
		}

	case "application/json":
		if err := json.Unmarshal(body, &data); err != nil {
			break
		}
	}

	if len(data) > 0 {
		c.Locals("body", data)
	}

	return c.Next()
}

func truncate(value string) string {
	if len(value) < maxStringLength {
		return value
	}

	return value[:maxStringLength-3] + "..."
}
