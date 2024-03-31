package bl

import (
	"context"
	"math/rand"
	"time"
)

const (
	aliasLength   = 7
	ValidDuration = time.Hour * 24 * 30
)

type UrlChecker interface {
	GetUrl(ctx context.Context, alias string) (string, error)
}

func GenerateUniqueAlias(ctx context.Context, checker UrlChecker) string {
	alias := getRandomString(aliasLength)
	for {
		_, err := checker.GetUrl(ctx, alias)
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
