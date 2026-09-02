package data

type PasswordHasher interface {
	Name() string
	Hash(password string) (string, error)
	Compare(password string, encodedHash string) error
	NeedsRehash(encodedHash string) bool
}
