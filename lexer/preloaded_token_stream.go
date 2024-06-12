package lexer

import "github.com/antlr4-go/antlr/v4"

type PreloadedTokenStream struct {
	tokens      []antlr.Token
	index       int
	tokenSource antlr.TokenSource
}

func NewPreloadedTokenStream(tokens []antlr.Token, tokenSource antlr.TokenSource) *PreloadedTokenStream {
	return &PreloadedTokenStream{
		tokens:      tokens,
		index:       0,
		tokenSource: tokenSource,
	}
}

func (p *PreloadedTokenStream) Consume() {
	if p.index >= len(p.tokens) {
		panic("cannot consume EOF")
	}
	p.index++
}

func (p *PreloadedTokenStream) LA(i int) int {
	if i == 0 {
		return 0 // undefined
	}
	if i < 0 {
		i++
		if (p.index + i - 1) < 0 {
			return antlr.TokenEOF
		}
	}
	if (p.index + i - 1) >= len(p.tokens) {
		return antlr.TokenEOF
	}
	return p.tokens[p.index+i-1].GetTokenType()
}

func (p *PreloadedTokenStream) Mark() int {
	return -1
}

func (p *PreloadedTokenStream) Release(marker int) {}

func (p *PreloadedTokenStream) Index() int {
	return p.index
}

func (p *PreloadedTokenStream) Seek(index int) {
	p.index = index
}

func (p *PreloadedTokenStream) Size() int {
	return len(p.tokens)
}

func (p *PreloadedTokenStream) GetSourceName() string {
	return p.tokenSource.GetSourceName()
}

func (p *PreloadedTokenStream) LT(k int) antlr.Token {
	if k == 0 {
		return nil
	}
	if k < 0 {
		k++
		if (p.index + k - 1) < 0 {
			return nil
		}
	}
	if (p.index + k - 1) >= len(p.tokens) {
		return nil
	}
	return p.tokens[p.index+k-1]
}

func (p *PreloadedTokenStream) Reset() {
	p.index = 0
}

func (p *PreloadedTokenStream) Get(index int) antlr.Token {
	return p.tokens[index]
}

func (p *PreloadedTokenStream) GetTokenSource() antlr.TokenSource {
	return p.tokenSource
}

func (p *PreloadedTokenStream) SetTokenSource(tokenSource antlr.TokenSource) {
	p.tokenSource = tokenSource
}

func (p *PreloadedTokenStream) GetAllText() string {
	return p.GetTextFromInterval(antlr.NewInterval(0, len(p.tokens)-1))
}

func (p *PreloadedTokenStream) GetTextFromInterval(interval antlr.Interval) string {
	start := interval.Start
	stop := interval.Stop

	if start < 0 || stop < 0 {
		return ""
	}

	if stop >= len(p.tokens) {
		stop = len(p.tokens) - 1
	}

	s := ""
	for i := start; i <= stop; i++ {
		t := p.tokens[i]
		if t.GetTokenType() == antlr.TokenEOF {
			break
		}
		s += t.GetText()
	}
	return s
}

func (p *PreloadedTokenStream) GetTextFromRuleContext(ctx antlr.RuleContext) string {
	return p.GetTextFromInterval(ctx.GetSourceInterval())
}

func (p *PreloadedTokenStream) GetTextFromTokens(start, end antlr.Token) string {
	if start == nil || end == nil {
		return ""
	}
	return p.GetTextFromInterval(antlr.NewInterval(start.GetTokenIndex(), end.GetTokenIndex()))
}
