package robotstxt

import (
	"errors"
	"io"
	"strconv"
	"strings"
)

type Parser struct {
	tokens []string
	pos    int
	agent  string
}

func NewParser(tokens []string) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) ParseAll() (result []Rule, err error) {
	var r *Rule
	err = nil
	for {
		r, err = p.ParseRule()
		if r != nil {
			result = append(result, *r)
		}
		if err == io.EOF {
			err = nil
			break
		}
	}
	return result, err
}

func (p *Parser) ParseRule() (r *Rule, err error) {
	t1, ok1 := p.popToken()
	if !ok1 {
		// proper EOF
		return nil, io.EOF
	}

	t2, ok2 := p.peekToken()
	switch strings.ToLower(t1) {
	case "\n":
		// Don't consume t2 and continue parsing
		return nil, nil
	case "user-agent":
		if !ok2 {
			// TODO: report error
			return nil, errors.New("Unexpected EOF at token #" + strconv.Itoa(p.pos) + " namely: \"" + t1 + "\"")
		}
		p.agent = t2
		p.popToken()
		// continue parsing
		return nil, nil
	case "disallow":
		if p.agent == "" {
			// TODO: report error
			return nil, errors.New("Disallow before User-agent.")
		}
		p.popToken()

		if t2 == "" {
			return &Rule{Agent: p.agent, Uri: t2, Allow: true}, nil
		} else {
			return &Rule{Agent: p.agent, Uri: t2, Allow: false}, nil
		}

	case "allow":
		if p.agent == "" {
			// TODO: report error
			return nil, errors.New("Disallow before User-agent.")
		}
		p.popToken()
		return &Rule{Agent: p.agent, Uri: t2, Allow: true}, nil
	}

	return nil, errors.New("Unknown token: " + strconv.Quote(t1))
}

func (p *Parser) popToken() (tok string, ok bool) {
	if p.pos >= len(p.tokens) {
		return "", false
	}
	tok = p.tokens[p.pos]
	p.pos++
	return tok, true
}

func (p *Parser) peekToken() (tok string, ok bool) {
	if p.pos >= len(p.tokens) {
		return "", false
	}
	return p.tokens[p.pos], true
}
