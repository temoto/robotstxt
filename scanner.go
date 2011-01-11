package robotstxt

import (
    "container/vector"
    "fmt"
    "go/token"
    "os"
    "strings"
    "utf8"
)


type ByteScanner struct {
    ErrorCount int
    Quiet      bool

    buf       []byte
    pos       token.Position
    lastChunk bool
    ch        int
    //
    //state string
}


const WhitespaceChars = " \t\v"


func NewByteScanner(srcname string, quiet bool) *ByteScanner {
    return &ByteScanner{
        Quiet: quiet,
        ch:    -1,
        pos:   token.Position{Filename: srcname},
        //state: "start",
    }
}

func (s *ByteScanner) Feed(input []byte, end bool) (bool, os.Error) {
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

func (s *ByteScanner) Scan() (string, os.Error) {
    //println("--- Scan(). Offset / len(s.buf): ", s.pos.Offset, len(s.buf))

    for {
        // Note Offset > len, not >=, so we can Scan last character.
        if s.lastChunk && s.pos.Offset > len(s.buf) {
            return "", os.EOF
        }

        s.skipSpace()

        if s.ch == -1 {
            return "", os.EOF
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
                return "", os.EOF
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

func (s *ByteScanner) ScanAll() ([]string, os.Error) {
    var results vector.StringVector
    for {
        t, err := s.Scan()
        if t != "" {
            results.Push(t)
        }
        if err == os.EOF {
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
        fmt.Fprintf(os.Stderr, "%s: %s\n", pos.String(), msg)
    }
}

func (s *ByteScanner) isEol() bool {
    return s.ch == '\n' || s.ch == '\r'
}

func (s *ByteScanner) isSpace() bool {
    return strings.Index(WhitespaceChars, string(s.ch)) >= 0
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
func (s *ByteScanner) nextChar() (int, os.Error) {
    //println("--- nextChar(). Offset / len(s.buf): ", s.pos.Offset, len(s.buf))

    if s.pos.Offset >= len(s.buf) {
        s.ch = -1
        return s.ch, os.EOF
    }
    s.pos.Column++
    if s.ch == '\n' {
        s.pos.Line++
        s.pos.Column = 1
    }
    r, w := int(s.buf[s.pos.Offset]), 1
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
