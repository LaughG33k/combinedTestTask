package apigw

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/goccy/go-json"
)

func GetKey(url string, timeout time.Duration) (PublicKey, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return PublicKey{}, err
	}

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(req)

	if err != nil {
		return PublicKey{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return PublicKey{}, errors.New("Failed to get public key")
	}

	bytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return PublicKey{}, err
	}

	var model PublicKey

	if err := json.Unmarshal(bytes, &model); err != nil {
		return PublicKey{}, err
	}

	return model, nil

}
