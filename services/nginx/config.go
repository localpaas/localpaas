package nginx

import (
	crossplane "github.com/nginxinc/nginx-go-crossplane"
	"github.com/tiendc/gofn"
)

type Config struct {
	inner *crossplane.Config
}

func NewConfig(directives ...*crossplane.Directive) *Config {
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
	return blocksByName(c.inner.Parsed, name, n)
}

func (c *Config) IterBlocksByName(name string, fn func(block *Block, i int) bool) {
	for i, block := range c.BlocksByName(name, -1) {
		if !fn(block, i) {
			return
		}
	}
}

func (c *Config) AddBlock(block *Block) {
	c.inner.Parsed = append(c.inner.Parsed, block.inner)
}

func (c *Config) AllDirectives() crossplane.Directives {
	return c.inner.Parsed
}

func (c *Config) IterAllDirectives(fn func(block *crossplane.Directive, i int) bool) {
	for i, directive := range c.AllDirectives() {
		if !fn(directive, i) {
			return
		}
	}
}

func (c *Config) Directives() crossplane.Directives {
	return directivesByName(c.inner.Parsed, "", -1)
}

func (c *Config) IterDirectives(fn func(block *crossplane.Directive, i int) bool) {
	for i, directive := range c.Directives() {
		if !fn(directive, i) {
			return
		}
	}
}

func (c *Config) DirectivesByName(name string, n int) crossplane.Directives {
	return directivesByName(c.inner.Parsed, name, n)
}

func (c *Config) IterDirectivesByName(name string, fn func(block *crossplane.Directive, i int) bool) {
	for i, directive := range c.DirectivesByName(name, -1) {
		if !fn(directive, i) {
			return
		}
	}
}

func (c *Config) AddDirectives(directives ...*crossplane.Directive) {
	c.inner.Parsed = append(c.inner.Parsed, directives...)
}

type Block struct {
	inner *crossplane.Directive
}

func NewBlock(name string, args []string, directives ...*crossplane.Directive) *Block {
	block := &Block{
		inner: &crossplane.Directive{
			Directive: name,
			Args:      args,
		},
	}
	block.AddDirectives(directives...)
	return block
}

func NewServerBlock(directives ...*crossplane.Directive) *Block {
	return NewBlock("server", nil, directives...)
}

func (c *Block) Blocks() []*Block {
	return blocksByName(c.inner.Block, "", 0)
}

func (c *Block) BlocksByName(name string, n int) []*Block {
	return blocksByName(c.inner.Block, name, n)
}

func (c *Block) IterBlocksByName(name string, fn func(block *Block, i int) bool) {
	for i, block := range c.BlocksByName(name, -1) {
		if !fn(block, i) {
			return
		}
	}
}

func (c *Block) AddBlock(block *Block) {
	c.inner.Block = append(c.inner.Block, block.inner)
}

func (c *Block) RemoveBlock(block *Block) {
	gofn.Remove(&c.inner.Block, block.inner)
}

func (c *Block) AllDirectives() crossplane.Directives {
	return c.inner.Block
}

func (c *Block) IterAllDirectives(fn func(block *crossplane.Directive, i int) bool) {
	for i, directive := range c.AllDirectives() {
		if !fn(directive, i) {
			return
		}
	}
}

func (c *Block) Directives() crossplane.Directives {
	return directivesByName(c.inner.Block, "", -1)
}

func (c *Block) IterDirectives(fn func(block *crossplane.Directive, i int) bool) {
	for i, directive := range c.Directives() {
		if !fn(directive, i) {
			return
		}
	}
}

func (c *Block) DirectivesByName(name string, n int) crossplane.Directives {
	return directivesByName(c.inner.Block, name, n)
}

func (c *Block) IterDirectivesByName(name string, fn func(block *crossplane.Directive, i int) bool) {
	for i, directive := range c.DirectivesByName(name, -1) {
		if !fn(directive, i) {
			return
		}
	}
}

func (c *Block) AddDirectives(directives ...*crossplane.Directive) {
	c.inner.Block = append(c.inner.Block, directives...)
}

func blocksByName(directives crossplane.Directives, name string, n int) (blocks []*Block) {
	for _, dir := range directives {
		if !dir.IsBlock() {
			continue
		}
		if name != "" && dir.Directive != name {
			continue
		}
		if n > 0 && len(blocks) >= n {
			break
		}
		blocks = append(blocks, &Block{inner: dir})
	}
	return
}

func directivesByName(directives crossplane.Directives, name string, n int) (res crossplane.Directives) {
	for _, dir := range directives {
		if dir.IsBlock() {
			continue
		}
		if name != "" && dir.Directive != name {
			continue
		}
		if n > 0 && len(res) >= n {
			break
		}
		res = append(res, dir)
	}
	return
}
