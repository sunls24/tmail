package api

import (
	"context"
	"tmail/config"
	"tmail/ent"

	"github.com/sunls24/gox/server"
)

type (
	dbKey     struct{}
	configKey struct{}
)

func ServerContext(srv *server.Server, cfg *config.Config, db *ent.Client) context.Context {
	srv.ContextValue(dbKey{}, db)
	srv.ContextValue(configKey{}, cfg)
	return srv.NewValueContext()
}

func Config(ctx context.Context) *config.Config {
	return ctx.Value(configKey{}).(*config.Config)
}

func DB(ctx context.Context) *ent.Client {
	return ctx.Value(dbKey{}).(*ent.Client)
}
