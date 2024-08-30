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

type Logout struct {
	RefreshRepo *repository.RefreshSessionRepository
	JwtParser   pkg.JwtParser
	Ctx         context.Context
	Timeout     time.Duration
}

func (h *Logout) Handle(w http.ResponseWriter, r *http.Request) {

	tokens := model.TokensModel{}

	iternal.Logger.Infof("logout start for %s", r.RemoteAddr)

	if err := json.NewDecoder(r.Body).Decode(&tokens); err != nil {
		iternal.Logger.Infof("cannot decode json from %s. Error: %s", r.RemoteAddr, err)
		http.Error(w, customerrors.BadRequst.Error(), 400)
		return
	}

	parsedJwt, err := h.JwtParser.ParseToken(tokens.Jwt)

	if err != nil {
		iternal.Logger.Infof("parsing jwt failed for %s. Error: %s", r.RemoteAddr, err)
		http.Error(w, "jwt parse failed", 401)
		return
	}

	tm, canc := context.WithTimeout(h.Ctx, h.Timeout)
	defer canc()

	ip, timelife, err := h.RefreshRepo.Get(tm, parsedJwt["uuid"].(string), tokens.Refresh)

	if err != nil {
		if err == customerrors.RefreshNotFound {
			http.Error(w, err.Error(), 401)
			return
		}

		iternal.Logger.Infof("failed to get session for %s. Error: %s", r.RemoteAddr, err)
		http.Error(w, "Iternal error", 500)
		return
	}

	if parsedJwt["ip"].(string) != ip {
		http.Error(w, "refresh token not supports the jwt", 401)
		return
	}

	if time.Now().Unix() > timelife {
		http.Error(w, "Token has exipired", 401)
		return
	}

	if err := h.RefreshRepo.Delete(h.Ctx, parsedJwt["uuid"].(string), tokens.Refresh); err != nil {
		iternal.Logger.Infof("failed delete session for %s. Error: %s", r.RemoteAddr, err)
		http.Error(w, "Iternal error", 500)
		return
	}

	iternal.Logger.Infof("logout completed for %s", r.RemoteAddr)

}
