package gts

import (
	"errors"
	"io"

	pars "gopkg.in/pars.v2"
)

// Constructor is the function signature for creating a decoder and object for
// use in DecoderParsers.
type Constructor func(io.Reader) (Decoder, interface{})

// DecoderParser creates a Parser using the Decoder and object constructed
// using the Constructor.
func DecoderParser(ctor Constructor) pars.Parser {
	return func(state *pars.State, result *pars.Result) error {
		state.Push()
		dec, v := ctor(state)
		if err := dec.Decode(v); err != nil {
			state.Pop()
			return err
		}
		result.SetValue(v)
		state.Drop()
		return nil
	}
}

// ParserScanner provides a convenient interface for continually matching a
// parser.
type ParserScanner struct {
	p   pars.Parser
	s   *pars.State
	res pars.Result
	err error
}

func NewParserScanner(r io.Reader, i interface{}) *ParserScanner {
	return &ParserScanner{pars.AsParser(i), pars.NewState(r), pars.Result{}, nil}
}

// Scan advances the scanner using the given parser.
func (s *ParserScanner) Scan() bool {
	if s.err != nil {
		return false
	}
	s.res, s.err = s.p.Parse(s.s)
	return s.err == nil || s.err == io.EOF
}

// Value returns the most recent result value generated by the parser.
func (s ParserScanner) Value() interface{} {
	return s.res.Value
}

// Err returns the first non-EOF error that was encountered by the Scanner.
func (s ParserScanner) Err() error {
	if s.err == nil || s.err == io.EOF {
		return nil
	}
	return s.err
}

// MultiParserScanner tests for a matching parser on the first Scan and will
// attempt to use the same parser for the subsequent calls to Scan.
type MultiParserScanner struct {
	pp  []pars.Parser
	idx int
	s   *pars.State
	res pars.Result
	err error
}

// NewMultiParserScanner creates a new MultiParserScanner.
func NewMultiParserScanner(r io.Reader, ii ...interface{}) *MultiParserScanner {
	pp := pars.AsParsers(ii...)
	return &MultiParserScanner{pp, -1, pars.NewState(r), pars.Result{}, nil}
}

// Scan advances the scanner using the given parsers. The first Scan can match
// any of the given parsers. All subsequent calls to Scan will require the same
// parser to match continually.
func (s *MultiParserScanner) Scan() bool {
	// Parser already returned an error.
	if s.err != nil {
		return false
	}

	// Find the appropriate parser for the first scan.
	if s.idx < 0 {
		s.s.Push()
		for i, p := range s.pp {
			s.res, s.err = p.Parse(s.s)
			if s.err == nil {
				s.s.Drop()
				s.idx = i
				return true
			}
		}
		s.err = errors.New("cannot interpret input bytes as a GenBank style object")
		return false
	}

	s.res, s.err = s.pp[s.idx].Parse(s.s)
	return s.err == nil || s.err == io.EOF
}

// Value returns the most recent result value generated by one of the parsers.
func (s MultiParserScanner) Value() interface{} {
	return s.res.Value
}

// Err returns the first non-EOF error that was encountered by the Scanner.
func (s MultiParserScanner) Err() error {
	if s.err == nil || s.err == io.EOF {
		return nil
	}
	return s.err
}
