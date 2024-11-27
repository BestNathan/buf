// Copyright 2020-2024 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package buflsp implements a language server for Protobuf.
//
// The main entry-point of this package is the Serve() function, which creates a new LSP server.
package buflsp

import (
	"context"
	"log/slog"

	"github.com/bufbuild/buf/private/pkg/slogext"
	"go.lsp.dev/jsonrpc2"
)

// wrapReplier wraps a jsonrpc2.Replier, allowing us to inject logging and tracing and so on.
func (l *lsp) wrapReplier(reply jsonrpc2.Replier, req jsonrpc2.Request) jsonrpc2.Replier {
	return func(ctx context.Context, result any, err error) error {
		if err != nil {
			l.logger.Warn(
				"JSON-RPC Responding Error",
				slog.String("Method", req.Method()),
				slogext.ErrorAttr(err),
			)
		} else {
			l.logger.Debug(
				"JSON-RPC Responding",
				slog.String("Method", req.Method()),
				slog.Any("Params", result),
			)
		}

		return reply(ctx, result, err)
	}
}

// connWrapper wraps a connection and logs calls and notifications.
//
// By default, the ClientDispatcher does not log the bodies of requests and responses, making
// for much lower-quality debugging.
type connWrapper struct {
	jsonrpc2.Conn

	logger *slog.Logger
}

func (c *connWrapper) Call(
	ctx context.Context, method string, params, result any) (id jsonrpc2.ID, err error) {
	c.logger.Debug(
		"JSON-RPC Call",
		slog.String("Method", method),
		slog.Any("Params", params),
	)

	id, err = c.Conn.Call(ctx, method, params, result)
	if err != nil {
		c.logger.Warn(
			"JSON-RPC Call Fail",
			slog.String("Method", method),
			slogext.ErrorAttr(err),
		)
	} else {
		c.logger.Debug(
			"JSON-RPC Call Success",
			slog.String("Method", method),
			slog.Any("Result", result),
		)
	}

	return
}

func (c *connWrapper) Notify(
	ctx context.Context, method string, params any) error {
	c.logger.Debug(
		"JSON-RPC Notify",
		slog.String("Method", method),
		slog.Any("Params", params),
	)

	err := c.Conn.Notify(ctx, method, params)
	if err != nil {
		c.logger.Warn(
			"JSON-RPC Notify Fail",
			slog.String("Method", method),
			slogext.ErrorAttr(err),
		)
	} else {
		c.logger.Debug(
			"JSON-RPC Notify Success",
			slog.String("Method", method),
		)
	}

	return err
}
