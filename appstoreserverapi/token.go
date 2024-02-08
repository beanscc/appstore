package appstoreserverapi

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	conf *Config

	// token lock
	mutex sync.RWMutex
	token string

	// token 过期时间
	expireAt time.Time
}

func NewToken(conf *Config) *Token {
	return &Token{
		conf: conf,
	}
}

func (t *Token) Get() (string, error) {
	if token := t.get(); "" != token {
		return token, nil
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	var err error
	t.token, t.expireAt, err = t.create()
	if err != nil {
		return "", err
	}
	return t.token, nil
}

// get return exist valid token
func (t *Token) get() string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	if "" != t.token && t.expireAt.After(time.Now()) {
		return t.token
	}
	return ""
}

// create got new token
// 文档：https://developer.apple.com/documentation/appstoreserverapi/generating_tokens_for_api_requests
func (t *Token) create() (string, time.Time, error) {
	now := time.Now()
	alg := jwt.SigningMethodES256
	exp := now.Add(55 * time.Minute)
	token := jwt.Token{
		Method: alg,
		Header: map[string]interface{}{
			"alg": alg.Name,
			"kid": t.conf.KeyID,
			"typ": "JWT",
		},
		Claims: jwt.MapClaims{
			"iss": t.conf.Issuer,
			"iat": now.Unix(),
			"exp": exp.Unix(),
			"aud": "appstoreconnect-v1",
			"bid": t.conf.BundleID,
		},
	}

	key, err := t.loadPrivateKey()
	if err != nil {
		return "", time.Time{}, err
	}

	s, err := token.SignedString(key)
	if err != nil {
		return "", time.Time{}, err
	}

	return s, exp, nil
}

func (t *Token) loadPrivateKey() (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(t.conf.PrivateKey)
	if block == nil {
		return nil, errors.New("appstore.serverapi.Token: private api key must be a PEM encoded PKCS8 key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pk, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("appstore.serverapi.Token: key is not a valid ECDSA private key")
	}

	return pk, nil
}
