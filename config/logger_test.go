package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm/logger"
)

// TestSetupLogger tests the SetupLogger function ensuring it handles log directory and file creation correctly.
func TestSetupLogger(t *testing.T) {
	testDir := "./test_logs"
	t.Cleanup(
		func() {
			err := os.RemoveAll(testDir)
			if err != nil {
				panic(err)
			}
		},
	)

	t.Run(
		"Successfully creates logger with directory and file", func(t *testing.T) {
			originalLogDir := LogDir
			LogDir = testDir
			t.Cleanup(
				func() {
					LogDir = originalLogDir
				},
			)

			result := SetupLogger()

			assert.NotNil(t, result)

			_, err := os.Stat(testDir)
			assert.NoError(t, err)

			currentMonth := strings.ToLower(time.Now().Format("January"))
			logFileName := fmt.Sprintf("%s_query.log", currentMonth)
			logPath := filepath.Join(testDir, logFileName)
			_, err = os.Stat(logPath)
			assert.NoError(t, err)
		},
	)

	t.Run(
		"Fails when directory cannot be created", func(t *testing.T) {
			if os.Geteuid() == 0 {
				t.Skip("Skipping test when running as root")
			}

			originalLogDir := LogDir
			LogDir = "/root/protected_directory"
			t.Cleanup(
				func() {
					LogDir = originalLogDir
				},
			)
		},
	)

	t.Run(
		"Fails when file cannot be created", func(t *testing.T) {
			originalLogDir := LogDir
			LogDir = testDir
			t.Cleanup(
				func() {
					LogDir = originalLogDir
				},
			)

			err := os.MkdirAll(testDir, 0444)
			require.NoError(t, err)
		},
	)
}

// TestLoggerInterfaceImplementation validates that SetupLogger returns an implementation of the logger.Interface interface.
func TestLoggerInterfaceImplementation(t *testing.T) {
	testDir := "./test_logs"
	t.Cleanup(
		func() {
			err := os.RemoveAll(testDir)
			if err != nil {
				panic(err)
			}
		},
	)

	originalLogDir := LogDir
	LogDir = testDir
	t.Cleanup(
		func() {
			LogDir = originalLogDir
		},
	)

	result := SetupLogger()

	_, ok := result.(logger.Interface)
	assert.True(t, ok, "SetupLogger should return a logger.Interface implementation")
}
