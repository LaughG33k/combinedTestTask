package yandexspeller

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/LaughG33k/notes/iternal/model"
)

const serviceURL = "http://speller.yandex.net/services/spellservice.json/checkText"

func CheckText(text string, timeout time.Duration) ([]model.CorrectedWord, error) {

	cl := &http.Client{Timeout: timeout}

	resp, err := cl.PostForm(serviceURL, url.Values{
		"text": {text},
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var speled []model.CorrectedWord

	json.Unmarshal(res, &speled)

	return speled, nil

}
