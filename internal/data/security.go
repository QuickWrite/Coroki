package data

type PasswordHasher interface {
	Name() string
	Hash(password string) (string, error)
	Compare(encodedHash string, password string) error
	NeedsRehash(encodedHash string) bool
}
