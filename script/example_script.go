package script

import (
	"fmt"

	"gorm.io/gorm"
)

// ExampleScript represents a script that operates using a gorm-based database connection.
type (
	ExampleScript struct {
		db *gorm.DB
	}
)

// NewExampleScript initializes and returns a new instance of ExampleScript with the provided gorm.DB connection.
func NewExampleScript(db *gorm.DB) *ExampleScript {
	return &ExampleScript{
		db: db,
	}
}

// Run executes the script logic and interacts with the underlying database, returning an error if execution fails.
func (s *ExampleScript) Run() error {
	fmt.Println("example script running")
	return nil
}
