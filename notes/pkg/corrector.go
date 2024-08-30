package pkg

import (
	"strings"

	"github.com/LaughG33k/notes/iternal/model"
)

func Correct(str string, in []model.CorrectedWord) string {

	res := ""

	wordsSlice := strings.Split(str, " ")
	words := make(map[string]int, len(wordsSlice))

	for i, v := range wordsSlice {

		words[v] = i

	}

	for _, v := range in {

		if pos, ok := words[v.Word]; ok {
			if len(v.Suggestions) > 0 {
				wordsSlice[pos] = v.Suggestions[0]
			}
		}

	}

	for _, v := range wordsSlice {
		res += v + " "
	}
	return res[:len(res)-1]

}
