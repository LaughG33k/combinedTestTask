package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	yandexspeller "github.com/LaughG33k/notes/client/yandexSpeller"
	"github.com/LaughG33k/notes/iternal"
	"github.com/LaughG33k/notes/iternal/model"
	"github.com/LaughG33k/notes/iternal/repository"
	"github.com/LaughG33k/notes/pkg"
	"github.com/goccy/go-json"
)

type CreateNote struct {
	NoteRepository *repository.Note
	Ctx            context.Context
	Timeout        time.Duration
}

func (h *CreateNote) Handle(w http.ResponseWriter, r *http.Request) {

	iternal.Logger.Infof("handler create request from %s", r.RemoteAddr)

	note := model.Note{OwnerUuid: r.Header.Get("User-Uuid")}

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	iternal.Logger.Debug(fmt.Sprintf("handler model note from %s. Params ", r.RemoteAddr), note)

	tm, canc := context.WithTimeout(h.Ctx, h.Timeout)
	defer canc()

	iternal.Logger.Infof("check Text start for %s", r.RemoteAddr)
	correctedWords, err := yandexspeller.CheckText(note.Content, h.Timeout)

	if err != nil {
		http.Error(w, "failed to fix errors. Try again", 500)
		iternal.Logger.Infof("failed to fix errors for %s request. Error: %s", r.RemoteAddr, err)
		return
	}

	note.Content = pkg.Correct(note.Content, correctedWords)

	if err := h.NoteRepository.Create(tm, note); err != nil {

		if err == context.DeadlineExceeded {
			http.Error(w, "timeout error", 500)
			iternal.Logger.Infof("failed to save note for %s. Error: %s", r.RemoteAddr, err)
			return
		}
		http.Error(w, "Iternal error", 500)
		iternal.Logger.Infof("failed to save note for %s. Error: %s", r.RemoteAddr, err)
		return
	}

	iternal.Logger.Infof("note success created for %s", r.RemoteAddr)

}
