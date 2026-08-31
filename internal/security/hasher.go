package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

// -----------------------------------------------------------------------------
// Bcrypt
// -----------------------------------------------------------------------------

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) BcryptHasher {
	return BcryptHasher{
		cost: cost,
	}
}

func (h BcryptHasher) Name() string {
	return "bcrypt"
}

func (h BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		h.cost,
	)

	if err != nil {
		return "", fmt.Errorf("hash password with bcrypt: %w", err)
	}

	return string(hash), nil
}

func (h BcryptHasher) Compare(password string, encodedHash string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(encodedHash),
		[]byte(password),
	)
}

func (h BcryptHasher) NeedsRehash(encodedHash string) bool {
	cost, err := bcrypt.Cost([]byte(encodedHash))
	if err != nil {
		return true
	}

	return cost != h.cost
}

// -----------------------------------------------------------------------------
// Argon2id
// -----------------------------------------------------------------------------

type Argon2idHasher struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	KeyLength   uint32
	SaltLength  uint32
}

func (h Argon2idHasher) Name() string {
	return "argon2id"
}

func (h Argon2idHasher) Hash(password string) (string, error) {
	if err := h.validate(); err != nil {
		return "", err
	}

	salt := make([]byte, h.SaltLength)

	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.Iterations,
		h.Memory,
		h.Parallelism,
		h.KeyLength,
	)

	// Store all parameters necessary to verify the hash.
	//
	// $argon2id$v=19$m=65536,t=3,p=4$<salt>$<hash>
	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		h.Memory,
		h.Iterations,
		h.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func (h Argon2idHasher) Compare(password string, encodedHash string) error {
	params, salt, expectedHash, err := parseArgon2idHash(encodedHash)
	if err != nil {
		return err
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Memory,
		params.Iterations,
		params.Parallelism,
		uint32(len(expectedHash)),
	)

	if subtle.ConstantTimeCompare(actualHash, expectedHash) != 1 {
		return errors.New("password does not match")
	}

	return nil
}

func (h Argon2idHasher) NeedsRehash(encodedHash string) bool {
	params, _, expectedHash, err := parseArgon2idHash(encodedHash)
	if err != nil {
		return true
	}

	return params.Memory != h.Memory ||
		params.Iterations != h.Iterations ||
		params.Parallelism != h.Parallelism ||
		uint32(len(expectedHash)) != h.KeyLength
}

func (h *Argon2idHasher) validate() error {
	if h.Memory == 0 {
		return errors.New("argon2id memory must be greater than zero")
	}

	if h.Iterations == 0 {
		return errors.New("argon2id iterations must be greater than zero")
	}

	if h.Parallelism == 0 {
		return errors.New("argon2id parallelism must be greater than zero")
	}

	if h.KeyLength == 0 {
		return errors.New("argon2id key length must be greater than zero")
	}

	if h.SaltLength == 0 {
		return errors.New("argon2id salt length must be greater than zero")
	}

	return nil
}

type argon2idParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
}

func parseArgon2idHash(
	encodedHash string,
) (argon2idParams, []byte, []byte, error) {
	var params argon2idParams

	parts := strings.Split(encodedHash, "$")

	if len(parts) != 6 {
		return params, nil, nil, errors.New("invalid argon2id hash format")
	}

	if parts[0] != "" {
		return params, nil, nil, errors.New("invalid argon2id hash format")
	}

	if parts[1] != "argon2id" {
		return params, nil, nil, errors.New("invalid argon2id algorithm")
	}

	if parts[2] != "v=19" {
		return params, nil, nil, errors.New("unsupported argon2id version")
	}

	if err := parseArgon2idParameters(parts[3], &params); err != nil {
		return params, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return params, nil, nil, fmt.Errorf("decode argon2id salt: %w", err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return params, nil, nil, fmt.Errorf("decode argon2id hash: %w", err)
	}

	if len(salt) == 0 {
		return params, nil, nil, errors.New("argon2id hash has empty salt")
	}

	if len(hash) == 0 {
		return params, nil, nil, errors.New("argon2id hash has empty hash")
	}

	return params, salt, hash, nil
}

func parseArgon2idParameters(
	value string,
	params *argon2idParams,
) error {
	values := strings.Split(value, ",")

	if len(values) != 3 {
		return errors.New("invalid argon2id parameters")
	}

	seenMemory := false
	seenIterations := false
	seenParallelism := false

	for _, value := range values {
		key, rawValue, ok := strings.Cut(value, "=")
		if !ok || rawValue == "" {
			return errors.New("invalid argon2id parameter")
		}

		switch key {
		case "m":
			if seenMemory {
				return errors.New("duplicate argon2id memory parameter")
			}

			memory, err := strconv.ParseUint(rawValue, 10, 32)
			if err != nil || memory == 0 {
				return errors.New("invalid argon2id memory parameter")
			}

			params.Memory = uint32(memory)
			seenMemory = true

		case "t":
			if seenIterations {
				return errors.New("duplicate argon2id iterations parameter")
			}

			iterations, err := strconv.ParseUint(rawValue, 10, 32)
			if err != nil || iterations == 0 {
				return errors.New("invalid argon2id iterations parameter")
			}

			params.Iterations = uint32(iterations)
			seenIterations = true

		case "p":
			if seenParallelism {
				return errors.New("duplicate argon2id parallelism parameter")
			}

			parallelism, err := strconv.ParseUint(rawValue, 10, 8)
			if err != nil || parallelism == 0 {
				return errors.New("invalid argon2id parallelism parameter")
			}

			params.Parallelism = uint8(parallelism)
			seenParallelism = true

		default:
			return fmt.Errorf("unknown argon2id parameter %q", key)
		}
	}

	if !seenMemory || !seenIterations || !seenParallelism {
		return errors.New("missing argon2id parameter")
	}

	return nil
}

// -----------------------------------------------------------------------------
// Plain text
// -----------------------------------------------------------------------------

// PlainTextHasher returns the input as the hash.
//
// WARNING: This is intentionally insecure and should only be used for
// debugging and tests.
type PlainTextHasher struct{}

func (h PlainTextHasher) Name() string {
	return "plain"
}

func (h PlainTextHasher) Hash(password string) (string, error) {
	return password, nil
}

func (h PlainTextHasher) Compare(password string, encodedHash string) error {
	if subtle.ConstantTimeCompare(
		[]byte(encodedHash),
		[]byte(password),
	) != 1 {
		return errors.New("password does not match")
	}

	return nil
}

func (h PlainTextHasher) NeedsRehash(encodedHash string) bool {
	return false
}
