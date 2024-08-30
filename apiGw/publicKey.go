package apigw

type PublicKey struct {
	Method string `json:"method"`
	Key    []byte `json:"key"`
}
