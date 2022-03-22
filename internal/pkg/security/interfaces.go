package security

type Manager interface {
	Hash(string) []byte
	ComparePasswords(string, string) error
}
