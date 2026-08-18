package bot

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type MockGateway struct {
	mu      sync.RWMutex
	events  chan Event
	online  bool
	started bool
	groups  []Group
	counter atomic.Uint64
}

func NewMockGateway() *MockGateway {
	members1, members2 := 12, 6
	return &MockGateway{
		events: make(chan Event, 128),
		groups: []Group{
			{ProviderID: "mock-room-engineering", Name: "研发讨论群", MemberCount: &members1},
			{ProviderID: "mock-room-operations", Name: "运营协作群", MemberCount: &members2},
		},
	}
}
func (m *MockGateway) Driver() string       { return "mock" }
func (m *MockGateway) Events() <-chan Event { return m.events }
func (m *MockGateway) Start(context.Context) error {
	m.mu.Lock()
	m.started = true
	m.mu.Unlock()
	return nil
}
func (m *MockGateway) Stop(context.Context) error {
	m.mu.Lock()
	m.online = false
	m.started = false
	m.mu.Unlock()
	return nil
}
func (m *MockGateway) BeginBinding(context.Context) error {
	m.mu.RLock()
	started := m.started
	m.mu.RUnlock()
	if !started {
		return errors.New("mock gateway is not started")
	}
	expires := time.Now().Add(2 * time.Minute)
	m.events <- Event{Type: EventScan, QRCode: fmt.Sprintf("mock-wechaty-binding-%d", time.Now().UnixNano()), QRExpiresAt: expires}
	return nil
}
func (m *MockGateway) CancelBinding(context.Context) error { return nil }
func (m *MockGateway) ListGroups(context.Context) ([]Group, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.online {
		return nil, errors.New("bot is offline")
	}
	out := append([]Group(nil), m.groups...)
	return out, nil
}
func (m *MockGateway) SendText(_ context.Context, groupID, text string) (string, error) {
	m.mu.RLock()
	online := m.online
	m.mu.RUnlock()
	if !online {
		return "", errors.New("bot is offline")
	}
	if text == "" {
		return "", errors.New("message is empty")
	}
	found := false
	for _, g := range m.groups {
		if g.ProviderID == groupID {
			found = true
			break
		}
	}
	if !found {
		return "", errors.New("group is unavailable")
	}
	return fmt.Sprintf("mock-out-%d", m.counter.Add(1)), nil
}
func (m *MockGateway) SimulateLogin(context.Context) error {
	m.mu.Lock()
	m.online = true
	m.mu.Unlock()
	m.events <- Event{Type: EventLogin, AccountID: "mock-account", DisplayName: "Mock WeChat Bot"}
	m.events <- Event{Type: EventGroups, Groups: append([]Group(nil), m.groups...)}
	return nil
}
func (m *MockGateway) SimulateLogout(_ context.Context, reason string) error {
	m.mu.Lock()
	m.online = false
	m.mu.Unlock()
	m.events <- Event{Type: EventLogout, Reason: reason}
	return nil
}
func (m *MockGateway) InjectMessage(_ context.Context, msg Message) error {
	m.mu.RLock()
	online := m.online
	m.mu.RUnlock()
	if !online {
		return errors.New("bot is offline")
	}
	if msg.ProviderMessageID == "" {
		msg.ProviderMessageID = fmt.Sprintf("mock-in-%d", m.counter.Add(1))
	}
	if msg.OccurredAt.IsZero() {
		msg.OccurredAt = time.Now()
	}
	if msg.MessageType == "" {
		msg.MessageType = "text"
	}
	m.events <- Event{Type: EventMessage, Message: &msg}
	return nil
}
