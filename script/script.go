package script

import (
	"errors"

	"gorm.io/gorm"
)

// Script runs a specified script by name using the given gorm.DB connection, returning an error if the script is not found.
func Script(scriptName string, db *gorm.DB) error {
	switch scriptName {
	case "example_script":
		exampleScript := NewExampleScript(db)
		return exampleScript.Run()
	default:
		return errors.New("script not found")
	}
}
