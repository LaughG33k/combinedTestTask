package model

type TokensModel struct {
	Jwt     string `json:"jwt"`
	Refresh string `json:"refresh"`
}
