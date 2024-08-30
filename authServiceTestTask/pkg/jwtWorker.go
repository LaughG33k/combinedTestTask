package pkg

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtGenerator interface {
	NewToken(claims jwt.MapClaims) (string, error)
}

type JwtParser interface {
	ParseToken(token string) (jwt.MapClaims, error)
}

type JwtWorker struct {
	key           any
	singMethod    jwt.SigningMethod
	TokenTimelife time.Duration
}

func NewJwtWorker(key []byte, signMethod jwt.SigningMethod) (*JwtWorker, error) {

	res := &JwtWorker{
		singMethod: signMethod,
		key:        key,
	}

	if _, ok := signMethod.(*jwt.SigningMethodRSA); ok {

		pk, errPriv := jwt.ParseRSAPrivateKeyFromPEM(key)
		pub, errPub := jwt.ParseRSAPublicKeyFromPEM(key)

		if errPriv == nil {
			res.key = pk
			return res, nil
		}

		if errPub == nil {
			res.key = pub
			return res, nil
		}

		return nil, errors.New("not rsa key")

	}

	return res, nil
}

func (g *JwtWorker) NewToken(claims jwt.MapClaims) (string, error) {

	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(g.TokenTimelife).Unix()

	token := jwt.NewWithClaims(g.singMethod, claims)

	jwt, err := token.SignedString(g.key)

	if err != nil {
		return "", err
	}

	return jwt, nil

}

func (g *JwtWorker) ParseToken(token string) (jwt.MapClaims, error) {

	res := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, res, func(t *jwt.Token) (interface{}, error) {
		if g.singMethod.Alg() != t.Method.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s != %s", t.Method.Alg(), g.singMethod.Alg())
		}
		return g.key, nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil

}
