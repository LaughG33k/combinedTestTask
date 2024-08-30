package pkg

import "encoding/base64"

func GenerateRT(length int) string {
	return base64.StdEncoding.EncodeToString([]byte(GenerateRandomString(length)))
}
