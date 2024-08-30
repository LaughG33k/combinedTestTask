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
	"github.com/golang-jwt/jwt"
)

type UpdateTokensHandler struct {
	RerfreshRepo *repository.RefreshSessionRepository
	Ctx          context.Context
	Timeout      time.Duration
	JwtParser    pkg.JwtParser
	JwtGenerator pkg.JwtGenerator
}

func (h *UpdateTokensHandler) Handle(w http.ResponseWriter, r *http.Request) {

	tokens := model.TokensModel{}
	iternal.Logger.Infof("start update sessions for %s", r.RemoteAddr)

	if err := json.NewDecoder(r.Body).Decode(&tokens); err != nil {
		iternal.Logger.Infof("decode jwt failed for %s. Error: %s", r.RemoteAddr, err)
		http.Error(w, "Bad request", 400)
		return
	}

	parsedJwt, err := h.JwtParser.ParseToken(tokens.Jwt)

	if err != nil {
		iternal.Logger.Infof("decode jwt failed for %s. Error: %s", r.RemoteAddr, err)
		http.Error(w, "decode jwt failed", 401)
		return
	}

	rt := pkg.GenerateRT(30)
	accessToken, err := h.JwtGenerator.NewToken(jwt.MapClaims{
		"uuid": parsedJwt["uuid"],
		"ip":   r.RemoteAddr,
	})

	if err != nil {
		iternal.Logger.Infof("failed creating new jwt for %s. Error: %s", r.RemoteAddr, err)
		http.Error(w, "Iternal error", 500)
		return
	}

	tm, canc := context.WithTimeout(h.Ctx, h.Timeout)
	defer canc()

	ip, timelife, err := h.RerfreshRepo.Get(tm, parsedJwt["uuid"].(string), tokens.Refresh)

	if err != nil {

		if err == customerrors.RefreshNotFound {
			http.Error(w, err.Error(), 401)
			return
		}

		iternal.Logger.Infof("failed getting user for %s. Error: %s", r.RemoteAddr, err)
		http.Error(w, "Iternal error", 500)
		return

	}

	if ip != parsedJwt["ip"] {
		http.Error(w, "refresh token not supports the jwt", 401)
		return
	}

	if time.Now().Unix() > timelife {
		http.Error(w, "Token has exipired", 401)
		return
	}

	if parsedJwt["ip"] != r.RemoteAddr {
		// todo send notify to email

		if err := h.RerfreshRepo.Delete(tm, parsedJwt["uuid"].(string), tokens.Refresh); err != nil {
			iternal.Logger.Infof("failed delete session for %s. Error: %s", r.RemoteAddr, err)
			http.Error(w, "Iternal error", 500)
			return
		}

		if err := h.RerfreshRepo.Create(tm, parsedJwt["uuid"].(string), r.RemoteAddr, rt, time.Now().Add(24*time.Hour*30).Unix()); err != nil {
			iternal.Logger.Infof("failed create session for %s. Error: %s", r.RemoteAddr, err)
			http.Error(w, "Iternal error", 500)
			return
		}

	} else {

		err := h.RerfreshRepo.Update(tm, parsedJwt["uuid"].(string), r.RemoteAddr, tokens.Refresh, rt, time.Now().Add(24*time.Hour*30).Unix())

		if err != nil {
			if err == customerrors.RefreshNotFound {
				http.Error(w, err.Error(), 401)
				return
			}

			iternal.Logger.Infof("failed update session for %s. Error: %s", r.RemoteAddr, err)

			http.Error(w, "Iternal error", 500)
			return
		}

	}

	tokens.Jwt = accessToken
	tokens.Refresh = rt

	json.NewEncoder(w).Encode(tokens)

	iternal.Logger.Infof("update sessions sucess for %s", r.RemoteAddr)

}
