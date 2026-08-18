package bot

import (
	"context"
	"time"
)

type EventType string

const (
	EventScan    EventType = "scan"
	EventLogin   EventType = "login"
	EventLogout  EventType = "logout"
	EventError   EventType = "error"
	EventMessage EventType = "message"
	EventGroups  EventType = "groups"
)

type Group struct {
	ProviderID  string
	Name        string
	MemberCount *int
}

type Message struct {
	ProviderMessageID string
	ProviderGroupID   string
	GroupName         string
	SenderID          string
	SenderName        string
	Text              string
	MessageType       string
	OccurredAt        time.Time
	Self              bool
}

type Event struct {
	Type        EventType
	QRCode      string
	QRExpiresAt time.Time
	AccountID   string
	DisplayName string
	Reason      string
	Message     *Message
	Groups      []Group
}

type Gateway interface {
	Driver() string
	Start(context.Context) error
	Stop(context.Context) error
	BeginBinding(context.Context) error
	CancelBinding(context.Context) error
	ListGroups(context.Context) ([]Group, error)
	SendText(context.Context, string, string) (string, error)
	Events() <-chan Event
}

type MockController interface {
	SimulateLogin(context.Context) error
	SimulateLogout(context.Context, string) error
	InjectMessage(context.Context, Message) error
}
