package pkg

import "math/rand"

var symbols []rune = []rune("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPSADFGHJKLZXCVBNM1234567890")

func GenerateRandomString(length int) string {

	res := ""

	for i := 0; i < length; i++ {
		res += string(symbols[rand.Intn(len(symbols))])
	}

	return res

}
