package apigw

import (
	"net/http"
)

type EndPoint struct {
	Paths         []string `yaml:"paths"`
	Addr          string   `yaml:"addr"`
	OnlyVerifConn bool     `yaml:"onlyVerifConn"`
}

func GenerateHandlers(proxy *Proxy, eps []EndPoint) {

	for _, v := range eps {

		for _, path := range v.Paths {

			addr := v.Addr
			verifConn := v.OnlyVerifConn

			Logger.Infof("started handler listening addr: %s and path: %s", v.Addr, path)
			http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				if err := proxy.Proxy(w, r, addr, verifConn); err != nil {
					Logger.Info(err)
				}
			})

		}

	}

}
