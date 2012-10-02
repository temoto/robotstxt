package robotstxt

import (
    "fmt"
    "go/token"
    "io"
    "os"
    "unicode/utf8"
)

type ByteScanner struct {
    ErrorCount int
    Quiet      bool

    buf       []byte
    pos       token.Position
    lastChunk bool
    ch        rune
    //
    //state string
}

var WhitespaceChars = []rune{' ', '\t', '\v'}

func NewByteScanner(srcname string, quiet bool) *ByteScanner {
    return &ByteScanner{
        Quiet: quiet,
        ch:    -1,
        pos:   token.Position{Filename: srcname},
        //state: "start",
    }
}

func (s *ByteScanner) Feed(input []byte, end bool) (bool, error) {
    s.buf = input
    s.pos.Offset = 0
    s.pos.Line = 1
    s.pos.Column = 1
    s.lastChunk = end
    // Read first char into look-ahead buffer `s.ch`.
    s.nextChar()
    return false, nil
}

func (s *ByteScanner) GetPosition() token.Position {
    return s.pos
}

func (s *ByteScanner) Scan() (string, error) {
    //println("--- Scan(). Offset / len(s.buf): ", s.pos.Offset, len(s.buf))

    bufsize := len(s.buf)
    for {
        // Note Offset > len, not >=, so we can Scan last character.
        if s.lastChunk && s.pos.Offset > bufsize {
            return "", io.EOF
        }

        s.skipSpace()

        if s.ch == -1 {
            return "", io.EOF
        }

        // EOL
        if s.isEol() {
            // skip subsequent newline chars
            for s.ch != -1 && s.isEol() {
                s.nextChar()
            }
            // emit newline as separate token
            return "\n", nil
        }

        // skip comments
        if s.ch == '#' {
            s.skipUntilEol()
            //            s.state = "start"
            if s.ch == -1 {
                return "", io.EOF
            }
            // emit newline as separate token
            return "\n", nil
        }

        // else we found something
        break
    }

    /*
       if s.state == "start" {
           s.state = "key"
       }
    */

    tok := string(s.ch)
    s.nextChar()
    for s.ch != -1 && !s.isSpace() && !s.isEol() {
        if s.ch == ':' {
            //            s.state = "pre-value"
            s.nextChar()
            break
        }

        tok += string(s.ch)
        s.nextChar()
    }
    return tok, nil
}

func (s *ByteScanner) ScanAll() ([]string, error) {
    var results []string
    for {
        t, err := s.Scan()
        if t != "" {
            results = append(results, t)
        }
        if err == io.EOF {
            break
        }
        if err != nil {
            return results, err
        }
    }
    return results, nil
}

func (s *ByteScanner) error(pos token.Position, msg string) {
    s.ErrorCount++
    if !s.Quiet {
        fmt.Fprintf(os.Stderr, "robotstxt from %s: %s\n", pos.String(), msg)
    }
}

func (s *ByteScanner) isEol() bool {
    return s.ch == '\n' || s.ch == '\r'
}

func (s *ByteScanner) isSpace() bool {
    for i, _ := range WhitespaceChars {
        if s.ch == WhitespaceChars[i] {
            return true
        }
    }
    return false
}

func (s *ByteScanner) skipSpace() {
    //println("--- string(ch): ", s.ch, ".")
    for s.ch != -1 && s.isSpace() {
        s.nextChar()
    }
}

func (s *ByteScanner) skipUntilEol() {
    //println("--- string(ch): ", s.ch, ".")
    for s.ch != -1 && !s.isEol() {
        s.nextChar()
    }
    // skip subsequent newline chars
    for s.ch != -1 && s.isEol() {
        s.nextChar()
    }
}

// Reads next Unicode char.
func (s *ByteScanner) nextChar() (rune, error) {
    //println("--- nextChar(). Offset / len(s.buf): ", s.pos.Offset, len(s.buf))

    if s.pos.Offset >= len(s.buf) {
        s.ch = -1
        return s.ch, io.EOF
    }
    s.pos.Column++
    if s.ch == '\n' {
        s.pos.Line++
        s.pos.Column = 1
    }
    r, w := rune(s.buf[s.pos.Offset]), 1
    if r >= 0x80 {
        r, w = utf8.DecodeRune(s.buf[s.pos.Offset:])
        if r == utf8.RuneError && w == 1 {
            s.error(s.pos, "illegal UTF-8 encoding")
        }
    }
    s.pos.Offset += w
    s.ch = r
    return s.ch, nil
}
