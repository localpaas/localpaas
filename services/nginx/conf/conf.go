package conf

import (
	"github.com/tufanbarisyildirim/gonginx/parser"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

func Parse(filepath string) (*parser.Parser, error) {
	p, err := parser.NewParser(filepath)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return p, nil
}
