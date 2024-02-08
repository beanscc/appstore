package serverapi

import (
	"context"
	"testing"
	"time"
)

func testConfig() *Config {
	return &Config{
		BundleID: `Your app’s bundle ID (Ex: “com.example.testbundleid2021”)`,
		Issuer:   `Your issuer ID from the Keys page in App Store Connect (Ex: "57246542-96fe-1a63-e053-0824d011072a")`,
		KeyID:    `Your private key ID from App Store Connect (Ex: 2X9R4HXF34)`,
		PrivateKey: []byte(`-----BEGIN PRIVATE KEY-----
YOUR PRIVATE KEY
-----END PRIVATE KEY-----`),
		Timeout: 20 * time.Second,
	}
}

func TestService_LookupOrder(t *testing.T) {
	token := NewToken(testConfig())
	ctx := context.Background()
	service := NewService(token).Debug(false)
	customerOrderID := `MTV70QV5J9`
	got, err := service.LookupOrder(ctx, customerOrderID)
	if err != nil {
		t.Errorf("TestService_LookupOrder failed. err:%v", err)
		return
	}

	for i, v := range got {
		t.Logf("TestService_LookupOrder got idx:%3d, v:%#v", i, v)
	}
}

func TestService_GetTransactionInfo(t *testing.T) {
	token := NewToken(testConfig())
	ctx := context.Background()
	service := NewService(token).Debug(false)
	transactionID := `350001859400409`
	got, err := service.GetTransactionInfo(ctx, transactionID)
	if err != nil {
		t.Errorf("TestService_GetTransactionInfo failed. err:%v", err)
		return
	}

	t.Logf("TestService_GetTransactionInfo: got:%#v", got)
}

func TestService_GetTransactionHistory(t *testing.T) {
	ctx := context.Background()
	token := NewToken(testConfig())
	service := NewService(token).Debug(false)
	req := GetTransactionHistoryReq{
		TransactionID: `350001859400409`,
		Query:         nil,
	}
	got, err := service.GetTransactionHistory(ctx, &req)
	if err != nil {
		t.Errorf("TestService_GetTransactionHistory failed. err:%v", err)
		return
	}

	total := 0
	for n := 0; ; n++ {
		loop := n + 1
		t.Logf("TestService_GetTransactionHistory loop:%d base:%#v", loop, got.TransactionHistoryBase)
		for i, v := range got.Transactions {
			t.Logf("TestService_GetTransactionHistory loop:%d, got idx:%3d, v:%#v", loop, i, v)
			total++
		}

		if got.HasMore {
			got, err = got.Next(context.Background())
		} else {
			break
		}
	}

	t.Logf("total:%d", total)
}

func TestService_GetAllSubscriptionStatuses(t *testing.T) {
	ctx := context.Background()
	token := NewToken(testConfig())
	service := NewService(token).Debug(false)

	transactionID := `350000614215995`
	status := []AutoRenewableSubscriptionStatus{
		// AutoRenewableSubscriptionStatusActive,
	}
	got, err := service.GetAllSubscriptionStatuses(ctx, transactionID, status)
	if err != nil {
		t.Errorf("TestService_GetAllSubscriptionStatuses failed. err:%v", err)
		return
	}

	for i, v := range got.Data {
		t.Logf("TestService_GetAllSubscriptionStatuses got:%3d, v:%#v", i, v.SubscriptionGroupIdentifier)
		for i1, v1 := range v.LastTransactions {
			renewInfo, _ := v1.SignedRenewalInfo.GetRenewInfo()
			transactionInfo, _ := v1.SignedTransactionInfo.GetTransaction()
			t.Logf("TestService_GetAllSubscriptionStatuses.LastTransactions i:%3d, v.originalTransactionId:%s, Status:%d, JWSRenewalInfo:%#v, JWSTransaction:%#v",
				i1, v1.OriginalTransactionID, v1.Status, renewInfo, transactionInfo)
		}
	}
}

func TestService_GetRefundHistory(t *testing.T) {
	ctx := context.Background()
	token := NewToken(testConfig())
	service := NewService(token).Debug(false)

	transactionID := `140002007488219`
	revision := ``
	got, err := service.GetRefundHistory(ctx, transactionID, revision)
	if err != nil {
		t.Errorf("TestService_GetRefundHistory failed. err:%v", err)
		return
	}

	total := 0
	for n := 0; ; n++ {
		loop := n + 1
		t.Logf("TestService_GetRefundHistory loop:%d, got revision:%s", loop, got.Revision)
		for i, v := range got.SignedTransactions {
			transaction, err := v.GetTransaction()
			if err != nil {
				t.Errorf("TestService_GetRefundHistory v.GetTransaction failed. err:%v", err)
				return
			}
			t.Logf("TestService_GetRefundHistory loop:%d, got idx:%3d, v:%#v", loop, i, transaction)
			total++
		}

		if got.HasMore {
			got, err = got.Next(context.Background())
		} else {
			break
		}
	}
	t.Logf("total:%d", total)
}
