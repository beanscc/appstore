package jws

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// AppleRootCA-G3.cer download from Apple PKI: https://www.apple.com/certificateauthority/
// openssl x509 -in AppleRootCA-G3.cer -inform DER -out AppleRootCA-G3.pem -outform PEM
const appleRootCertificate = `
-----BEGIN CERTIFICATE-----
MIICQzCCAcmgAwIBAgIILcX8iNLFS5UwCgYIKoZIzj0EAwMwZzEbMBkGA1UEAwwS
QXBwbGUgUm9vdCBDQSAtIEczMSYwJAYDVQQLDB1BcHBsZSBDZXJ0aWZpY2F0aW9u
IEF1dGhvcml0eTETMBEGA1UECgwKQXBwbGUgSW5jLjELMAkGA1UEBhMCVVMwHhcN
MTQwNDMwMTgxOTA2WhcNMzkwNDMwMTgxOTA2WjBnMRswGQYDVQQDDBJBcHBsZSBS
b290IENBIC0gRzMxJjAkBgNVBAsMHUFwcGxlIENlcnRpZmljYXRpb24gQXV0aG9y
aXR5MRMwEQYDVQQKDApBcHBsZSBJbmMuMQswCQYDVQQGEwJVUzB2MBAGByqGSM49
AgEGBSuBBAAiA2IABJjpLz1AcqTtkyJygRMc3RCV8cWjTnHcFBbZDuWmBSp3ZHtf
TjjTuxxEtX/1H7YyYl3J6YRbTzBPEVoA/VhYDKX1DyxNB0cTddqXl5dvMVztK517
IDvYuVTZXpmkOlEKMaNCMEAwHQYDVR0OBBYEFLuw3qFYM4iapIqZ3r6966/ayySr
MA8GA1UdEwEB/wQFMAMBAf8wDgYDVR0PAQH/BAQDAgEGMAoGCCqGSM49BAMDA2gA
MGUCMQCD6cHEFl4aXTQY2e3v9GwOAEZLuN+yRhHFD/3meoyhpmvOwgPUnPWTxnS4
at+qIxUCMG1mihDK1A3UT82NQz60imOlM27jbdoXt2QfyFMm+YhidDkLF1vLUagM
6BgD56KyKA==
-----END CERTIFICATE-----
`

// JWS App Store in JSON Web Signature (JWS) format
type JWS struct {
	// raw token
	token string

	Header *Header
}

// Header JWSDecodedHeader
// https://developer.apple.com/documentation/appstoreserverapi/jwsdecodedheader
type Header struct {
	Alg string   `json:"alg"`
	X5C []string `json:"x5c"`
}

// Certificate 从 header x5c 证书链中获取证书
// 依次是：
//   - Certificate(0) Apple leaf certificate
//   - Certificate(1) Apple intermediate certificate
//   - Certificate(2) Apple root certificate
func (h *Header) Certificate(index int) (*x509.Certificate, error) {
	bytes, err := base64.StdEncoding.DecodeString(h.X5C[index])
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(bytes)
}

// Verify 验证 x5c 证书链，并返回用于 jws 签名的公钥
func (h *Header) Verify() (*ecdsa.PublicKey, error) {
	leaf, err := h.Certificate(0)
	if err != nil {
		return nil, err
	}

	intermediate, err := h.Certificate(1)
	if err != nil {
		return nil, err
	}

	root, err := h.Certificate(2)
	if err != nil {
		return nil, err
	}

	// 验证证书链
	if err := h.verify(leaf, intermediate, root); err != nil {
		return nil, err
	}

	return leaf.PublicKey.(*ecdsa.PublicKey), nil
}

func (h *Header) verify(leaf, intermediate, root *x509.Certificate) error {
	// 验证证书链
	opts := x509.VerifyOptions{
		Roots:         x509.NewCertPool(),
		Intermediates: x509.NewCertPool(),
	}
	opts.Roots.AddCert(root)
	opts.Intermediates.AddCert(intermediate)

	_, err := leaf.Verify(opts)
	if err != nil {
		return err
	}

	// // debug 证书
	// for _, ch := range chains {
	// 	for _, c := range ch {
	// 		fmt.Printf("[leaf cert verify] issuer:%s, subject: name:%s, org:%s, org-uni:%s \n",
	// 			c.Issuer.CommonName, c.Subject.CommonName, c.Subject.Organization, c.Subject.OrganizationalUnit)
	// 	}
	// }

	/* output:
	[leaf cert verify] issuer:Apple Worldwide Developer Relations Certification Authority, subject: name:Prod ECC Mac App Store and iTunes Store Receipt Signing, org:[Apple Inc.], org-uni:[Apple Worldwide Developer Relations]
	[leaf cert verify] issuer:Apple Root CA - G3, subject: name:Apple Worldwide Developer Relations Certification Authority, org:[Apple Inc.], org-uni:[G6]
	[leaf cert verify] issuer:Apple Root CA - G3, subject: name:Apple Root CA - G3, org:[Apple Inc.], org-uni:[Apple Certification Authority]
	*/

	// 使用 apple 官方根证书验证回调中的 root 证书
	// Apple PKI 官方提供的根证书
	rootFromApplePKI := x509.NewCertPool()
	if ok := rootFromApplePKI.AppendCertsFromPEM([]byte(appleRootCertificate)); !ok {
		return errors.New("failed to append apple root certificate")
	}

	_, err = root.Verify(x509.VerifyOptions{Roots: rootFromApplePKI})
	if err != nil {
		return err
	}

	// // debug 证书链
	// for _, ch := range chains {
	// 	for _, c := range ch {
	// 		fmt.Printf("[root cert verify] issuer:%s, subject: name:%s, org:%s, org-uni:%s \n",
	// 			c.Issuer.CommonName, c.Subject.CommonName, c.Subject.Organization, c.Subject.OrganizationalUnit,
	// 		)
	// 	}
	// }

	return nil
}

// Parse return JWS by app store jws token
func Parse(token string) (*JWS, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid app store JWS token")
	}

	headerByte, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}

	var h Header
	if err := json.Unmarshal(headerByte, &h); err != nil {
		return nil, err
	}

	return &JWS{
		token:  token,
		Header: &h,
	}, nil
}

// VerifyAndBind 验证 jsw 签名及 x5c 证书链，并绑定 jws payload 到 claims
//
//	type claims struct {
//		jwt.RegisteredClaims
//		CustomClaims
//	}
func (jws *JWS) VerifyAndBind(claims jwt.Claims) error {
	_, err := jwt.ParseWithClaims(jws.token, claims, func(token *jwt.Token) (interface{}, error) {
		return jws.Header.Verify()
	})

	return err
}
