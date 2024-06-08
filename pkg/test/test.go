package test

import (
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

// Reads a given file to bytes
func ReadFixtureFile(t *testing.T, inputFile string) []byte {
	raw, err := os.ReadFile(inputFile)

	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	return raw
}

func ReadFixtureFileString(t *testing.T, inputFile string) string {
	raw := ReadFixtureFile(t, inputFile)
	return string(raw)
}

func ReadFixtureFileLines(t *testing.T, inputFile string) []string {
	content := ReadFixtureFileString(t, inputFile)
	return strings.Split(content, "\n")
}

func CaptureLog() *strings.Builder {
	buf := &strings.Builder{}
	logrus.SetOutput(buf)
	return buf
}

func ResetCaptureLog() {
	logrus.SetOutput(os.Stdout)
}
