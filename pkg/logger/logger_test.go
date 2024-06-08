// This package includes an environment-specific log format setting.
package logger

import (
	"testing"

	logrus "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/williampsena/bugs-channel/pkg/test"
)

func TestSetupDevelopment(t *testing.T) {
	t.Setenv("GO_ENV", "development")

	Setup()

	buf := test.CaptureLog()

	logrus.Error("Logrus testing")

	assert.Contains(t, buf.String(), "level=error msg=\"Logrus testing\"")

	test.ResetCaptureLog()
}

func TestSetupProduction(t *testing.T) {
	t.Setenv("GO_ENV", "production")

	Setup()

	buf := test.CaptureLog()

	logrus.Info("Logrus testing")

	assert.Contains(t, buf.String(), "\"level\":\"info\",\"msg\":\"Logrus testing\"")

	test.ResetCaptureLog()
}
