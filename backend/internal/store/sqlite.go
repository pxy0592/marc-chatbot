package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"

	"github.com/marc-pango/marc-chatbot/backend/internal/domain"
)

var ErrNotFound = errors.New("not found")

type Store struct{ db *sql.DB }

func Open(path string) (*Store, error) {
	if path != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, err
		}
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	s := &Store{db: db}
	if err := s.migrate(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate(ctx context.Context) error {
	const schema = `
PRAGMA journal_mode=WAL;
PRAGMA foreign_keys=ON;
CREATE TABLE IF NOT EXISTS bot_instances (
 id TEXT PRIMARY KEY, public_account_id TEXT, display_name TEXT, status TEXT NOT NULL,
 status_reason TEXT, qr_code TEXT, qr_expires_at TEXT, last_seen_at TEXT, updated_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS binding_sessions (
 id TEXT PRIMARY KEY, status TEXT NOT NULL, qr_code TEXT, created_at TEXT NOT NULL,
 expires_at TEXT NOT NULL, completed_at TEXT, failure_reason TEXT
);
CREATE TABLE IF NOT EXISTS groups (
 id TEXT PRIMARY KEY, provider_group_id TEXT NOT NULL UNIQUE, name TEXT NOT NULL,
 member_count INTEGER, enabled INTEGER NOT NULL DEFAULT 0, available INTEGER NOT NULL DEFAULT 1,
 discovered_at TEXT NOT NULL, updated_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS group_messages (
 id TEXT PRIMARY KEY, provider_message_id TEXT NOT NULL UNIQUE, group_id TEXT NOT NULL,
 sender_id TEXT NOT NULL, sender_name TEXT NOT NULL, direction TEXT NOT NULL,
 message_type TEXT NOT NULL, content TEXT NOT NULL, occurred_at TEXT NOT NULL,
 received_at TEXT NOT NULL, processing_status TEXT NOT NULL, self_message INTEGER NOT NULL,
 FOREIGN KEY(group_id) REFERENCES groups(id)
);
CREATE INDEX IF NOT EXISTS idx_messages_group_received ON group_messages(group_id, received_at DESC);
CREATE TABLE IF NOT EXISTS outbound_requests (
 id TEXT PRIMARY KEY, idempotency_key TEXT NOT NULL UNIQUE, group_id TEXT NOT NULL,
 content TEXT NOT NULL, requested_by TEXT NOT NULL, status TEXT NOT NULL,
 provider_message_id TEXT, failure_reason TEXT, created_at TEXT NOT NULL, updated_at TEXT NOT NULL,
 FOREIGN KEY(group_id) REFERENCES groups(id)
);
CREATE TABLE IF NOT EXISTS audit_events (
 id TEXT PRIMARY KEY, actor TEXT NOT NULL, action TEXT NOT NULL, target_type TEXT NOT NULL,
 target_id TEXT NOT NULL, result TEXT NOT NULL, detail TEXT, created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_audit_created ON audit_events(created_at DESC);
`
	_, err := s.db.ExecContext(ctx, schema)
	return err
}

func ts(t time.Time) string { return t.UTC().Format(time.RFC3339Nano) }
func nts(t *time.Time) any {
	if t == nil {
		return nil
	}
	return ts(*t)
}
func ns(v *string) any {
	if v == nil {
		return nil
	}
	return *v
}
func parseTime(v string) time.Time { t, _ := time.Parse(time.RFC3339Nano, v); return t }
func parseTimePtr(v sql.NullString) *time.Time {
	if !v.Valid {
		return nil
	}
	t := parseTime(v.String)
	return &t
}
func strPtr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	x := v.String
	return &x
}
func boolInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func (s *Store) GetBot(ctx context.Context) (domain.BotInstance, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, public_account_id, display_name, status, status_reason, qr_code, qr_expires_at, last_seen_at, updated_at FROM bot_instances WHERE id='primary'`)
	var b domain.BotInstance
	var pub, name, reason, qr, qrExp, last sql.NullString
	var updated string
	if err := row.Scan(&b.ID, &pub, &name, &b.Status, &reason, &qr, &qrExp, &last, &updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.BotInstance{}, ErrNotFound
		}
		return domain.BotInstance{}, err
	}
	b.PublicAccountID, b.DisplayName, b.StatusReason, b.QRCode = strPtr(pub), strPtr(name), strPtr(reason), strPtr(qr)
	b.QRExpiresAt, b.LastSeenAt, b.UpdatedAt = parseTimePtr(qrExp), parseTimePtr(last), parseTime(updated)
	return b, nil
}

func (s *Store) SaveBot(ctx context.Context, b domain.BotInstance) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO bot_instances(id,public_account_id,display_name,status,status_reason,qr_code,qr_expires_at,last_seen_at,updated_at)
VALUES(?,?,?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET public_account_id=excluded.public_account_id,display_name=excluded.display_name,status=excluded.status,status_reason=excluded.status_reason,qr_code=excluded.qr_code,qr_expires_at=excluded.qr_expires_at,last_seen_at=excluded.last_seen_at,updated_at=excluded.updated_at`,
		b.ID, ns(b.PublicAccountID), ns(b.DisplayName), b.Status, ns(b.StatusReason), ns(b.QRCode), nts(b.QRExpiresAt), nts(b.LastSeenAt), ts(b.UpdatedAt))
	return err
}

func (s *Store) SaveBinding(ctx context.Context, b domain.BindingSession) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO binding_sessions(id,status,qr_code,created_at,expires_at,completed_at,failure_reason) VALUES(?,?,?,?,?,?,?)
ON CONFLICT(id) DO UPDATE SET status=excluded.status,qr_code=excluded.qr_code,expires_at=excluded.expires_at,completed_at=excluded.completed_at,failure_reason=excluded.failure_reason`,
		b.ID, b.Status, ns(b.QRCode), ts(b.CreatedAt), ts(b.ExpiresAt), nts(b.CompletedAt), ns(b.FailureReason))
	return err
}

func (s *Store) CurrentBinding(ctx context.Context) (domain.BindingSession, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id,status,qr_code,created_at,expires_at,completed_at,failure_reason FROM binding_sessions WHERE status IN ('pending','scanned') ORDER BY created_at DESC LIMIT 1`)
	var b domain.BindingSession
	var qr, completed, failure sql.NullString
	var created, expires string
	if err := row.Scan(&b.ID, &b.Status, &qr, &created, &expires, &completed, &failure); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return b, ErrNotFound
		}
		return b, err
	}
	b.QRCode, b.CreatedAt, b.ExpiresAt, b.CompletedAt, b.FailureReason = strPtr(qr), parseTime(created), parseTime(expires), parseTimePtr(completed), strPtr(failure)
	return b, nil
}

func (s *Store) UpsertGroup(ctx context.Context, g domain.Group) (domain.Group, error) {
	_, err := s.db.ExecContext(ctx, `INSERT INTO groups(id,provider_group_id,name,member_count,enabled,available,discovered_at,updated_at) VALUES(?,?,?,?,?,?,?,?)
ON CONFLICT(provider_group_id) DO UPDATE SET name=excluded.name,member_count=excluded.member_count,available=excluded.available,updated_at=excluded.updated_at`,
		g.ID, g.ProviderGroupID, g.Name, g.MemberCount, boolInt(g.Enabled), boolInt(g.Available), ts(g.DiscoveredAt), ts(g.UpdatedAt))
	if err != nil {
		return g, err
	}
	return s.GetGroupByProviderID(ctx, g.ProviderGroupID)
}

func scanGroup(scanner interface{ Scan(...any) error }) (domain.Group, error) {
	var g domain.Group
	var count sql.NullInt64
	var enabled, available int
	var discovered, updated string
	if err := scanner.Scan(&g.ID, &g.ProviderGroupID, &g.Name, &count, &enabled, &available, &discovered, &updated); err != nil {
		return g, err
	}
	if count.Valid {
		n := int(count.Int64)
		g.MemberCount = &n
	}
	g.Enabled, g.Available, g.DiscoveredAt, g.UpdatedAt = enabled == 1, available == 1, parseTime(discovered), parseTime(updated)
	return g, nil
}

const groupColumns = `id,provider_group_id,name,member_count,enabled,available,discovered_at,updated_at`

func (s *Store) GetGroup(ctx context.Context, id string) (domain.Group, error) {
	g, err := scanGroup(s.db.QueryRowContext(ctx, `SELECT `+groupColumns+` FROM groups WHERE id=?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return g, ErrNotFound
	}
	return g, err
}
func (s *Store) GetGroupByProviderID(ctx context.Context, id string) (domain.Group, error) {
	g, err := scanGroup(s.db.QueryRowContext(ctx, `SELECT `+groupColumns+` FROM groups WHERE provider_group_id=?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return g, ErrNotFound
	}
	return g, err
}
func (s *Store) SetGroupEnabled(ctx context.Context, id string, enabled bool) (domain.Group, error) {
	res, err := s.db.ExecContext(ctx, `UPDATE groups SET enabled=?,updated_at=? WHERE id=?`, boolInt(enabled), ts(time.Now()), id)
	if err != nil {
		return domain.Group{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.Group{}, ErrNotFound
	}
	return s.GetGroup(ctx, id)
}
func (s *Store) ListGroups(ctx context.Context, enabled *bool) ([]domain.Group, error) {
	q := `SELECT ` + groupColumns + ` FROM groups`
	args := []any{}
	if enabled != nil {
		q += ` WHERE enabled=?`
		args = append(args, boolInt(*enabled))
	}
	q += ` ORDER BY name`
	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Group{}
	for rows.Next() {
		g, err := scanGroup(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

func (s *Store) InsertMessage(ctx context.Context, m domain.GroupMessage) (bool, error) {
	res, err := s.db.ExecContext(ctx, `INSERT OR IGNORE INTO group_messages(id,provider_message_id,group_id,sender_id,sender_name,direction,message_type,content,occurred_at,received_at,processing_status,self_message) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`, m.ID, m.ProviderMessageID, m.GroupID, m.SenderID, m.SenderName, m.Direction, m.MessageType, m.Content, ts(m.OccurredAt), ts(m.ReceivedAt), m.ProcessingStatus, boolInt(m.SelfMessage))
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n == 1, nil
}
func (s *Store) ListMessages(ctx context.Context, groupID string, direction domain.MessageDirection, limit int) ([]domain.GroupMessage, error) {
	q := `SELECT id,provider_message_id,group_id,sender_id,sender_name,direction,message_type,content,occurred_at,received_at,processing_status,self_message FROM group_messages WHERE 1=1`
	args := []any{}
	if groupID != "" {
		q += ` AND group_id=?`
		args = append(args, groupID)
	}
	if direction != "" {
		q += ` AND direction=?`
		args = append(args, direction)
	}
	q += ` ORDER BY received_at DESC LIMIT ?`
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.GroupMessage{}
	for rows.Next() {
		var m domain.GroupMessage
		var occurred, received string
		var self int
		if err := rows.Scan(&m.ID, &m.ProviderMessageID, &m.GroupID, &m.SenderID, &m.SenderName, &m.Direction, &m.MessageType, &m.Content, &occurred, &received, &m.ProcessingStatus, &self); err != nil {
			return nil, err
		}
		m.OccurredAt, m.ReceivedAt, m.SelfMessage = parseTime(occurred), parseTime(received), self == 1
		out = append(out, m)
	}
	return out, rows.Err()
}

func (s *Store) GetOutboundByKey(ctx context.Context, key string) (domain.OutboundMessageRequest, error) {
	return s.scanOutbound(s.db.QueryRowContext(ctx, `SELECT id,idempotency_key,group_id,content,requested_by,status,provider_message_id,failure_reason,created_at,updated_at FROM outbound_requests WHERE idempotency_key=?`, key))
}
func (s *Store) scanOutbound(scanner interface{ Scan(...any) error }) (domain.OutboundMessageRequest, error) {
	var o domain.OutboundMessageRequest
	var provider, failure sql.NullString
	var created, updated string
	err := scanner.Scan(&o.ID, &o.IdempotencyKey, &o.GroupID, &o.Content, &o.RequestedBy, &o.Status, &provider, &failure, &created, &updated)
	if errors.Is(err, sql.ErrNoRows) {
		return o, ErrNotFound
	}
	if err != nil {
		return o, err
	}
	o.ProviderMessageID, o.FailureReason, o.CreatedAt, o.UpdatedAt = strPtr(provider), strPtr(failure), parseTime(created), parseTime(updated)
	return o, nil
}
func (s *Store) SaveOutbound(ctx context.Context, o domain.OutboundMessageRequest) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO outbound_requests(id,idempotency_key,group_id,content,requested_by,status,provider_message_id,failure_reason,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?)
ON CONFLICT(id) DO UPDATE SET status=excluded.status,provider_message_id=excluded.provider_message_id,failure_reason=excluded.failure_reason,updated_at=excluded.updated_at`, o.ID, o.IdempotencyKey, o.GroupID, o.Content, o.RequestedBy, o.Status, ns(o.ProviderMessageID), ns(o.FailureReason), ts(o.CreatedAt), ts(o.UpdatedAt))
	return err
}

func (s *Store) InsertAudit(ctx context.Context, a domain.AuditEvent) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO audit_events(id,actor,action,target_type,target_id,result,detail,created_at) VALUES(?,?,?,?,?,?,?,?)`, a.ID, a.Actor, a.Action, a.TargetType, a.TargetID, a.Result, ns(a.Detail), ts(a.CreatedAt))
	return err
}
func (s *Store) ListAudit(ctx context.Context, limit int) ([]domain.AuditEvent, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,actor,action,target_type,target_id,result,detail,created_at FROM audit_events ORDER BY created_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.AuditEvent{}
	for rows.Next() {
		var a domain.AuditEvent
		var detail sql.NullString
		var created string
		if err := rows.Scan(&a.ID, &a.Actor, &a.Action, &a.TargetType, &a.TargetID, &a.Result, &detail, &created); err != nil {
			return nil, err
		}
		a.Detail, a.CreatedAt = strPtr(detail), parseTime(created)
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *Store) DebugCount(ctx context.Context, table string) (int, error) {
	allowed := map[string]bool{"group_messages": true, "outbound_requests": true, "audit_events": true}
	if !allowed[table] {
		return 0, fmt.Errorf("unsupported table")
	}
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+table).Scan(&n)
	return n, err
}
