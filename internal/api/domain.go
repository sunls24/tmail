package api

import "context"

func DomainList(ctx context.Context) ([]string, error) {
	return Config(ctx).DomainList, nil
}
