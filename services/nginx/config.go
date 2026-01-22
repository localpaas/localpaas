package nginx

import (
	crossplane "github.com/localpaas/nginx-go-crossplane"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

type Config struct {
	inner *crossplane.Config
}

func NewConfig(directives ...*Directive) *Config {
	cfg := &Config{
		inner: &crossplane.Config{},
	}
	cfg.AddDirectives(directives...)
	return cfg
}

func (c *Config) Blocks() []*Block {
	return c.BlocksByName("", 0)
}

func (c *Config) BlocksByName(name string, n int) []*Block {
	return blocksByName(c.allCrossplaneDirectives(), name, n)
}

func (c *Config) IterBlocksByName(name string, fn func(*Block, int) (bool, error)) error {
	return iterBlocksByName(c.allCrossplaneDirectives(), name, fn)
}

func (c *Config) AddBlock(block *Block) {
	c.inner.Parsed = append(c.inner.Parsed, block.Directive)
}

func (c *Config) allCrossplaneDirectives() crossplane.Directives {
	return c.inner.Parsed
}

func (c *Config) AllDirectives() []*Directive {
	return convertDirectives(c.allCrossplaneDirectives())
}

func (c *Config) IterAllDirectives(fn func(*Directive, int) (bool, error)) error {
	for i, dir := range c.allCrossplaneDirectives() {
		shouldContinue, err := fn(&Directive{Directive: dir}, i)
		if err != nil {
			return apperrors.Wrap(err)
		}
		if !shouldContinue {
			break
		}
	}
	return nil
}

func (c *Config) Directives() []*Directive {
	return directivesByName(c.allCrossplaneDirectives(), "", -1)
}

func (c *Config) IterDirectives(fn func(*Directive, int) (bool, error)) error {
	return iterDirectivesByName(c.allCrossplaneDirectives(), "", fn)
}

func (c *Config) DirectivesByName(name string, n int) []*Directive {
	return directivesByName(c.inner.Parsed, name, n)
}

func (c *Config) IterDirectivesByName(name string, fn func(*Directive, int) (bool, error)) error {
	return iterDirectivesByName(c.allCrossplaneDirectives(), name, fn)
}

func (c *Config) AddDirectives(directives ...*Directive) {
	for _, dir := range directives {
		c.inner.Parsed = append(c.inner.Parsed, dir.Directive)
	}
}
