package robotstxt

import (
	"strconv"
	"testing"
)

func TestScan001(t *testing.T) {
	sc := newByteScanner("test-001", false)
	if _, err := sc.Scan(); err == nil {
		t.Fatal("Empty ByteScanner should fail on Scan.")
	}
}

func TestScan002(t *testing.T) {
	sc := newByteScanner("test-002", false)
	if err := sc.Feed([]byte("foo"), true); err != nil {
		t.Fatal(err)
	}
	tok, err := sc.Scan()
	t.Logf("Scan tok=%v err=%v", tok, err)
	if err != nil {
		t.Fatal(err)
	}
}

func TestScan004(t *testing.T) {
	sc := newByteScanner("test-004", false)
	if err := sc.Feed([]byte("\u2010"), true); err != nil {
		t.Fatal(err)
	}
	tok, err := sc.Scan()
	t.Logf("Scan tok=%v err=%v", tok, err)
	if err != nil {
		t.Fatal(err)
	}
}

func TestScan005(t *testing.T) {
	sc := newByteScanner("test-005", true)
	if err := sc.Feed([]byte("\xd9\xd9"), true); err != nil {
		t.Fatal(err)
	}
	tok, err := sc.Scan()
	t.Logf("Scan tok=%v err=%v", tok, err)
	if err != nil {
		t.Fatal(err)
	}
	if sc.ErrorCount != 2 {
		t.Fatal("Expecting ErrorCount be exactly 2.")
	}
}

func TestScan006(t *testing.T) {
	sc := newByteScanner("test-006", false)
	s := "# comment \r\nSomething: Somewhere\r\n"
	if err := sc.Feed([]byte(s), true); err != nil {
		t.Fatal(err)
	}
	tokens, err := sc.ScanAll()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("len(tokens)=%d", len(tokens))
	if len(tokens) != 4 {
		t.Fatal("Expecting exactly 4 tokens.")
	}
	if tokens[0] != "\n" || tokens[1] != "Something" || tokens[2] != "Somewhere" || tokens[3] != "\n" {
		t.Fatal("Wrong tokens read:", strconv.Quote(tokens[0]), strconv.Quote(tokens[1]), strconv.Quote(tokens[2]), strconv.Quote(tokens[3]))
	}
}

func TestScan007(t *testing.T) {
	sc := newByteScanner("test-007", false)
	s := "# comment \r\n# more comments\n\nDisallow:\r"
	if err := sc.Feed([]byte(s), true); err != nil {
		t.Fatal(err)
	}
	tokens, err := sc.ScanAll()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("len(tokens)=%d", len(tokens))
	if len(tokens) != 4 {
		t.Fatal("Expecting exactly 4 tokens.")
	}
	if tokens[0] != "\n" || tokens[1] != "\n" || tokens[2] != "Disallow" || tokens[3] != "\n" {
		t.Fatal("Wrong tokens read:", strconv.Quote(tokens[0]), strconv.Quote(tokens[1]), strconv.Quote(tokens[2]), strconv.Quote(tokens[3]))
	}
}

func TestScanUnicode8BOM(t *testing.T) {
	sc := newByteScanner("test-bom", false)
	if err := sc.Feed([]byte(robotsTextVanityfair), true); err != nil {
		t.Fatal(err)
	}
	tokens, err := sc.ScanAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) == 0 {
		t.Fatal("Read zero tokens.")
	}
	if tokens[0] != "User-agent" {
		t.Fatal("Expecting first token: User-agent")
	}
}
