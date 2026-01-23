package nginx

import (
	crossplane "github.com/localpaas/nginx-go-crossplane"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

type Block struct {
	*crossplane.Directive
}

func NewBlock(name string, args []string, directives ...*Directive) *Block {
	block := &Block{
		Directive: &crossplane.Directive{
			Directive: name,
			Args:      args,
		},
	}
	block.AddDirectives(directives...)
	return block
}

func NewServerBlock(directives ...*Directive) *Block {
	return NewBlock("server", nil, directives...)
}

func (b *Block) Blocks() []*Block {
	return blocksByName(b.Block, "", 0)
}

func (b *Block) BlocksByName(name string, n int) []*Block {
	return blocksByName(b.Block, name, n)
}

func (b *Block) IterBlocksByName(name string, fn func(*Block, int) (bool, error)) error {
	return iterBlocksByName(b.Block, name, fn)
}

func (b *Block) AddBlock(block *Block) {
	b.Block = append(b.Block, block.Directive)
}

// AddLocationBlock adds `location` block to a `server` block.
// If the block is not a `server` one, an error will be returned.
func (b *Block) AddLocationBlock(location string, args []string) (*Block, error) {
	if b.Directive == nil || b.Directive.Directive != "server" {
		return nil, apperrors.Wrap(ErrServerBlockRequired)
	}
	for _, block := range b.BlocksByName("location", -1) {
		if block.Args[0] == location {
			return block, nil
		}
	}
	block := NewBlock("location", append([]string{location}, args...))
	b.AddBlock(block)
	return block, nil
}

func (b *Block) GetBlock(name string, args []string) *Block {
	idx := b.GetBlockIndex(name, args)
	if idx >= 0 {
		return &Block{Directive: b.Block[idx]}
	}
	return nil
}

func (b *Block) GetBlockIndex(name string, args []string) int {
	for i, directive := range b.Block {
		if directive.IsBlock() && directive.Directive == name && gofn.Equal(directive.Args, args) {
			return i
		}
	}
	return -1
}

func (b *Block) RemoveBlock(name string, args []string) {
	idx := b.GetBlockIndex(name, args)
	if idx >= 0 {
		gofn.RemoveAt(&b.Block, idx)
	}
}

func (b *Block) AllDirectives() []*Directive {
	return convertDirectives(b.Block)
}

func (b *Block) IterAllDirectives(fn func(*Directive, int) (bool, error)) error {
	for i, dir := range b.Block {
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

func (b *Block) Directives() []*Directive {
	return directivesByName(b.Block, "", -1)
}

func (b *Block) IterDirectives(fn func(*Directive, int) (bool, error)) error {
	return iterDirectivesByName(b.Block, "", fn)
}

func (b *Block) DirectivesByName(name string, n int) []*Directive {
	return directivesByName(b.Block, name, n)
}

func (b *Block) IterDirectivesByName(name string, fn func(*Directive, int) (bool, error)) error {
	return iterDirectivesByName(b.Block, name, fn)
}

func (b *Block) AddDirectives(directives ...*Directive) {
	for _, dir := range directives {
		b.Block = append(b.Block, dir.Directive)
	}
}

func (b *Block) AddDirectiveNamesIfNotExist(names ...string) {
	for _, name := range names {
		if b.ContainPartialMatchedDirective(name, nil) {
			continue
		}
		b.AddDirectives(NewDirective(name, nil))
	}
}

func (b *Block) AddDirectivesIfNotExist(directives ...*Directive) {
	for _, directive := range directives {
		if b.ContainDirective(directive) {
			continue
		}
		b.AddDirectives(directive)
	}
}

func (b *Block) SetDirectiveArgs(name string, args []string, n int) {
	for _, directive := range b.DirectivesByName(name, n) {
		directive.Args = args
	}
}

func (b *Block) ContainDirective(dir *Directive) bool {
	return b.GetDirectiveIndex(dir.Directive.Directive, dir.Args) != -1
}

func (b *Block) GetDirective(name string, args []string) *Directive {
	idx := b.GetDirectiveIndex(name, args)
	if idx >= 0 {
		return &Directive{Directive: b.Block[idx]}
	}
	return nil
}

func (b *Block) GetDirectiveIndex(name string, args []string) int {
	for i, directive := range b.Block {
		if !directive.IsBlock() && directive.Directive == name && gofn.Equal(directive.Args, args) {
			return i
		}
	}
	return -1
}

func (b *Block) ContainPartialMatchedDirective(name string, args []string) bool {
	return b.GetPartialMatchedDirectiveIndex(name, args) != -1
}

func (b *Block) GetPartialMatchedDirective(name string, args []string) *Directive {
	idx := b.GetPartialMatchedDirectiveIndex(name, args)
	if idx >= 0 {
		return &Directive{Directive: b.Block[idx]}
	}
	return nil
}

func (b *Block) GetPartialMatchedDirectiveIndex(name string, args []string) int {
	partialMatchIdx := -1
	for i, dir := range b.Block {
		if dir.IsBlock() || dir.Directive != name {
			continue
		}
		if gofn.Equal(dir.Args, args) { // Full matching
			return i
		}
		if partialMatchIdx == -1 && len(dir.Args) > len(args) && gofn.Equal(dir.Args[:len(args)], args) {
			partialMatchIdx = i
		}
	}
	return partialMatchIdx
}

func (b *Block) RemoveDirectives(directives ...*Directive) {
	for _, dir := range directives {
		idx := b.GetDirectiveIndex(dir.Directive.Directive, dir.Args)
		if idx >= 0 {
			gofn.RemoveAt(&b.Block, idx)
		}
	}
}

func (b *Block) RemovePartialMatchedDirectives(directives ...*Directive) {
	for _, dir := range directives {
		idx := b.GetPartialMatchedDirectiveIndex(dir.Directive.Directive, dir.Args)
		if idx >= 0 {
			gofn.RemoveAt(&b.Block, idx)
		}
	}
}

func (b *Block) RemoveDirectiveNames(names ...string) {
	for _, name := range names {
		idx := b.GetPartialMatchedDirectiveIndex(name, nil)
		if idx >= 0 {
			gofn.RemoveAt(&b.Block, idx)
		}
	}
}

func blocksByName(directives crossplane.Directives, name string, n int) (blocks []*Block) {
	blocks = make([]*Block, 0, len(directives))
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
		blocks = append(blocks, &Block{Directive: dir})
	}
	return
}

func iterBlocksByName(
	directives crossplane.Directives,
	name string,
	fn func(*Block, int) (bool, error),
) error {
	count := 0
	for _, dir := range directives {
		if !dir.IsBlock() {
			continue
		}
		if name != "" && dir.Directive != name {
			continue
		}
		shouldContinue, err := fn(&Block{Directive: dir}, count)
		if err != nil {
			return apperrors.Wrap(err)
		}
		if !shouldContinue {
			break
		}
		count++
	}
	return nil
}

func directivesByName(directives crossplane.Directives, name string, n int) (res []*Directive) {
	res = make([]*Directive, 0, len(directives))
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
		res = append(res, &Directive{Directive: dir})
	}
	return
}

func iterDirectivesByName(
	directives crossplane.Directives,
	name string,
	fn func(*Directive, int) (bool, error),
) error {
	count := 0
	for _, dir := range directives {
		if dir.IsBlock() {
			continue
		}
		if name != "" && dir.Directive != name {
			continue
		}
		shouldContinue, err := fn(&Directive{Directive: dir}, count)
		if err != nil {
			return apperrors.Wrap(err)
		}
		if !shouldContinue {
			break
		}
		count++
	}
	return nil
}

func convertDirectives(ds crossplane.Directives) []*Directive {
	directives := make([]*Directive, 0, len(ds))
	for _, dir := range ds {
		directives = append(directives, &Directive{Directive: dir})
	}
	return directives
}
