package compiler

import (
	"fmt"
	"strconv"
)

type Parser struct {
	tokens       <-chan *Token
	currentToken *Token
	nextToken    *Token
}

func NewParser(tokens <-chan *Token) *Parser {
	p := &Parser{tokens, nil, nil}
	p.advance()
	return p
}

func (p *Parser) advance() {
	p.currentToken = p.nextToken
	p.nextToken = <-p.tokens
}

func (p *Parser) accept(tokType string) bool {
	if p.nextToken != nil && p.nextToken.TokType == tokType {
		p.advance()
		return true
	}

	return false
}

func (p *Parser) expect(tokType string) error {
	if !p.accept(tokType) {
		return fmt.Errorf("expected %s, got token type %s instead", tokType, p.currentToken.TokType)
	}
	return nil
}

func (p *Parser) ParseProgram() (*Program, error) {
	blocks := make(map[int]*NodeBlock)

	for !p.accept(EOF) {
		blk, err := p.parseNodeBlock()
		if err != nil {
			return nil, err
		}

		blocks[blk.ID] = blk
	}

	return &Program{blocks}, nil
}

func (p *Parser) parseNodeBlock() (*NodeBlock, error) {
	var err error

	err = p.expect(AT)

	if err != nil {
		return nil, err
	}

	err = p.expect(DIGIT)
	if err != nil {
		return nil, err
	}

	nodeID, err := strconv.Atoi(p.currentToken.Value)
	if err != nil {
		return nil, err
	}

	err = p.expect(LINEBREAK)
	if err != nil {
		return nil, err
	}

	var statements []*Statement

	for p.nextToken.TokType != AT && p.nextToken.TokType != EOF {
		var stmt *Statement
		stmt, err = p.ParseStatement()
		if err != nil {
			return nil, err
		}

		statements = append(statements, stmt)
	}

	return &NodeBlock{nodeID, statements}, nil
}

func (p *Parser) ParseStatement() (*Statement, error) {
	var label string

	if p.accept(LABEL) {
		label = p.currentToken.Value
	}

	p.accept(LINEBREAK)

	instr, err := p.ParseInstruction()
	if err != nil {
		return nil, err
	}

	return &Statement{
		Label:       label,
		Instruction: instr,
	}, nil
}

func (p *Parser) ParseInstruction() (*Instruction, error) {
	var err error

	err = p.expect(WORD)

	if err != nil {
		return nil, err
	}

	instruction := p.currentToken.Value

	var args []string

	// TODO: Check argument count according to instruction
	for p.accept(WORD) || p.accept(DIGIT) {
		args = append(args, p.currentToken.Value)
	}

	err = p.expect(LINEBREAK)
	if err != nil {
		return nil, err
	}

	return &Instruction{instruction, args}, nil
}
