package util

import "errors"

var tokenMap = map[string]int64{
	"token1": 1,
	"token2": 2,
	"token3": 3,
	"token4": 4,
	"token5": 5,
	"token6": 6,
	"token7": 7,
	"token8": 8,
}

var EmptyToken = errors.New("empty token")

func VerifyToken(token string) (uid int64, err error) {
	if val, ok := tokenMap[token]; ok {
		return val, nil
	}
	return 0, EmptyToken
}
