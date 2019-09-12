package text

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"image"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const (
	apiEndpoint  = "https://api.ocr.space/parse/image"
	apiKey       = "30b85f0dda88957"
	apiKeyHeader = "apiKey"
)

type (
	apiResponse struct {
		HasError     bool     `json:"IsErroredOnProcessing"`
		ErrorMessage []string `json:"ErrorMessage"`
		Result       []struct {
			ErrorMessage string `json:"ErrorMessage"`
			Overlay      struct {
				Lines []struct {
					Words []struct {
						Left   float32
						Top    float32
						Height float32
						Width  float32
					}
				} `json:"Lines"`
			} `json:"TextOverlay"`
		} `json:"ParsedResults"`
	}
)

func Find(filePath string) ([]image.Rectangle, error) {
	response, err := apiCall(filePath)
	if err != nil {
		return nil, err
	}

	if len(response.Result[0].Overlay.Lines) == 0 {
		return nil, nil
	}

	var result []image.Rectangle
	for _, l := range response.Result[0].Overlay.Lines {
		for _, w := range l.Words {
			result = append(result, image.Rectangle{
				Min: image.Point{
					X: int(w.Left),
					Y: int(w.Top),
				},
				Max: image.Point{
					X: int(w.Left + w.Width),
					Y: int(w.Top + w.Height),
				},
			})
		}
	}
	return result, nil
}

func apiCall(filePath string) (response apiResponse, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		err = errors.Wrap(err, "failed to open file")
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Errorf("failed to close file %q, %q", filePath, err.Error())
		}
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		err = errors.Wrap(err, "failed to create writer")
		return
	}

	if _, err = io.Copy(part, file); err != nil {
		err = errors.Wrap(err, "failed to Copy data from file")
		return
	}

	if err = writer.WriteField("isOverlayRequired", "True"); err != nil {
		err = errors.Wrap(err, "failed to WriteField")
		return
	}

	if err = writer.Close(); err != nil {
		err = errors.Wrap(err, "failed to Close writer")
		return
	}

	req, err := http.NewRequest("POST", apiEndpoint, body)
	if err != nil {
		err = errors.Wrap(err, "failed to create POST-request")
		return
	}

	req.Header.Set(apiKeyHeader, apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.Wrap(err, "failed to Do request")
		return
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf(" failed to Close response body: %q", err.Error())
		}
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "failed to Read response body")
		return
	}

	if err = json.Unmarshal(data, &response); err != nil {
		err = errors.Wrap(err, "failed to Unmarshal api response")
		return
	}

	if response.HasError {
		err = errors.Errorf("got API error: %q", response.ErrorMessage[0])
		return
	}

	return
}
