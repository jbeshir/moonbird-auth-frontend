package aengine

import "context"

type PermissionChecker interface {
	CheckRead(ctx context.Context, kind, key string) (bool, error)
	CheckWrite(ctx context.Context, kind, key string) (bool, error)
}
