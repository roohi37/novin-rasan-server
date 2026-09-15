/*
 * Novin Rasan: persistent offline push queue DAO (sync service).
 *
 * Copyright (c) 2021-present, Teamgram Studio (https://teamgram.io).
 * All rights reserved.
 */

package dao

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	_ "github.com/go-sql-driver/mysql"
)

// PushQueueDao stores pending pushes for users that have no online session.
type PushQueueDao struct {
	db  *sql.DB
	dsn string
}

// NewPushQueueDao opens a MySQL pool for the push_queue table. When dsn is
// empty the DAO is disabled and every method is a safe no-op.
func NewPushQueueDao(dsn string) *PushQueueDao {
	d := &PushQueueDao{dsn: dsn}
	if dsn == "" {
		logx.Info("pushQueue(sync): disabled (empty PushQueueDSN)")
		return d
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logx.Errorf("pushQueue(sync): open error: %v", err)
		return d
	}
	db.SetMaxOpenConns(8)
	db.SetMaxIdleConns(4)
	db.SetConnMaxLifetime(time.Hour)
	if err = db.Ping(); err != nil {
		logx.Errorf("pushQueue(sync): ping error: %v", err)
		return d
	}
	d.db = db
	logx.Info("pushQueue(sync): initialized")
	return d
}

// Enabled reports whether the DAO has a live DB pool.
func (d *PushQueueDao) Enabled() bool {
	return d != nil && d.db != nil
}

// SaveOfflinePush persists one push payload for later replay.
func (d *PushQueueDao) SaveOfflinePush(ctx context.Context, userId, authKeyId, sessionId, msgId int64, data []byte) error {
	if !d.Enabled() {
		return nil
	}
	_, err := d.db.ExecContext(ctx,
		"INSERT INTO push_queue (user_id, auth_key_id, session_id, msg_id, data, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		userId, authKeyId, sessionId, msgId, data, time.Now())
	if err != nil {
		logx.WithContext(ctx).Errorf("pushQueue(sync): save error: %v", err)
	}
	return err
}

// PendingPush is one row of the push_queue table.
type PendingPush struct {
	Id         int64
	UserId     int64
	AuthKeyId  int64
	SessionId  int64
	MsgId      int64
	Data       []byte
	CreatedAt  time.Time
	DeliveredAt sql.NullTime
}

// FetchPending returns up to limit undelivered pushes for a user.
func (d *PushQueueDao) FetchPending(ctx context.Context, userId int64, limit int) ([]*PendingPush, error) {
	if !d.Enabled() {
		return nil, nil
	}
	rows, err := d.db.QueryContext(ctx,
		"SELECT id, user_id, auth_key_id, session_id, msg_id, data, created_at FROM push_queue WHERE user_id = ? AND delivered_at IS NULL ORDER BY id ASC LIMIT ?",
		userId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*PendingPush
	for rows.Next() {
		p := &PendingPush{}
		if err = rows.Scan(&p.Id, &p.UserId, &p.AuthKeyId, &p.SessionId, &p.MsgId, &p.Data, &p.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// MarkDelivered marks a queued push as delivered.
func (d *PushQueueDao) MarkDelivered(ctx context.Context, id int64) error {
	if !d.Enabled() {
		return nil
	}
	_, err := d.db.ExecContext(ctx,
		"UPDATE push_queue SET delivered_at = ? WHERE id = ?", time.Now(), id)
	return err
}
