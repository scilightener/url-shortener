package bl

import (
	"math/rand"
	"time"
)

const (
	aliasLength   = 7
	ValidDuration = time.Hour * 24 * 30
)

type UrlChecker interface {
	GetUrl(alias string) (string, error)
}

func GenerateUniqueAlias(checker UrlChecker) string {
	alias := getRandomString(aliasLength)
	for {
		_, err := checker.GetUrl(alias)
		if err != nil {
			break
		}
		alias = getRandomString(aliasLength)
	}
	return alias
}

func getRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetValidUntilUTC() time.Time {
	return time.Now().Add(ValidDuration).UTC()
}
