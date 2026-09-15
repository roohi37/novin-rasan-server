/*
 * Created from 'scheme.tl' by 'mtprotoc'
 *
 * Copyright (c) 2021-present,  Teamgram Studio (https://teamgram.io).
 *  All rights reserved.
 *
 * Author: teamgramio (teamgram.io@gmail.com)
 */

package config

import (
	"github.com/teamgram/teamgram-server/pkg/conf"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	AuthSession      zrpc.RpcClientConf
	StatusClient     zrpc.RpcClientConf
	GatewayClient    zrpc.RpcClientConf
	BFFProxyClients  conf.BFFProxyClients
	UseStreamGateway bool `json:",default=false"`

	// PushQueueDSN - Novin Rasan: DSN of the push_queue table used to persist
	// offline pushes. Empty string disables the offline-push feature.
	PushQueueDSN string `json:",optional"`
}

// Routine routine.
type Routine struct {
	Size uint64
	Chan uint64
}
