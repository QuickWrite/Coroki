package security

import (
	"errors"
	"fmt"
	"strings"

	"github.com/QuickWrite/Coroki/internal/data"
)

type passwordService struct {
	hashers map[string]data.PasswordHasher
	current *data.PasswordHasher
}

func NewPasswordService(
	current data.PasswordHasher,
	hashers ...data.PasswordHasher,
) data.PasswordService {
	m := make(map[string]data.PasswordHasher, len(hashers))

	for _, hasher := range hashers {
		m[hasher.Name()] = hasher
	}

	return passwordService{
		hashers: m,
		current: &current,
	}
}

func (s passwordService) Hash(password string) (string, error) {
	hash, err := (*s.current).Hash(password)
	if err != nil {
		return "", err
	}

	return formatEncodedHash((*s.current).Name(), hash), nil
}

func (s passwordService) Compare(password, encodedHash string) error {
	algorithm, hash, err := parseEncodedHash(encodedHash)
	if err != nil {
		return err
	}

	hasher, ok := s.hashers[algorithm]
	if !ok {
		return fmt.Errorf("unsupported password algorithm %q", algorithm)
	}

	return hasher.Compare(password, hash)
}

func (s passwordService) NeedsRehash(encodedHash string) bool {
	algorithm, hash, err := parseEncodedHash(encodedHash)
	if err != nil {
		return true
	}

	// If the algorithm isn't supported anymore, the password
	// should be considered in need of rehashing.
	hasher, ok := s.hashers[algorithm]
	if !ok {
		return true
	}

	// We always want passwords on the current algorithm.
	if algorithm != (*s.current).Name() {
		return true
	}

	return hasher.NeedsRehash(hash)
}

func formatEncodedHash(algorithm, hash string) string {
	return "{" + algorithm + "}" + hash
}

func parseEncodedHash(encoded string) (algorithm, hash string, err error) {
	if !strings.HasPrefix(encoded, "{") {
		return "", "", errors.New("password hash has no algorithm")
	}

	end := strings.IndexByte(encoded, '}')
	if end == -1 {
		return "", "", errors.New("invalid password hash format")
	}

	algorithm = encoded[1:end]
	hash = encoded[end+1:]

	if algorithm == "" {
		return "", "", errors.New("password hash has empty algorithm")
	}

	if hash == "" {
		return "", "", errors.New("password hash has empty hash")
	}

	return algorithm, hash, nil
}
