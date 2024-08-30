package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/LaughG33k/authServiceTestTask/iternal"
	customerrors "github.com/LaughG33k/authServiceTestTask/iternal/errors"
	"github.com/LaughG33k/authServiceTestTask/iternal/model"
	"github.com/LaughG33k/authServiceTestTask/iternal/repository"
	"github.com/LaughG33k/authServiceTestTask/pkg"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt"
)

type LoginHandler struct {
	UserRepo     *repository.UserRepository
	RefreshRepo  *repository.RefreshSessionRepository
	Ctx          context.Context
	Timeout      time.Duration
	JwtGenerator pkg.JwtGenerator
}

func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {

	iternal.Logger.Infof("start login for %s", r.RemoteAddr)

	loginModel := model.LogingModel{}

	if err := json.NewDecoder(r.Body).Decode(&loginModel); err != nil {
		iternal.Logger.Infof("bad request for login from %s", r.RemoteAddr)
		http.Error(w, customerrors.BadRequst.Error(), 400)
		return
	}

	if err := pkg.ValidateLoginPassword(loginModel.Login, loginModel.Password); err != nil {
		iternal.Logger.Infof("bad args from %s for login", r.RemoteAddr)
		http.Error(w, err.Error(), 400)
		return
	}

	rt := pkg.GenerateRT(30)

	tm, canc := context.WithTimeout(h.Ctx, h.Timeout)

	defer canc()

	exists, uuid, err := h.UserRepo.CheckLP(tm, loginModel.Login, loginModel.Password)

	if err != nil {
		if err == customerrors.UserNotFound {
			iternal.Logger.Info("user not found for login request %s", r.RemoteAddr)
			http.Error(w, err.Error(), 404)
			return
		}

		iternal.Logger.Infof("cannot check user by login and password for %s. Error: %s", r.RemoteAddr, err)

		http.Error(w, "Iternal error", 500)
		return
	}

	if !exists {
		iternal.Logger.Infof("incorect pass for")
		http.Error(w, "Incorect password", 400)
		return
	}

	accessToken, err := h.JwtGenerator.NewToken(jwt.MapClaims{
		"uuid": uuid,
		"ip":   r.RemoteAddr,
	})

	if err != nil {
		iternal.Logger.Infof("cannot generate jwt for %s. Error: %s", r.RemoteAddr, err)
		http.Error(w, "Iternal error", 500)
		return
	}

	if err := h.RefreshRepo.Create(tm, uuid, strings.Split(r.RemoteAddr, ":")[0], rt, time.Now().Add(24*time.Hour*30).Unix()); err != nil {
		iternal.Logger.Infof("cannot create session for %s. Error: %s", r.RemoteAddr, err)
		http.Error(w, "Iternal error", 500)
		return
	}

	tokensModel := model.TokensModel{
		Jwt:     accessToken,
		Refresh: rt,
	}

	iternal.Logger.Infof("login for %s success", r.RemoteAddr)

	json.NewEncoder(w).Encode(tokensModel)

}
