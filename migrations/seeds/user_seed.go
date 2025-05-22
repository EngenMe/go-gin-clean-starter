package seeds

import (
	"encoding/json"
	"io"
	"os"
	"path"

	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/dto"
	"github.com/Caknoooo/go-gin-clean-starter/entity"
	"github.com/Caknoooo/go-gin-clean-starter/helpers"
)

// ListUserSeeder is a function that seeds user data into the database from a JSON file located in the project directory.
var ListUserSeeder = func(db *gorm.DB) error {
	projectDir, err := helpers.GetProjectRoot()
	if err != nil {
		return err
	}

	jsonFilePath := path.Join(projectDir, "migrations/json/users.json")
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		return err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			panic(err)
		}
	}(jsonFile)

	type SeedUserRequest struct {
		dto.UserCreateRequest
		Role       string `json:"role" binding:"required,oneof=user admin"`
		IsVerified bool   `json:"is_verified"`
	}

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var seedUsers []SeedUserRequest
	if err := json.Unmarshal(jsonData, &seedUsers); err != nil {
		return err
	}

	hasTable := db.Migrator().HasTable(&entity.User{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.User{}); err != nil {
			return err
		}
	}

	for _, seedUser := range seedUsers {
		user := entity.User{
			Name:        seedUser.Name,
			PhoneNumber: seedUser.PhoneNumber,
			Email:       seedUser.Email,
			Password:    seedUser.Password,
			Role:        seedUser.Role,
			IsVerified:  seedUser.IsVerified,
		}

		var existingUser entity.User
		isData := db.Where("email = ?", user.Email).Find(&existingUser).RowsAffected

		if isData == 0 {
			if err := db.Create(&user).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
