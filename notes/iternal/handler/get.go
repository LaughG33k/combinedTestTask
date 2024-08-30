package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/LaughG33k/notes/iternal"
	"github.com/LaughG33k/notes/iternal/repository"
)

type GetNotes struct {
	Ctx       context.Context
	NotesRepo *repository.Note
	Timeout   time.Duration
}

func (h *GetNotes) Handle(w http.ResponseWriter, r *http.Request) {

	iternal.Logger.Infof("handle get req from %s", r.RemoteAddr)

	uuid := r.Header.Get("User-Uuid")

	iternal.Logger.Infof("get uuid: %s from %s", uuid, r.RemoteAddr)

	if uuid == "" {
		iternal.Logger.Infof("uuid is invalid from %s", r.RemoteAddr)
		http.Error(w, "Uuid can't be empty", 400)
		return
	}

	tm, canc := context.WithTimeout(h.Ctx, h.Timeout)
	defer canc()

	iternal.Logger.Infof("start started receiving notes for %s", r.RemoteAddr)
	notes, err := h.NotesRepo.Get(tm, uuid)

	if err != nil {
		iternal.Logger.Infof("failed receiving notes for %s. Error: %s", r.RemoteAddr, err)
		http.Error(w, "Iternal error", 500)
		return
	}

	if len(notes) == 0 {
		iternal.Logger.Infof("notes not found for %s", r.RemoteAddr)
		http.Error(w, "notes not found", 404)
		return
	}

	iternal.Logger.Infof("notes sent to %s", r.RemoteAddr)
	json.NewEncoder(w).Encode(notes)

}
