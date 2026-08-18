package bot

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/wechaty/go-wechaty/wechaty"
	wp "github.com/wechaty/go-wechaty/wechaty-puppet"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

type WechatyGateway struct {
	token   string
	events  chan Event
	mu      sync.RWMutex
	bot     *wechaty.Wechaty
	started bool
}

func NewWechatyGateway(token string) *WechatyGateway {
	return &WechatyGateway{token: token, events: make(chan Event, 128)}
}
func (w *WechatyGateway) Driver() string       { return "wechaty" }
func (w *WechatyGateway) Events() <-chan Event { return w.events }

// Start prepares the adapter. The real Puppet connection is intentionally
// deferred until BeginBinding so scan events always have an active binding session.
func (w *WechatyGateway) Start(context.Context) error { return nil }

func (w *WechatyGateway) BeginBinding(context.Context) error {
	w.mu.RLock()
	started := w.started
	w.mu.RUnlock()
	if started {
		return nil
	}
	b := wechaty.NewWechaty(wechaty.WithPuppetOption(wp.Option{Token: w.token}))
	b.OnScan(func(_ *wechaty.Context, qr string, status schemas.ScanStatus, _ string) {
		if status == schemas.ScanStatusWaiting || status == schemas.ScanStatusTimeout {
			w.events <- Event{Type: EventScan, QRCode: qr, QRExpiresAt: time.Now().Add(2 * time.Minute)}
		}
	}).OnLogin(func(_ *wechaty.Context, c *user.ContactSelf) {
		w.events <- Event{Type: EventLogin, AccountID: c.ID(), DisplayName: c.Name()}
		if groups, err := w.ListGroups(context.Background()); err == nil {
			w.events <- Event{Type: EventGroups, Groups: groups}
		}
	}).OnLogout(func(_ *wechaty.Context, _ *user.ContactSelf, reason string) {
		w.events <- Event{Type: EventLogout, Reason: reason}
	}).
		OnError(func(_ *wechaty.Context, err error) { w.events <- Event{Type: EventError, Reason: err.Error()} }).
		OnMessage(func(_ *wechaty.Context, m *user.Message) {
			room := m.Room()
			if room == nil {
				return
			}
			talker := m.Talker()
			senderID, senderName := "unknown", "Unknown"
			if talker != nil {
				senderID = talker.ID()
				senderName = talker.Name()
			}
			kind := "unsupported"
			if m.Type() == schemas.MessageTypeText {
				kind = "text"
			}
			w.events <- Event{Type: EventMessage, Message: &Message{ProviderMessageID: m.ID(), ProviderGroupID: room.ID(), GroupName: room.Topic(), SenderID: senderID, SenderName: senderName, Text: m.Text(), MessageType: kind, OccurredAt: m.Date(), Self: m.Self()}}
		})
	w.mu.Lock()
	if w.started {
		w.mu.Unlock()
		return nil
	}
	w.bot = b
	w.started = true
	w.mu.Unlock()
	if err := b.Start(); err != nil {
		w.mu.Lock()
		w.bot = nil
		w.started = false
		w.mu.Unlock()
		return err
	}
	return nil
}
func (w *WechatyGateway) Stop(context.Context) error {
	w.mu.Lock()
	b := w.bot
	w.bot = nil
	w.started = false
	w.mu.Unlock()
	if b != nil && b.Puppet() != nil {
		b.Puppet().Stop()
	}
	return nil
}
func (w *WechatyGateway) CancelBinding(ctx context.Context) error { return w.Stop(ctx) }
func (w *WechatyGateway) ListGroups(context.Context) ([]Group, error) {
	w.mu.RLock()
	b := w.bot
	started := w.started
	w.mu.RUnlock()
	if !started || b == nil {
		return nil, errors.New("wechaty is not started")
	}
	rooms := b.Room().FindAll(nil)
	out := make([]Group, 0, len(rooms))
	for _, room := range rooms {
		members, err := room.MemberAll(nil)
		var count *int
		if err == nil {
			n := len(members)
			count = &n
		}
		out = append(out, Group{ProviderID: room.ID(), Name: room.Topic(), MemberCount: count})
	}
	return out, nil
}
func (w *WechatyGateway) SendText(_ context.Context, groupID, text string) (string, error) {
	w.mu.RLock()
	b := w.bot
	started := w.started
	w.mu.RUnlock()
	if !started || b == nil {
		return "", errors.New("wechaty is not started")
	}
	room := b.Room().Load(groupID)
	if err := room.Ready(false); err != nil {
		return "", err
	}
	msg, err := room.Say(text)
	if err != nil {
		return "", err
	}
	if msg == nil {
		return "", nil
	}
	if identified, ok := msg.(interface{ ID() string }); ok {
		return identified.ID(), nil
	}
	return "", nil
}
