package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/LaughG33k/authServiceTestTask/iternal"
	customerrors "github.com/LaughG33k/authServiceTestTask/iternal/errors"
	"github.com/LaughG33k/authServiceTestTask/iternal/model"
	"github.com/LaughG33k/authServiceTestTask/iternal/repository"
	"github.com/LaughG33k/authServiceTestTask/pkg"
	"github.com/goccy/go-json"
)

type RegHandler struct {
	UserRepo *repository.UserRepository
	Ctx      context.Context
	Timeout  time.Duration
}

func (h *RegHandler) Handle(w http.ResponseWriter, r *http.Request) {

	regModel := model.RegistrationModel{}

	iternal.Logger.Infof("start registration for %s", r.RemoteAddr)

	if err := json.NewDecoder(r.Body).Decode(&regModel); err != nil {
		iternal.Logger.Infof("failed decode json from %s. Error: %s", r.RemoteAddr, err)
		http.Error(w, customerrors.BadRequst.Error(), 400)
		return
	}

	if err := pkg.ValidateRegistration(regModel); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	tm, canc := context.WithTimeout(h.Ctx, h.Timeout)

	defer canc()

	if err := h.UserRepo.CreateUser(tm, regModel); err != nil {
		if err == customerrors.UserAlreadyExists {
			http.Error(w, err.Error(), 404)
			return
		}

		iternal.Logger.Infof("Failed create new user for %s. Error: %s", r.RemoteAddr, err)
		http.Error(w, "Iternal error", 500)
		return
	}

	iternal.Logger.Infof("registration completed for %s", r.RemoteAddr)
}
