package apigw

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

type JwtParser struct {
	signMethod string
	key        any
}

func NewParser(key []byte, signMethod string) (*JwtParser, error) {

	p := &JwtParser{
		signMethod: signMethod,
		key:        key,
	}

	if signMethod[0] == 'R' && signMethod[1] == 'S' {

		pk, err := jwt.ParseRSAPublicKeyFromPEM(key)

		if err != nil {
			return nil, err
		}

		p.key = pk

	}

	return p, nil

}

func (g *JwtParser) ParseToken(token string) (jwt.MapClaims, error) {

	res := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, res, func(t *jwt.Token) (interface{}, error) {
		if g.signMethod != t.Method.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s != %s", t.Method.Alg(), g.signMethod)
		}
		return g.key, nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil

}
