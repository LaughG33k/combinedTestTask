package handler

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/LaughG33k/authServiceTestTask/iternal"
	"github.com/goccy/go-json"
)

type SharePublicKey struct {
	Path               string
	Method             string
	key                []byte
	KeyTimelifeInCache time.Duration
	sync.Mutex
}

func (h *SharePublicKey) Handle(w http.ResponseWriter, r *http.Request) {

	iternal.Logger.Infof("start sharing key for %s", r.RemoteAddr)

	key := h.getKey()

	if key == nil {
		h.setAutoDelete()
		if err := h.readFileSetToCache(); err != nil {
			iternal.Logger.Infof("failed to read key. Error: %s", err)
			http.Error(w, "Iternal error", 500)
		}

		json.NewEncoder(w).Encode(map[string]any{
			"key":    h.getKey(),
			"method": h.Method,
		})

		fmt.Println(h.Method)

		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"key":    key,
		"method": h.Method,
	})

	iternal.Logger.Infof("sharing key compleyed for %s", r.RemoteAddr)

}

func (h *SharePublicKey) getKey() []byte {
	h.Lock()
	defer h.Unlock()

	return h.key

}

func (h *SharePublicKey) setAutoDelete() {

	h.Lock()
	defer h.Unlock()

	if h.key != nil {
		return
	}

	go func() {
		time.Sleep(h.KeyTimelifeInCache)
		h.Lock()
		h.key = nil
		h.Unlock()
	}()

}

func (h *SharePublicKey) readFileSetToCache() error {

	h.Lock()
	defer h.Unlock()

	if h.key != nil {
		return nil
	}

	bytes, err := os.ReadFile(h.Path)

	if err != nil {
		return err
	}

	h.key = bytes

	return nil

}
