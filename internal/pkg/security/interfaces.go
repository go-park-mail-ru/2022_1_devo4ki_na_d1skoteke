package security

type Manager interface {
	Hash(string) string
	ComparePasswords(string, string) error
}
