package auth

import (
	"auth-service/pkg/hash"
	"errors"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

var (
	ErrSecretKeyIsEmpty = errors.New("secret key is empty")
)

type TokenManager interface {
	GenerateSessionToken(hash hash.Hasher) (string, error)
}

type Manager struct {
	secretKey string
}

func NewManager(secretKey string) (*Manager, error) {
	if secretKey == "" {
		return nil, ErrSecretKeyIsEmpty
	}
	return &Manager{
		secretKey: secretKey,
	}, nil
}

func (m *Manager) GenerateSessionToken(hasher hash.Hasher) (string, error) {
	session := generateRandomString(24)
	session = session + m.secretKey
	return hasher.Hash(session)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder

	rand.Seed(uint64(time.Now().UnixNano()))

	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charset))
		sb.WriteByte(charset[randomIndex])
	}

	return sb.String()
}
