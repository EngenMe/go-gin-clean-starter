package container

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// BaseSuite is a wrapper around suite.Suite providing utility methods for test environment setup and cleanup.
type BaseSuite struct {
	suite.Suite
}

// SetupSuite initializes the suite environment by loading environment variables from the .env.test file.
func (s *BaseSuite) SetupSuite() {
	LoadTestEnv()
}

// SetupEnv sets multiple environment variables from the provided map and ensures no errors during the process.
func (s *BaseSuite) SetupEnv(vars map[string]string) {
	err := SetEnv(vars)
	require.NoError(s.T(), err)
}

// CleanupEnv removes the specified environment variables by their keys and fails the test if an error occurs.
func (s *BaseSuite) CleanupEnv(keys []string) {
	err := UnsetEnv(keys)
	require.NoError(s.T(), err)
}
