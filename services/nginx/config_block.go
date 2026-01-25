package nginx

import (
	"strings"

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

func NewLocationBlock(args []string, directives ...*Directive) *Block {
	return NewBlock("location", args, directives...)
}

func (b *Block) AsConfig() *Config {
	return NewConfig(b.AsDirective())
}

func (b *Block) AsDirective() *Directive {
	return &Directive{Directive: b.Directive}
}

func (b *Block) getDirectiveIndex(isBlock bool, name string, args []string, partialMatch bool) int {
	for i, dir := range b.Block {
		if dir.IsBlock() != isBlock || dir.Directive != name {
			continue
		}
		if gofn.Equal(dir.Args, args) { // fully match
			return i
		}
		if partialMatch && len(dir.Args) > len(args) && gofn.Equal(dir.Args[:len(args)], args) { // partially match
			return i
		}
	}
	return -1
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

func (b *Block) AddBlockAt(block *Block, i int) {
	b.Block = gofn.Splice(b.Block, i, 0, block.Directive)
}

func (b *Block) GetBlock(name string, args []string, partialMatch bool) *Block {
	idx := b.getDirectiveIndex(true, name, args, partialMatch)
	if idx >= 0 {
		return &Block{Directive: b.Block[idx]}
	}
	return nil
}

func (b *Block) GetBlockIndex(name string, args []string, partialMatch bool) int {
	return b.getDirectiveIndex(true, name, args, partialMatch)
}

func (b *Block) RemoveBlock(name string, args []string, partialMatch bool) {
	idx := b.GetBlockIndex(name, args, partialMatch)
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

// AddDirectivesInSection inserts directives into the position between 2 comments:
// `# @lp_section section` and `# @lp_section_end section`
func (b *Block) AddDirectivesInSection(section string, directives ...*Directive) {
	idx := b.GetCommentIndex(section, true)
	if idx == -1 {
		return
	}
	if strings.HasPrefix(section, "@lp_section ") {
		idx++
	}
	items := make([]*crossplane.Directive, 0, len(directives))
	for _, dir := range directives {
		items = append(items, dir.Directive)
	}
	b.Block = gofn.Splice(b.Block, idx, 0, items...)
}

func (b *Block) RemoveSectionComments(section string) {
	idx := b.GetCommentIndex(section, true)
	if idx == -1 {
		return
	}
	count := 1
	if idx+1 < len(b.Block) {
		next := b.Block[idx+1]
		if next.Directive == "#" && next.Comment != nil &&
			strings.HasPrefix(strings.TrimSpace(*next.Comment), "@lp_section_end ") {
			count++
		}
	}
	b.Block = gofn.Splice(b.Block, idx, count)
}

func (b *Block) AddDirectiveNamesIfNotExist(names ...string) {
	for _, name := range names {
		if b.GetDirectiveIndex(name, nil, true) >= 0 {
			continue
		}
		b.AddDirectives(NewDirective(name, nil))
	}
}

func (b *Block) AddDirectivesIfNotExist(directives ...*Directive) {
	for _, dir := range directives {
		if b.GetDirectiveIndex(dir.Directive.Directive, dir.Args, false) >= 0 {
			continue
		}
		b.AddDirectives(dir)
	}
}

func (b *Block) SetDirectiveArgs(name string, args []string, n int) {
	for _, directive := range b.DirectivesByName(name, n) {
		directive.Args = args
	}
}

func (b *Block) SetVariable(name string, value string) *Directive {
	dir := b.GetDirective("set", []string{name}, true)
	if dir != nil {
		dir.Args = []string{name, value}
	} else {
		dir = NewDirective("set", []string{name, value})
		b.AddDirectives(dir)
	}
	return dir
}

func (b *Block) GetDirective(name string, args []string, partialMatch bool) *Directive {
	idx := b.GetDirectiveIndex(name, args, partialMatch)
	if idx >= 0 {
		return &Directive{Directive: b.Block[idx]}
	}
	return nil
}

func (b *Block) GetDirectiveIndex(name string, args []string, partialMatch bool) int {
	return b.getDirectiveIndex(false, name, args, partialMatch)
}

func (b *Block) RemoveDirectives(directives ...*Directive) {
	for _, dir := range directives {
		idx := b.GetDirectiveIndex(dir.Directive.Directive, dir.Args, false)
		if idx >= 0 {
			gofn.RemoveAt(&b.Block, idx)
		}
	}
}

func (b *Block) RemoveDirective(name string, args []string, partialMatch bool) {
	idx := b.GetDirectiveIndex(name, args, partialMatch)
	if idx >= 0 {
		gofn.RemoveAt(&b.Block, idx)
	}
}

func (b *Block) SetDirective(name string, args []string, partialMatch bool) {
	idx := b.GetDirectiveIndex(name, args, partialMatch)
	if idx >= 0 {
		gofn.RemoveAt(&b.Block, idx)
	}
}

func (b *Block) GetCommentIndex(comment string, partialMatch bool) int {
	for i, dir := range b.Block {
		if dir.Directive != "#" || dir.Comment == nil {
			continue
		}
		dirComment := strings.TrimSpace(*dir.Comment)
		if !strings.HasPrefix(dirComment, comment) {
			continue
		}
		if !partialMatch && len(dirComment) != len(comment) {
			continue
		}
		return i
	}
	return -1
}

func (b *Block) GetComment(comment string, partialMatch bool) *Directive {
	idx := b.GetCommentIndex(comment, partialMatch)
	if idx >= 0 {
		return &Directive{Directive: b.Block[idx]}
	}
	return nil
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
