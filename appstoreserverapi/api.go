package appstoreserverapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

// LookupOrder Get a customer’s in-app purchases from a receipt using the order ID.
// api: https://developer.apple.com/documentation/appstoreserverapi/look_up_order_id
func (s *Service) LookupOrder(ctx context.Context, customerOrderID string) ([]Transaction, error) {
	_, body, err := s.get(ctx, "/inApps/v1/lookup/"+customerOrderID, nil)
	if err != nil {
		return nil, err
	}

	type LookupOrderResp struct {
		// 0 contains an array of one or more signed transactions for the in-app purchase based on the order ID.
		// 1 1 doesn't contain a signed transactions array.
		Status             int              `json:"status"`
		SignedTransactions []JWSTransaction `json:"signedTransactions"`
	}

	var res LookupOrderResp
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if 1 == res.Status {
		return nil, nil
	}

	return JWSTransactions(res.SignedTransactions).GetTransactions()
}

// GetTransactionInfo Get information about a single transaction for your app
// https://developer.apple.com/documentation/appstoreserverapi/get_transaction_info
func (s *Service) GetTransactionInfo(ctx context.Context, transactionID string) (*Transaction, error) {
	_, body, err := s.get(ctx, "/inApps/v1/transactions/"+transactionID, nil)
	if err != nil {
		return nil, err
	}

	type Response struct {
		SignedTransactionInfo JWSTransaction `json:"signedTransactionInfo"`
	}

	var res Response
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return res.SignedTransactionInfo.GetTransaction()
}

type GetTransactionHistoryReq struct {
	TransactionID string
	Query         *GetTransactionHistoryReqQuery
}

type GetTransactionHistoryReqQuery struct {
	// A token you provide to get the next set of up to 20 transactions. All responses include a revision token.
	// Use the revision token from the previous HistoryResponse
	Revision string `form:"revision"`
	// An optional start date of the timespan for the transaction history records you’re requesting.
	// The startDate needs to precede the endDate if you specify both dates.
	// The results include a transaction if its purchaseDate is equal to or greater than the startDate
	// The start date of a timespan, expressed in UNIX time, in milliseconds
	StartDate int64 `form:"startDate"`
	// An optional end date of the timespan for the transaction history records you’re requesting.
	// Choose an endDate that’s later than the startDate if you specify both dates.
	// Using an endDate in the future is valid.
	// The results include a transaction if its purchaseDate is earlier than the endDate.
	EndDate int64 `form:"endDate"`

	// An optional filter that indicates the product identifier to include in the transaction history.
	// Your query may specify more than one productID
	ProductID []string `form:"productId"`

	// An optional filter that indicates the product type to include in the transaction history.
	// Your query may specify more than one productType
	ProductType []ProductType `form:"productType"`

	// An optional sort order for the transaction history records.
	// The response sorts the transaction records by their recently modified date.
	// The default value is ASCENDING, so you receive the oldest records first
	Sort string `form:"sort"`

	// An optional filter that indicates the subscription group identifier to include in the transaction history.
	// Your query may specify more than one subscriptionGroupIdentifier
	SubscriptionGroupIdentifier []string `form:"subscriptionGroupIdentifier"`

	// An optional filter that limits the transaction history by the in-app ownership type
	InAppOwnershipType InAppOwnershipType `form:"inAppOwnershipType"`

	// An optional Boolean value that indicates whether the response includes only revoked transactions
	// when the value is true, or contains only nonrevoked transactions when the value is false.
	// By default, the request doesn't include this parameter
	Revoked *bool `form:"revoked"`
}

func (r *GetTransactionHistoryReqQuery) Values() url.Values {
	if r == nil {
		return nil
	}

	query := url.Values{}
	if r.Revision != "" {
		query.Add("revision", r.Revision)
	}

	if r.StartDate > 0 {
		query.Add("startDate", strconv.FormatInt(r.StartDate, 10))
	}

	if r.EndDate > 0 {
		query.Add("endDate", strconv.FormatInt(r.EndDate, 10))
	}

	for _, v := range r.ProductID {
		query.Add("productId", v)
	}

	for _, v := range r.ProductType {
		query.Add("productType", string(v))
	}

	if r.Sort != "" {
		query.Add("sort", r.Sort)
	}

	for _, v := range r.SubscriptionGroupIdentifier {
		query.Add("subscriptionGroupIdentifier", v)
	}

	if r.InAppOwnershipType != "" {
		query.Add("inAppOwnershipType", string(r.InAppOwnershipType))
	}

	if r.Revoked != nil {
		query.Add("revoked", strconv.FormatBool(*r.Revoked))
	}

	return query
}

// GetTransactionHistoryResp A response that contains the customer’s transaction history for an app
type GetTransactionHistoryResp struct {
	AppAppleID         int64            `json:"appAppleId"`
	BundleID           string           `json:"bundleId"`
	Environment        Environment      `json:"environment"`
	HasMore            bool             `json:"hasMore"`
	Revision           string           `json:"revision"`
	SignedTransactions []JWSTransaction `json:"signedTransactions"`

	// ===== 用于调用 next ====
	service *Service
	req     *GetTransactionHistoryReq
}

// GetTransactions 获取当前返回数据中的交易信息
func (resp *GetTransactionHistoryResp) GetTransactions() ([]Transaction, error) {
	return JWSTransactions(resp.SignedTransactions).GetTransactions()
}

// Next GetTransactionHistory 的下一页
func (resp *GetTransactionHistoryResp) Next(ctx context.Context) (*GetTransactionHistoryResp, error) {
	if !resp.HasMore {
		return nil, nil
	}

	req := resp.req
	if req == nil {
		return nil, errors.New("appstore.appstoreserverapi.GetTransactionHistoryResp: invalid req")
	}

	query := req.Query
	if query == nil {
		query = &GetTransactionHistoryReqQuery{Revision: resp.Revision}
	} else {
		query.Revision = resp.Revision
	}

	req.Query = query
	return resp.service.GetTransactionHistory(ctx, req)
}

// GetTransactionHistory Get a customer’s in-app purchase transaction history for your app
// https://developer.apple.com/documentation/appstoreserverapi/get_transaction_history
func (s *Service) GetTransactionHistory(ctx context.Context, req *GetTransactionHistoryReq) (*GetTransactionHistoryResp, error) {
	_, body, err := s.get(ctx, "/inApps/v1/history/"+req.TransactionID, req.Query.Values())
	if err != nil {
		return nil, err
	}

	var out GetTransactionHistoryResp
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}

	// 调用 next 使用
	out.service = s
	out.req = req

	return &out, nil
}

type GetAllSubscriptionStatusesResp struct {
	Environment Environment                       `json:"environment"`
	AppAppleID  int64                             `json:"appAppleId"`
	BundleID    string                            `json:"bundleId"`
	Data        []SubscriptionGroupIdentifierItem `json:"data"`
}

// GetAllSubscriptionStatuses Get the statuses for all of a customer’s auto-renewable subscriptions in your app.
//   - status: An optional filter that indicates the status of subscriptions to include in the response.
//     Your query may specify more than one status query parameter
//
// https://developer.apple.com/documentation/appstoreserverapi/get_all_subscription_statuses
func (s *Service) GetAllSubscriptionStatuses(ctx context.Context, transactionID string, status []AutoRenewableSubscriptionStatus) (*GetAllSubscriptionStatusesResp, error) {
	query := url.Values{}
	for _, v := range status {
		query.Add("status", strconv.FormatInt(int64(v), 10))
	}

	_, body, err := s.get(ctx, "/inApps/v1/subscriptions/"+transactionID, query)
	if err != nil {
		return nil, err
	}

	var out GetAllSubscriptionStatusesResp
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

type GetRefundHistoryResp struct {
	HasMore            bool             `json:"hasMore"`
	Revision           string           `json:"revision"`
	SignedTransactions []JWSTransaction `json:"signedTransactions"`

	// ==== request ====
	service       *Service
	transactionID string
}

// GetTransactions 获取当前返回数据中的交易信息
func (resp *GetRefundHistoryResp) GetTransactions() ([]Transaction, error) {
	return JWSTransactions(resp.SignedTransactions).GetTransactions()
}

// Next GetTransactionHistoryResp 的下一页
func (resp *GetRefundHistoryResp) Next(ctx context.Context) (*GetRefundHistoryResp, error) {
	if !resp.HasMore {
		return nil, nil
	}

	return resp.service.GetRefundHistory(ctx, resp.transactionID, resp.Revision)
}

// GetRefundHistory https://developer.apple.com/documentation/appstoreserverapi/get_refund_history
func (s *Service) GetRefundHistory(ctx context.Context, transactionID string, revision string) (*GetRefundHistoryResp, error) {
	query := url.Values{}
	if revision != "" {
		query.Add("revision", revision)
	}

	_, body, err := s.get(ctx, "/inApps/v2/refund/lookup/"+transactionID, query)
	if err != nil {
		return nil, err
	}

	var out GetRefundHistoryResp
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}

	out.transactionID = transactionID
	out.service = s

	return &out, nil
}
