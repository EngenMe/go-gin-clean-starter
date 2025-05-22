package routes

import (
	"bufio"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// LOG_DIR specifies the directory path where query log files are stored.
	LOG_DIR = "./config/logs/query_log"

	// LOG_HTML represents the filename of the HTML template used to render log data in the application.
	LOG_HTML = "logs.html"
)

// LoggerRoute sets up the routes for logging functionality, including rendering logs and viewing logs by month.
func LoggerRoute(router *gin.Engine) {
	router.LoadHTMLFiles(LOG_HTML)

	router.GET("/logs/:month", Logger)
	router.GET("/logs", Logger)
}

// Logger handles HTTP requests to display log files filtered by month, rendering them using a predefined HTML template.
func Logger(c *gin.Context) {
	month := c.Param("month")
	if month == "" {
		month = time.Now().Format("January")
	}

	logFileName := strings.ToLower(month) + "_query.log"
	logFile := filepath.Join(LOG_DIR, logFileName)

	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		c.HTML(
			http.StatusOK, filepath.Base(LOG_HTML), gin.H{
				"Month": month,
				"Logs":  nil,
			},
		)
		return
	}

	file, err := os.Open(logFile)
	if err != nil {
		c.HTML(
			http.StatusInternalServerError, filepath.Base(LOG_HTML), gin.H{
				"Month": month,
				"Logs":  nil,
			},
		)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	var logs []string
	scanner := bufio.NewScanner(file)
	var block strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			if block.Len() > 0 {
				logs = append(logs, block.String())
				block.Reset()
			}
		} else {
			block.WriteString(line + "\n")
		}
	}

	if block.Len() > 0 {
		logs = append(logs, block.String())
	}

	if err := scanner.Err(); err != nil {
		c.HTML(
			http.StatusInternalServerError, filepath.Base(LOG_HTML), gin.H{
				"Month": month,
				"Logs":  nil,
			},
		)
		return
	}

	reverseLogs := ReverseSlice(logs)

	c.HTML(
		http.StatusOK, filepath.Base(LOG_HTML), gin.H{
			"Month": month,
			"Logs":  reverseLogs,
		},
	)
}

// ReverseSlice reverses the order of elements in a slice of strings and returns the reversed slice.
func ReverseSlice(input []string) []string {
	length := len(input)
	reversed := make([]string, length)
	for i, v := range input {
		reversed[length-1-i] = v
	}
	return reversed
}
