// Copyright 2026 Novin Rasan Authors
//  All rights reserved.
//
// Novin Rasan offline-push queue DAO.
//

package dao

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/zeromicro/go-zero/core/logx"
)

// PushQueueDao persists pushes for sessions that are currently offline so the
// client can fetch them on the next reconnect.
type PushQueueDao struct {
	db *sql.DB
}

// NewPushQueueDao opens the push_queue database. Returns nil if dsn is empty or
// the connection cannot be established (offline-push simply stays disabled).
func NewPushQueueDao(dsn string) *PushQueueDao {
	if dsn == "" {
		logx.Info("pushQueue: disabled (no DSN configured)")
		return nil
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logx.Errorf("pushQueue: open error: %v", err)
		return nil
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Hour)

	if err = db.Ping(); err != nil {
		logx.Errorf("pushQueue: ping error: %v", err)
		_ = db.Close()
		return nil
	}
	logx.Info("pushQueue: initialized")
	return &PushQueueDao{db: db}
}

// Enabled reports whether the DAO has a live DB pool.
func (d *PushQueueDao) Enabled() bool {
	return d != nil && d.db != nil
}

// SaveOfflinePush inserts a pending push row for the given user.
func (d *PushQueueDao) SaveOfflinePush(ctx context.Context, userId, authKeyId, sessionId, msgId int64, data []byte) error {
	if d == nil || d.db == nil {
		return nil
	}
	_, err := d.db.ExecContext(ctx,
		"INSERT INTO push_queue (user_id, auth_key_id, session_id, msg_id, data, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		userId, authKeyId, sessionId, msgId, data, time.Now(),
	)
	if err != nil {
		logx.WithContext(ctx).Errorf("pushQueue: insert error (user_id: %d, msg_id: %d): %v", userId, msgId, err)
	}
	return err
}

// FetchPending returns pending (undelivered) pushes for a user, oldest first.
func (d *PushQueueDao) FetchPending(ctx context.Context, userId int64, limit int) ([]*PendingPush, error) {
	if d == nil || d.db == nil {
		return nil, nil
	}
	rows, err := d.db.QueryContext(ctx,
		"SELECT id, user_id, auth_key_id, session_id, msg_id, data FROM push_queue WHERE user_id = ? AND delivered_at IS NULL ORDER BY id ASC LIMIT ?",
		userId, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*PendingPush
	for rows.Next() {
		p := &PendingPush{}
		if err = rows.Scan(&p.Id, &p.UserId, &p.AuthKeyId, &p.SessionId, &p.MsgId, &p.Data); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// MarkDelivered flags the given push rows as delivered.
func (d *PushQueueDao) MarkDelivered(ctx context.Context, ids []int64) error {
	if d == nil || d.db == nil || len(ids) == 0 {
		return nil
	}
	now := time.Now()
	for _, id := range ids {
		if _, err := d.db.ExecContext(ctx,
			"UPDATE push_queue SET delivered_at = ? WHERE id = ?", now, id); err != nil {
			logx.WithContext(ctx).Errorf("pushQueue: mark delivered error (id: %d): %v", id, err)
		}
	}
	return nil
}

// PendingPush is one row of the push_queue table.
type PendingPush struct {
	Id        int64
	UserId    int64
	AuthKeyId  int64
	SessionId int64
	MsgId     int64
	Data      []byte
}
