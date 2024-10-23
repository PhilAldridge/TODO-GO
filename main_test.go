package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestTodoCli(t *testing.T) {
	t.Run("Should exit the app when 'exit' command recieved", func(t *testing.T) {
		reader := strings.NewReader("exit\n")
		writer := &bytes.Buffer{}

		todoCli(reader, writer)

		got := writer.String()
		wantSuffix := "Exiting App\n"
		if !strings.HasSuffix(got, wantSuffix) {
			t.Errorf("todoCli() output did not end with %q, got %q", wantSuffix, got)
		}
	})
}
