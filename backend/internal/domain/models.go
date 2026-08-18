package domain

import "time"

type BotStatus string

const (
	BotUnbound      BotStatus = "unbound"
	BotAwaitingScan BotStatus = "awaiting_scan"
	BotConnecting   BotStatus = "connecting"
	BotOnline       BotStatus = "online"
	BotReconnecting BotStatus = "reconnecting"
	BotOffline      BotStatus = "offline"
	BotError        BotStatus = "error"
)

type BotInstance struct {
	ID              string     `json:"id"`
	PublicAccountID *string    `json:"publicAccountId"`
	DisplayName     *string    `json:"displayName"`
	Status          BotStatus  `json:"status"`
	StatusReason    *string    `json:"statusReason"`
	QRCode          *string    `json:"qrCode"`
	QRExpiresAt     *time.Time `json:"qrExpiresAt"`
	LastSeenAt      *time.Time `json:"lastSeenAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

type BindingStatus string

const (
	BindingPending   BindingStatus = "pending"
	BindingScanned   BindingStatus = "scanned"
	BindingConfirmed BindingStatus = "confirmed"
	BindingExpired   BindingStatus = "expired"
	BindingCancelled BindingStatus = "cancelled"
	BindingFailed    BindingStatus = "failed"
)

type BindingSession struct {
	ID            string        `json:"id"`
	Status        BindingStatus `json:"status"`
	QRCode        *string       `json:"qrCode"`
	CreatedAt     time.Time     `json:"createdAt"`
	ExpiresAt     time.Time     `json:"expiresAt"`
	CompletedAt   *time.Time    `json:"completedAt"`
	FailureReason *string       `json:"failureReason"`
}

type Group struct {
	ID              string    `json:"id"`
	ProviderGroupID string    `json:"providerGroupId"`
	Name            string    `json:"name"`
	MemberCount     *int      `json:"memberCount"`
	Enabled         bool      `json:"enabled"`
	Available       bool      `json:"available"`
	DiscoveredAt    time.Time `json:"discoveredAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type MessageDirection string

const (
	DirectionInbound  MessageDirection = "inbound"
	DirectionOutbound MessageDirection = "outbound"
)

type MessageType string

const (
	MessageText        MessageType = "text"
	MessageUnsupported MessageType = "unsupported"
	MessageSystem      MessageType = "system"
)

type ProcessingStatus string

const (
	ProcessingReceived ProcessingStatus = "received"
	ProcessingIgnored  ProcessingStatus = "ignored"
	ProcessingSent     ProcessingStatus = "sent"
	ProcessingFailed   ProcessingStatus = "failed"
)

type GroupMessage struct {
	ID                string           `json:"id"`
	ProviderMessageID string           `json:"providerMessageId"`
	GroupID           string           `json:"groupId"`
	SenderID          string           `json:"senderId"`
	SenderName        string           `json:"senderName"`
	Direction         MessageDirection `json:"direction"`
	MessageType       MessageType      `json:"messageType"`
	Content           string           `json:"content"`
	OccurredAt        time.Time        `json:"occurredAt"`
	ReceivedAt        time.Time        `json:"receivedAt"`
	ProcessingStatus  ProcessingStatus `json:"processingStatus"`
	SelfMessage       bool             `json:"selfMessage"`
}

type OutboundStatus string

const (
	OutboundPending   OutboundStatus = "pending"
	OutboundSending   OutboundStatus = "sending"
	OutboundSucceeded OutboundStatus = "succeeded"
	OutboundFailed    OutboundStatus = "failed"
)

type OutboundMessageRequest struct {
	ID                string         `json:"id"`
	IdempotencyKey    string         `json:"idempotencyKey"`
	GroupID           string         `json:"groupId"`
	Content           string         `json:"content"`
	RequestedBy       string         `json:"requestedBy"`
	Status            OutboundStatus `json:"status"`
	ProviderMessageID *string        `json:"providerMessageId"`
	FailureReason     *string        `json:"failureReason"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
}

type AuditResult string

const (
	AuditSuccess AuditResult = "success"
	AuditFailure AuditResult = "failure"
)

type AuditEvent struct {
	ID         string      `json:"id"`
	Actor      string      `json:"actor"`
	Action     string      `json:"action"`
	TargetType string      `json:"targetType"`
	TargetID   string      `json:"targetId"`
	Result     AuditResult `json:"result"`
	Detail     *string     `json:"detail"`
	CreatedAt  time.Time   `json:"createdAt"`
}
