package shell

import (
	"bufio"
	"bytes"
	"regexp"
)

type LineBufferWriter struct {
	lines    []string
	capacity int
	index    int
	full     bool
	retained string
	Handler  func(string)
}

func NewLineBufferWriter(n int) *LineBufferWriter {
	return &LineBufferWriter{
		lines:    make([]string, n),
		capacity: n,
		retained: "",
	}
}

var trimRightRegex = regexp.MustCompile("[\r\n]+$")

func trimNewLine(s string) string {
	return trimRightRegex.ReplaceAllString(s, "")

}

func (w *LineBufferWriter) Write(p []byte) (n int, err error) {
	reader := bufio.NewReader(bytes.NewReader(p))
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			w.retained = line
			break
		}
		if len(w.retained) > 0 {
			line = w.retained + line
			w.retained = ""
		}
		line = trimNewLine(line)
		w.onNewLine(line) // callback
		w.lines[w.index] = line
		w.index = (w.index + 1) % w.capacity
		if w.index == 0 {
			w.full = true
		}
	}
	return len(p), nil
}

func (w *LineBufferWriter) Flush() {
	if len(w.retained) > 0 {
		w.Write([]byte("\n"))
	}
}

func (w *LineBufferWriter) onNewLine(line string) {
	if w.Handler != nil {
		w.Handler(line)
	}
}

func (w *LineBufferWriter) GetLines() []string {
	if !w.full {
		return w.lines[:w.index]
	}
	lines := make([]string, w.capacity)
	copy(lines, w.lines[w.index:])
	copy(lines[w.capacity-w.index:], w.lines[:w.index])
	return lines
}
