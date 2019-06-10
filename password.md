## Password Hashing with Argon2

```go
package main

// REFERENCES: https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("invalid hash format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

type Config struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func NewConfig() Config {
	return Config{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
}

type PasswordManager struct {
	config Config
}

func NewPasswordManager(cfg Config) *PasswordManager {
	return &PasswordManager{
		config: cfg,
	}
}

func (p *PasswordManager) Hash(password string) (string, error) {
	if trimmed := strings.TrimSpace(password); len(trimmed) < 8 {
		return "", errors.New("password is required")
	}

	salt, err := generateRandomBytes(p.config.saltLength)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		p.config.iterations,
		p.config.memory,
		p.config.parallelism,
		p.config.keyLength,
	)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		p.config.memory,
		p.config.iterations,
		p.config.parallelism,
		b64Salt,
		b64Hash,
	)
	return encodedHash, nil
}
func (p *PasswordManager) Compare(password, encodedHash string) (bool, error) {
	cfg, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}
	otherHash := argon2.IDKey([]byte(password), salt, cfg.iterations, cfg.memory, cfg.parallelism, cfg.keyLength)
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (cfg *Config, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}
	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}

	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}
	cfg = &Config{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &cfg.memory, &cfg.iterations, &cfg.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}
	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	cfg.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	cfg.keyLength = uint32(len(hash))
	return cfg, salt, hash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func main() {
	pwdMgr := NewPasswordManager(NewConfig())
	f := func(s string) bool {
		hash, err := pwdMgr.Hash(s)
		if err != nil {
			return false
		}
		match, err := pwdMgr.Compare(s, hash)
		if err != nil {
			return false
		}
		return match
	}
	log.Println(f("caramel10%"))
	/*if err := quick.Check(f, nil); err != nil {
		log.Fatal(err)
	}*/
}
```
