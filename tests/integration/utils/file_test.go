package utils_test

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/Caknoooo/go-gin-clean-starter/utils"
)

// FileUploadIntegrationTestSuite is a test suite for integration testing of file upload functionalities.
type FileUploadIntegrationTestSuite struct {
	suite.Suite
	testDir string
}

// SetupSuite prepares the test suite environment by creating the necessary test directory and initializing paths.
func (suite *FileUploadIntegrationTestSuite) SetupSuite() {
	suite.testDir = "./test_assets_integration"
	utils.PATH = suite.testDir
	err := os.MkdirAll(suite.testDir, 0755)
	if err != nil {
		panic(err)
	}
}

// TearDownSuite cleans up resources after all tests in the suite have run, including deleting the test directory.
func (suite *FileUploadIntegrationTestSuite) TearDownSuite() {
	err := os.RemoveAll(suite.testDir)
	if err != nil {
		panic(err)
	}
	utils.PATH = "assets"
}

// createTestFile creates a new test file with the given content and returns its multipart file header.
func (suite *FileUploadIntegrationTestSuite) createTestFile(content string) *multipart.FileHeader {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	require.NoError(suite.T(), err)
	_, err = part.Write([]byte(content))
	require.NoError(suite.T(), err)
	err = writer.Close()
	if err != nil {
		return nil
	}

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	file, header, err := req.FormFile("file")
	require.NoError(suite.T(), err)
	err = file.Close()
	if err != nil {
		return nil
	}

	return header
}

// TestUploadFile_Integration verifies the successful uploading of files and handling of nested directories using integration tests.
func (suite *FileUploadIntegrationTestSuite) TestUploadFile_Integration() {
	tests := []struct {
		name        string
		path        string
		fileContent string
		wantErr     bool
	}{
		{
			name:        "Successfully upload file",
			path:        "images/test123",
			fileContent: "integration test content",
			wantErr:     false,
		},
		{
			name:        "Create nested directory structure",
			path:        "images/nested/test123",
			fileContent: "nested content",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		suite.Run(
			tt.name, func() {
				fileHeader := suite.createTestFile(tt.fileContent)
				err := utils.UploadFile(fileHeader, tt.path)

				if tt.wantErr {
					assert.Error(suite.T(), err)
				} else {
					assert.NoError(suite.T(), err)

					parts := strings.Split(tt.path, "/")
					fileID := parts[len(parts)-1]
					dirPath := filepath.Join(suite.testDir, strings.Join(parts[:len(parts)-1], "/"))
					filePath := filepath.Join(dirPath, fileID)

					_, err = os.Stat(filePath)
					assert.NoError(suite.T(), err)

					content, err := os.ReadFile(filePath)
					require.NoError(suite.T(), err)
					assert.Equal(suite.T(), tt.fileContent, string(content))
				}
			},
		)
	}
}

// TestFileUploadIntegrationTestSuite runs the integration test suite for file upload functionality using testify's suite package.
func TestFileUploadIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(FileUploadIntegrationTestSuite))
}
