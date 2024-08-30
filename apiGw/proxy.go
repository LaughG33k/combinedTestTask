package apigw

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

type Proxy struct {
	t *http.Transport
	p *JwtParser
}

func InitProxy(cfg Cfg, parser *JwtParser) *Proxy {

	lAddr, err := net.ResolveTCPAddr("tcp", cfg.Addr)

	if err != nil {
		panic(err)
	}

	return &Proxy{
		t: &http.Transport{
			MaxIdleConns: 100,
			DialContext: (&net.Dialer{
				LocalAddr: &net.TCPAddr{
					IP: lAddr.IP,
				},
				Timeout: time.Duration(cfg.TimeoutInSec) * time.Second,
			}).DialContext,

			IdleConnTimeout:     time.Duration(cfg.IdleConnTimeoutInSec) * time.Second,
			MaxConnsPerHost:     cfg.MaxIdleConns,
			MaxIdleConnsPerHost: cfg.MaxIdleConns,
		},
		p: parser,
	}

}

func (p *Proxy) Proxy(w http.ResponseWriter, r *http.Request, destHost string, verifConn bool) error {

	Logger.Info(fmt.Sprintf("Start proxing from %s to %s", r.RemoteAddr, destHost+r.URL.Path))

	claims, err := p.verificateJwt(r)

	if verifConn {

		if err != nil {
			http.Error(w, err.Error(), 401)
			return err
		}

		r.Header.Set("User-Uuid", claims["uuid"].(string))
	}

	client := &http.Client{}
	r.URL.Host = destHost
	r.Host = destHost

	req := &http.Request{
		Method:           r.Method,
		URL:              &url.URL{Scheme: "http", Host: destHost, Path: r.URL.Path, ForceQuery: r.URL.ForceQuery, RawQuery: r.URL.RawQuery, RawFragment: r.URL.RawFragment},
		Header:           r.Header,
		ContentLength:    r.ContentLength,
		TransferEncoding: r.TransferEncoding,
		Form:             r.Form,
		PostForm:         r.PostForm,
		TLS:              r.TLS,
		Proto:            r.Proto,
		ProtoMajor:       r.ProtoMajor,
		ProtoMinor:       r.ProtoMinor,
		Body:             r.Body,
	}

	resp, err := client.Do(req)

	Logger.Infof("sent the req from %s to %s ", r.RemoteAddr, destHost+r.URL.Path)

	if err != nil {
		Logger.Infof("failed send the req from %s to %s. Error: %s ", r.RemoteAddr, destHost+r.URL.Path, err)
		return err
	}

	defer resp.Body.Close()

	for i, v := range resp.Header {

		for _, vals := range v {
			Logger.Log(logrus.DebugLevel, fmt.Sprintf("set to header response for %s", r.RemoteAddr))
			w.Header().Set(i, vals)
		}

	}

	bytes, err := io.ReadAll(resp.Body)

	if err != nil {
		Logger.Infof("failied reading body from %s to %s. Error: ", destHost, r.RemoteAddr, err)
		return err
	}

	Logger.Infof("proxing data from %s to %s completed", r.RemoteAddr, destHost)
	Logger.Log(logrus.DebugLevel, fmt.Sprintf("status %d and content length %d", resp.StatusCode, len(bytes)))

	w.WriteHeader(resp.StatusCode)
	w.Write(bytes)
	return nil

}

func (p *Proxy) verificateJwt(r *http.Request) (jwt.MapClaims, error) {

	Logger.Infof("start authorization for %s", r.RemoteAddr)

	jwt := r.Header.Get("Authorization")

	claims, err := p.p.ParseToken(jwt)

	if err != nil {
		Logger.Infof("failed authorization for %s. Error: %s", r.RemoteAddr, err)
		return nil, errors.New("not authorized")
	}

	Logger.Infof("authorization success for %s", r.RemoteAddr)

	return claims, nil
}
