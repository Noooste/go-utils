package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/google/uuid"
	"math/rand"
	"strings"
)

const (
	hexTable = "0123456789abcdef"
)

type Oauth2 struct {
	codeVerifier  string
	codeChallenge string
}

func NewOauth2() *Oauth2 {
	return &Oauth2{}
}

var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateRandomString(length uint8) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

func (o *Oauth2) GenerateCodeVerifier(length uint8) {
	o.codeVerifier = GenerateRandomString(length)
	o.codeChallenge = GetHashOf(o.codeVerifier)
}

func (o *Oauth2) GetCodeVerifier() string {
	return o.codeVerifier
}

func (o *Oauth2) GetCodeChallenge() string {
	return o.codeChallenge
}

func (o *Oauth2) GetCodeChallengeMethod() string {
	return "S256"
}

func (o *Oauth2) GetState() string {
	return uuid.New().String()
}

func (o *Oauth2) GetNonce() string {
	return uuid.New().String()
}

func GetHexString(length int) string {
	var result = make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = hexTable[rand.Intn(15)]
	}
	return string(result)
}

func GetHashOf(text string) string {
	var result = sha256.Sum256([]byte(text))
	return strings.TrimRight(base64.URLEncoding.EncodeToString(result[:]), "=")
}

func GetHashOfWithSalt(text, salt string, timestamp int) string {
	bText := []byte(text + salt)
	for i, el := range bText {
		bText[i] = el + byte(timestamp%i)
	}
	hash := md5.Sum(bText)
	return hex.EncodeToString(hash[:])
}

func NewSHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
