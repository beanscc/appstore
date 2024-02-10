package appstoreserverapi

import (
	"github.com/beanscc/appstore/jws"
	"github.com/golang-jwt/jwt/v5"
)

// Environment https://developer.apple.com/documentation/appstoreserverapi/environment
type Environment string

const (
	EnvironmentSandbox    Environment = "Sandbox"
	EnvironmentProduction Environment = "Production"
)

// InAppOwnershipType https://developer.apple.com/documentation/appstoreserverapi/inappownershiptype
type InAppOwnershipType string

const (
	InAppOwnershipTypeFamilyShared InAppOwnershipType = "FAMILY_SHARED"
	InAppOwnershipTypePurchased    InAppOwnershipType = "PURCHASED"
)

// OfferDiscountType https://developer.apple.com/documentation/appstoreserverapi/offerdiscounttype
type OfferDiscountType string

const (
	OfferDiscountTypeFreeTrial  OfferDiscountType = "FREE_TRIAL"
	OfferDiscountTypePayAsYouGo OfferDiscountType = "PAY_AS_YOU_GO"
	OfferDiscountTypePayUpFront OfferDiscountType = "PAY_UP_FRONT"
)

// OfferType https://developer.apple.com/documentation/appstoreserverapi/offertype
type OfferType int

const (
	OfferTypeIntroductory OfferType = 1
	OfferTypePromotional  OfferType = 2
	OfferTypeSubscription OfferType = 3
)

// TransactionReason https://developer.apple.com/documentation/appstoreserverapi/transactionreason
type TransactionReason string

const (
	TransactionReasonPurchase TransactionReason = "PURCHASE"
	TransactionReasonRenewal  TransactionReason = "RENEWAL"
)

// TransactionType https://developer.apple.com/documentation/appstoreserverapi/type
type TransactionType string

const (
	TransactionTypeAutoRenewableSubscription TransactionType = "Auto-Renewable Subscription"
	TransactionTypeNonRenewingSubscription   TransactionType = "Non-Renewing Subscription"
	TransactionTypeConsumable                TransactionType = "Consumable"
	TransactionTypeNonConsumable             TransactionType = "Non-Consumable"
)

type ProductType string

const (
	ProductTypeAutoRenewable ProductType = "AUTO_RENEWABLE"
	ProductTypeNonRenewable  ProductType = "NON_RENEWABLE"
	ProductTypeConsumable    ProductType = "CONSUMABLE"
	ProductTypeNonConsumable ProductType = "NON_CONSUMABLE"
)

type Sort string

const (
	SortAsc  Sort = "ASCENDING"
	SortDesc Sort = "DESCENDING"
)

// Transaction A decoded payload that contains transaction information
// https://developer.apple.com/documentation/appstoreserverapi/jwstransactiondecodedpayload
type Transaction struct {
	// A UUID you create at the time of purchase that associates the transaction with a customer on your own service.
	// If your app doesn’t provide an appAccountToken, this string is empty.
	AppAccountToken string `json:"appAccountToken,omitempty"`
	// The bundle identifier of the app.
	BundleID string `json:"bundleId,omitempty"`
	// The server environment, either sandbox or production
	Environment Environment `json:"environment,omitempty"`

	// A string that describes whether the transaction was purchased by the customer,
	// or is available to them through Family Sharing
	InAppOwnershipType InAppOwnershipType `json:"inAppOwnershipType,omitempty"`

	// The UNIX time, in milliseconds, that represents the purchase date of the original transaction identifier
	OriginalPurchaseDate int64 `json:"originalPurchaseDate,omitempty"`
	// The transaction identifier of the original purchase
	OriginalTransactionID string `json:"originalTransactionId,omitempty"`

	// The unique identifier of the product
	ProductID string `json:"productId,omitempty"`
	// The three-letter ISO 4217 currency code associated with the price parameter.
	// This value is present only if price is present
	Currency string `json:"currency,omitempty"`
	// An integer value that represents the price multiplied by 1000 of the in-app purchase or subscription offer
	// you configured in App Store Connect and that the system records at the time of the purchase.
	// For more information, see price. The currency parameter indicates the currency of this price
	Price int `json:"price,omitempty"`
	// The number of consumable products the customer purchased.
	Quantity int `json:"quantity,omitempty"`
	// The unique identifier of the transaction
	TransactionID string `json:"transactionId,omitempty"`
	// The type of the in-app purchase
	Type TransactionType `json:"type,omitempty"`
	// The reason for the purchase transaction, which indicates whether it’s a customer’s purchase or a renewal
	// for an auto-renewable subscription that the system initiates
	TransactionReason TransactionReason `json:"transactionReason,omitempty"`
	// The UNIX time, in milliseconds, that the App Store charged the customer’s account for a purchase,
	// restored product, subscription, or subscription renewal after a lapse
	PurchaseDate int64 `json:"purchaseDate,omitempty"`

	// The UNIX time, in milliseconds, that the App Store refunded the transaction or revoked it from Family Sharing
	RevocationDate int64 `json:"revocationDate,omitempty"`
	// The reason that the App Store refunded the transaction or revoked it from Family Sharing.
	RevocationReason *int `json:"revocationReason,omitempty"`

	// The UNIX time, in milliseconds, that the App Store signed the JSON Web Signature (JWS) data
	SignedDate int64 `json:"signedDate,omitempty"`
	// The three-letter code that represents the country or region associated with the App Store storefront for the purchase
	Storefront string `json:"storefront,omitempty"`
	// An Apple-defined value that uniquely identifies the App Store storefront associated with the purchase
	StorefrontID string `json:"storefrontId,omitempty"`

	// ===== only to auto-renewable subscriptions ====

	// The identifier of the subscription group to which the subscription belongs.
	SubscriptionGroupIdentifier string `json:"subscriptionGroupIdentifier,omitempty"`
	// A Boolean value that indicates whether the customer upgraded to another subscription
	IsUpgraded bool `json:"isUpgraded,omitempty"`
	// The UNIX time, in milliseconds, that the subscription expires or renews.
	ExpiresDate int64 `json:"expiresDate,omitempty"`
	// The payment mode you configure for the subscription offer
	OfferDiscountType OfferDiscountType `json:"offerDiscountType,omitempty"`
	// The identifier that contains the offer code or the promotional offer identifier
	// The offerIdentifier applies only when the offerType has a value of 2 or 3.
	// The offerIdentifier provides details about the subscription offer in effect for the transaction.
	// Its value is either the offer code or the promotional offer
	OfferIdentifier string `json:"offerIdentifier,omitempty"`
	// A value that represents the promotional offer type.
	OfferType OfferType `json:"offerType,omitempty"`
	// The unique identifier of subscription purchase events across devices, including subscription renewals.
	WebOrderLineItemID string `json:"webOrderLineItemId,omitempty"`
}

type JWSTransaction string

func (s JWSTransaction) GetTransaction() (*Transaction, error) {
	val, err := jws.Parse(string(s))
	if err != nil {
		return nil, err
	}

	type Payload struct {
		jwt.RegisteredClaims
		Transaction
	}
	var out Payload
	if err := val.VerifyAndBind(&out); err != nil {
		return nil, err
	}

	return &out.Transaction, nil
}

type JWSTransactions []JWSTransaction

func (ts JWSTransactions) GetTransactions() ([]Transaction, error) {
	transactions := make([]Transaction, 0, len(ts))
	for _, v := range ts {
		transaction, err := v.GetTransaction()
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, *transaction)
	}
	return transactions, nil
}

// AutoRenewableSubscriptionStatus 自动续订状态
// https://developer.apple.com/documentation/appstoreserverapi/status
type AutoRenewableSubscriptionStatus int32

const (
	// The auto-renewable subscription is active
	AutoRenewableSubscriptionStatusActive AutoRenewableSubscriptionStatus = 1
	// The auto-renewable subscription is expired
	AutoRenewableSubscriptionStatusExpired AutoRenewableSubscriptionStatus = 2
	// The auto-renewable subscription is in a billing retry period
	AutoRenewableSubscriptionStatusInBillingRetryPeriod AutoRenewableSubscriptionStatus = 3
	// The auto-renewable subscription is in a Billing Grace Period
	AutoRenewableSubscriptionStatusInBillingGracePeriod AutoRenewableSubscriptionStatus = 4
	// The auto-renewable subscription is revoked.
	// The App Store refunded the transaction or revoked it from Family Sharing.
	AutoRenewableSubscriptionStatusRevoked AutoRenewableSubscriptionStatus = 5
)

// AutoRenewStatus https://developer.apple.com/documentation/appstoreserverapi/autorenewstatus
type AutoRenewStatus int32

const (
	AutoRenewStatusOff AutoRenewStatus = 0
	AutoRenewStatusOn  AutoRenewStatus = 1
)

// ExpirationIntent https://developer.apple.com/documentation/appstoreserverapi/expirationintent
type ExpirationIntent int32

const (
	ExpirationIntentCustomerCanceled                  ExpirationIntent = 1
	ExpirationIntentBillingError                      ExpirationIntent = 2
	ExpirationIntentCustomerDoNotConsentPriceIncrease ExpirationIntent = 3
	ExpirationIntentProductNotAvailable               ExpirationIntent = 4
	ExpirationIntentOtherReason                       ExpirationIntent = 5
)

// RenewalInfo subscription renewal information for an auto-renewable subscription
// https://developer.apple.com/documentation/appstoreserverapi/jwsrenewalinfodecodedpayload
type RenewalInfo struct {
	AutoRenewProductID          string           `json:"autoRenewProductId"`
	AutoRenewStatus             AutoRenewStatus  `json:"autoRenewStatus"`
	Environment                 Environment      `json:"environment"`
	ExpirationIntent            ExpirationIntent `json:"expirationIntent"`
	GracePeriodExpiresDate      int64            `json:"gracePeriodExpiresDate"`
	IsInBillingRetryPeriod      bool             `json:"isInBillingRetryPeriod"`
	OfferIdentifier             string           `json:"offerIdentifier"`
	OfferType                   OfferType        `json:"offerType"`
	OriginalTransactionID       string           `json:"originalTransactionId"`
	PriceIncreaseStatus         *int32           `json:"priceIncreaseStatus"`
	ProductID                   string           `json:"productId"`
	RecentSubscriptionStartDate int64            `json:"recentSubscriptionStartDate"`
	RenewalDate                 int64            `json:"renewalDate"`
	SignedDate                  int64            `json:"signedDate"`
}

type JWSRenewalInfo string

func (s JWSRenewalInfo) GetRenewInfo() (*RenewalInfo, error) {
	val, err := jws.Parse(string(s))
	if err != nil {
		return nil, err
	}
	type Payload struct {
		jwt.RegisteredClaims
		RenewalInfo
	}
	var out Payload
	if err := val.VerifyAndBind(&out); err != nil {
		return nil, err
	}

	return &out.RenewalInfo, nil
}

type SubscriptionGroupIdentifierItem struct {
	SubscriptionGroupIdentifier string                         `json:"subscriptionGroupIdentifier"`
	LastTransactions            []SubscriptionLastTransactions `json:"lastTransactions"`
}

type SubscriptionLastTransactions struct {
	OriginalTransactionID string                          `json:"originalTransactionId"`
	Status                AutoRenewableSubscriptionStatus `json:"status"`
	SignedRenewalInfo     JWSRenewalInfo                  `json:"signedRenewalInfo"`
	SignedTransactionInfo JWSTransaction                  `json:"signedTransactionInfo"`
}

// NotificationType The type that describes the in-app purchase event for which the App Store sends the version 2 notification
// https://developer.apple.com/documentation/appstoreservernotifications/notificationtype
type NotificationType string

const (
	NotificationTypeConsumptionRequest   NotificationType = "CONSUMPTION_REQUEST"
	NotificationTypeDidChangeRenewPref   NotificationType = "DID_CHANGE_RENEWAL_PREF"
	NotificationTypeDidChangeRenewStatus NotificationType = "DID_CHANGE_RENEWAL_STATUS"
	NotificationTypeDidFailToRenew       NotificationType = "DID_FAIL_TO_RENEW"
	NotificationTypeDidRenew             NotificationType = "DID_RENEW"
	NotificationTypeExpired              NotificationType = "EXPIRED"
	NotificationTypeGracePeriodExpired   NotificationType = "GRACE_PERIOD_EXPIRED"
	NotificationTypeOfferRedeemed        NotificationType = "OFFER_REDEEMED"
	NotificationTypePriceIncrease        NotificationType = "PRICE_INCREASE"
	NotificationTypeRefund               NotificationType = "REFUND"
	NotificationTypeRefundDeclined       NotificationType = "REFUND_DECLINED"
	NotificationTypeRefundReversed       NotificationType = "REFUND_REVERSED"
	NotificationTypeRenewalExtended      NotificationType = "RENEWAL_EXTENDED"
	NotificationTypeRenewalExtension     NotificationType = "RENEWAL_EXTENSION"
	NotificationTypeRevoke               NotificationType = "REVOKE"
	NotificationTypeSubscribed           NotificationType = "SUBSCRIBED"
	NotificationTypeTest                 NotificationType = "TEST"
)

// NotificationSubtype A string that provides details about select notification types in version 2
// https://developer.apple.com/documentation/appstoreservernotifications/subtype
type NotificationSubtype string

const (
	NotificationSubtypeAccepted          NotificationSubtype = "ACCEPTED"
	NotificationSubtypeAutoRenewDisabled NotificationSubtype = "AUTO_RENEW_DISABLED"
	NotificationSubtypeAutoRenewEnabled  NotificationSubtype = "AUTO_RENEW_ENABLED"
	NotificationSubtypeBillingRecovery   NotificationSubtype = "BILLING_RECOVERY"
	NotificationSubtypeBillingRetry      NotificationSubtype = "BILLING_RETRY"
	NotificationSubtypeDowngrade         NotificationSubtype = "DOWNGRADE"
	NotificationSubtypeFailure           NotificationSubtype = "FAILURE"
	NotificationSubtypeGracePeriod       NotificationSubtype = "GRACE_PERIOD"
	NotificationSubtypeInitialBuy        NotificationSubtype = "INITIAL_BUY"
	NotificationSubtypePending           NotificationSubtype = "PENDING"
	NotificationSubtypePriceIncrease     NotificationSubtype = "PRICE_INCREASE"
	NotificationSubtypeProductNotForSale NotificationSubtype = "PRODUCT_NOT_FOR_SALE"
	NotificationSubtypeResubscribe       NotificationSubtype = "RESUBSCRIBE"
	NotificationSubtypeSummary           NotificationSubtype = "SUMMARY"
	NotificationSubtypeUpgrade           NotificationSubtype = "UPGRADE"
	NotificationSubtypeVoluntary         NotificationSubtype = "VOLUNTARY"
)

type NotificationHistoryItem struct {
	SendAttempts  []NotificationSendAttemptItem `json:"sendAttempts"`
	SignedPayload JWSNotification               `json:"signedPayload"`
}

type JWSNotification string

func (s JWSNotification) GetNotification() (*Notification, error) {
	val, err := jws.Parse(string(s))
	if err != nil {
		return nil, err
	}
	type Payload struct {
		jwt.RegisteredClaims
		Notification
	}
	var out Payload
	if err := val.VerifyAndBind(&out); err != nil {
		return nil, err
	}

	return &out.Notification, nil
}

// Notification the version 2 notification data
// https://developer.apple.com/documentation/appstoreservernotifications/responsebodyv2decodedpayload
type Notification struct {
	NotificationType NotificationType    `json:"notificationType"`
	Subtype          NotificationSubtype `json:"subtype"`
	Data             NotificationData    `json:"data"`
	Summary          NotificationSummary `json:"summary"`
	Version          string              `json:"version"`
	SignedDate       int64               `json:"signedDate"`
	NotificationUUID string              `json:"notificationUUID"`
}

type NotificationData struct {
	AppAppleID            int64                           `json:"appAppleId"`
	BundleID              string                          `json:"bundleId"`
	BundleVersion         string                          `json:"bundleVersion"`
	Environment           Environment                     `json:"environment"`
	SignedRenewalInfo     JWSRenewalInfo                  `json:"signedRenewalInfo"`
	SignedTransactionInfo JWSTransaction                  `json:"signedTransactionInfo"`
	Status                AutoRenewableSubscriptionStatus `json:"status"`
}

// NotificationSummary https://developer.apple.com/documentation/appstoreservernotifications/summary
type NotificationSummary struct {
	RequestIdentifier      string      `json:"requestIdentifier"`
	Environment            Environment `json:"environment"`
	AppAppleID             int64       `json:"appAppleId"`
	BundleID               string      `json:"bundleId"`
	ProductID              string      `json:"productId"`
	StorefrontCountryCodes []string    `json:"storefrontCountryCodes"`
	FailedCount            int64       `json:"failedCount"`
	SucceededCount         int64       `json:"succeededCount"`
}

type NotificationSendAttemptItem struct {
	AttemptDate       int64                         `json:"attemptDate"`
	SendAttemptResult NotificationSendAttemptResult `json:"sendAttemptResult"`
}

// NotificationSendAttemptResult The success or error information the App Store server records when it attempts to send an App Store server notification to your server
// https://developer.apple.com/documentation/appstoreserverapi/sendattemptresult
type NotificationSendAttemptResult string

const (
	NotificationSendAttemptResultSuccess                      NotificationSendAttemptResult = "SUCCESS"
	NotificationSendAttemptResultCircularRedirect             NotificationSendAttemptResult = "CIRCULAR_REDIRECT"
	NotificationSendAttemptResultInvalidResponse              NotificationSendAttemptResult = "INVALID_RESPONSE"
	NotificationSendAttemptResultNoResponse                   NotificationSendAttemptResult = "NO_RESPONSE"
	NotificationSendAttemptResultOther                        NotificationSendAttemptResult = "OTHER"
	NotificationSendAttemptResultPrematureClose               NotificationSendAttemptResult = "PREMATURE_CLOSE"
	NotificationSendAttemptResultSocketIssue                  NotificationSendAttemptResult = "SOCKET_ISSUE"
	NotificationSendAttemptResultTimeout                      NotificationSendAttemptResult = "TIMED_OUT"
	NotificationSendAttemptResultTLSIssue                     NotificationSendAttemptResult = "TLS_ISSUE"
	NotificationSendAttemptResultUnsuccessfulHttpResponseCode NotificationSendAttemptResult = "UNSUCCESSFUL_HTTP_RESPONSE_CODE"
	NotificationSendAttemptResultUnsupportedCharset           NotificationSendAttemptResult = "UNSUPPORTED_CHARSET"
)
