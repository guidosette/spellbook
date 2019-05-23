package page

import "context"

type Extender interface {
	BeforeCreate(ctx context.Context, resource Resource) error
}
