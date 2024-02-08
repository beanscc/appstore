package jws

import (
	"encoding/json"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWS_VerifyAndBind(t *testing.T) {
	token := `App Store JSON Web Signature format token string`
	jws, err := Parse(token)
	if err != nil {
		t.Errorf("TestJWS_VerifyAndBind Parse token failed. err:%v", err)
		return
	}
	// https://developer.apple.com/documentation/appstoreserverapi/jwstransactiondecodedpayload
	type JWSTransaction struct {
		jwt.RegisteredClaims

		AppAccountToken string `json:"appAccountToken"`
		TransactionId   string `json:"transactionId"`
		BundleId        string `json:"bundleId"`
		Type            string `json:"type"`
		Environment     string `json:"environment"`
		Price           int    `json:"price"`
		Currency        string `json:"currency"`
	}

	var payload JWSTransaction
	if err := jws.VerifyAndBind(&payload); err != nil {
		t.Errorf("TestJWS_VerifyAndBind failed. err:%v", err)
		return
	}

	t.Logf("TestJWS_VerifyAndBind got:%#v", payload)

	py, _ := json.Marshal(payload)
	t.Logf("f:%s", py)
}
