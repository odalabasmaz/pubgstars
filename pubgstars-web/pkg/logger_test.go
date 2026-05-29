package pkg

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(fn func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestLog(t *testing.T) {
	out := captureStdout(func() { Log("hello") })
	want := fmt.Sprintf("Log:  %s\n", "hello")
	if !strings.Contains(out, "hello") {
		t.Errorf("Log() output %q does not contain %q", out, want)
	}
}

func TestPrint(t *testing.T) {
	out := captureStdout(func() { Print("world") })
	if !strings.Contains(out, "world") {
		t.Errorf("Print() output %q does not contain %q", out, "world")
	}
}
