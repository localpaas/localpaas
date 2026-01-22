package nginx

import (
	crossplane "github.com/localpaas/nginx-go-crossplane"
)

type Directive struct {
	*crossplane.Directive
}

func NewDirective(name string, args []string) *Directive {
	return &Directive{
		Directive: &crossplane.Directive{
			Directive: name,
			Args:      args,
		},
	}
}
