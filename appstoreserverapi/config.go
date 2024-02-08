package appstoreserverapi

import (
	"time"
)

// Config config for generate token for api
// 文档：https://developer.apple.com/documentation/appstoreserverapi/generating_tokens_for_api_requests
type Config struct {
	// Your app’s bundle ID (Ex: “com.example.testbundleid2021”)
	BundleID string
	// Your issuer ID from the Keys page in App Store Connect (Ex: "57246542-96fe-1a63-e053-0824d011072a")
	Issuer string
	// Your private key ID from App Store Connect (Ex: 2X9R4HXF34)
	KeyID string
	// Your private key content
	PrivateKey []byte

	// http request timeout
	Timeout time.Duration
}
