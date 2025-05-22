package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/constants"
	"github.com/Caknoooo/go-gin-clean-starter/dto"
	"github.com/Caknoooo/go-gin-clean-starter/entity"
	"github.com/Caknoooo/go-gin-clean-starter/helpers"
	"github.com/Caknoooo/go-gin-clean-starter/repository"
	"github.com/Caknoooo/go-gin-clean-starter/utils"
)

type (
	// UserService defines the methods required for managing user operations and authentication.
	// Register creates a new user with the provided creation request.
	// GetAllUserWithPagination retrieves paginated lists of users based on the request criteria.
	// GetUserById fetches user details using a specific user ID.
	// GetUserByEmail retrieves user details using their email address.
	// SendVerificationEmail sends a verification email to the specified user.
	// VerifyEmail processes email verification requests and returns a response.
	// Update modifies user information based on the update request and user ID.
	// Delete removes a user from the system using their user ID.
	// Verify authenticates a user and generates a token based on login request details.
	// RefreshToken generates a new access token using a valid refresh token.
	// RevokeRefreshToken revokes all refresh tokens for a specified user ID.
	UserService interface {
		Register(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error)
		GetAllUserWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.UserPaginationResponse, error)
		GetUserById(ctx context.Context, userId string) (dto.UserResponse, error)
		GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, error)
		SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error
		VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error)
		Update(ctx context.Context, req dto.UserUpdateRequest, userId string) (dto.UserUpdateResponse, error)
		Delete(ctx context.Context, userId string) error
		Verify(ctx context.Context, req dto.UserLoginRequest) (dto.TokenResponse, error)
		RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (dto.TokenResponse, error)
		RevokeRefreshToken(ctx context.Context, userID string) error
	}

	// userService is a struct that implements the UserService interface and manages user-related operations.
	userService struct {
		userRepo         repository.UserRepository
		refreshTokenRepo repository.RefreshTokenRepository
		jwtService       JWTService
		db               *gorm.DB
	}
)

// NewUserService initializes and returns a new instance of UserService with the provided dependencies.
func NewUserService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtService JWTService,
	db *gorm.DB,
) UserService {
	return &userService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
		db:               db,
	}
}

const (
	// LOCAL_URL defines the base URL used for local development, pointing to the localhost server running on port 3000.
	LOCAL_URL = "http://localhost:3000"

	// VERIFY_EMAIL_ROUTE specifies the endpoint for email verification during the user registration process.
	VERIFY_EMAIL_ROUTE = "register/verify_email"
)

// SafeRollback ensures that a transaction is safely rolled back in case of a panic and re-panics after rollback.
func SafeRollback(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
		panic(r)
	}
}

// Register handles user registration by creating a new user, verifying email uniqueness, and sending a verification email.
func (s *userService) Register(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error) {
	var filename string

	_, flag, err := s.userRepo.CheckEmail(ctx, nil, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.UserResponse{}, err
	}

	if flag {
		return dto.UserResponse{}, dto.ErrEmailAlreadyExists
	}

	if req.Image != nil {
		imageId := uuid.New()
		ext := utils.GetExtensions(req.Image.Filename)

		filename = fmt.Sprintf("profile/%s.%s", imageId, ext)
		if err := utils.UploadFile(req.Image, filename); err != nil {
			return dto.UserResponse{}, err
		}
	}

	user := entity.User{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		ImageUrl:    filename,
		Role:        constants.ENUM_ROLE_USER,
		Email:       req.Email,
		Password:    req.Password,
		IsVerified:  false,
	}

	userReg, err := s.userRepo.Register(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, dto.ErrCreateUser
	}

	draftEmail, err := makeVerificationEmail(userReg.Email)
	if err != nil {
		return dto.UserResponse{}, err
	}

	err = utils.SendMail(userReg.Email, draftEmail["subject"], draftEmail["body"])
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:          userReg.ID.String(),
		Name:        userReg.Name,
		PhoneNumber: userReg.PhoneNumber,
		ImageUrl:    userReg.ImageUrl,
		Role:        userReg.Role,
		Email:       userReg.Email,
		IsVerified:  userReg.IsVerified,
	}, nil
}

// makeVerificationEmail generates a verification email draft with a secure token embedded in the verification link.
// It returns a map containing the subject and body of the email or an error if the operation fails.
func makeVerificationEmail(receiverEmail string) (map[string]string, error) {
	expired := time.Now().Add(time.Hour * 24).Format("2006-01-02 15:04:05")
	plainText := receiverEmail + "_" + expired
	token, err := utils.AESEncrypt(plainText)
	if err != nil {
		return nil, err
	}

	verifyLink := LOCAL_URL + "/" + VERIFY_EMAIL_ROUTE + "?token=" + token

	templatePath := "utils/email-template/base_mail.html"
	projectPath, err := helpers.GetProjectRoot()
	if err != nil {
		return nil, err
	}
	templatePath = path.Join(projectPath, templatePath)
	readHtml, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}

	data := struct {
		Email  string
		Verify string
	}{
		Email:  receiverEmail,
		Verify: verifyLink,
	}

	tmpl, err := template.New("custom").Parse(string(readHtml))
	if err != nil {
		return nil, err
	}

	var strMail bytes.Buffer
	if err := tmpl.Execute(&strMail, data); err != nil {
		return nil, err
	}

	draftEmail := map[string]string{
		"subject": "Cakno - Go Gin Template",
		"body":    strMail.String(),
	}

	return draftEmail, nil
}

// SendVerificationEmail sends a verification email to the user specified in the request.
// It retrieves the user by email, creates a verification email draft, and sends the email.
func (s *userService) SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error {
	user, err := s.userRepo.GetUserByEmail(ctx, nil, req.Email)
	if err != nil {
		return dto.ErrEmailNotFound
	}

	draftEmail, err := makeVerificationEmail(user.Email)
	if err != nil {
		return err
	}

	err = utils.SendMail(user.Email, draftEmail["subject"], draftEmail["body"])
	if err != nil {
		return err
	}

	return nil
}

// VerifyEmail verifies a user's email using a token, updating the user's verified status if the token is valid and unexpired.
func (s *userService) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error) {
	decryptedToken, err := utils.AESDecrypt(req.Token)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
	}

	if !strings.Contains(decryptedToken, "_") {
		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
	}

	decryptedTokenSplit := strings.Split(decryptedToken, "_")
	email := decryptedTokenSplit[0]
	expired := decryptedTokenSplit[1]

	now := time.Now()
	expiredTime, err := time.Parse("2006-01-02 15:04:05", expired)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
	}

	if expiredTime.Sub(now) < 0 {
		return dto.VerifyEmailResponse{
			Email:      email,
			IsVerified: false,
		}, dto.ErrTokenExpired
	}

	user, err := s.userRepo.GetUserByEmail(ctx, nil, email)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrUserNotFound
	}

	if user.IsVerified {
		return dto.VerifyEmailResponse{}, dto.ErrAccountAlreadyVerified
	}

	updatedUser, err := s.userRepo.Update(
		ctx, nil, entity.User{
			ID:         user.ID,
			IsVerified: true,
		},
	)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrUpdateUser
	}

	return dto.VerifyEmailResponse{
		Email:      email,
		IsVerified: updatedUser.IsVerified,
	}, nil
}

// GetAllUserWithPagination retrieves paginated user data based on the provided pagination request and context.
func (s *userService) GetAllUserWithPagination(
	ctx context.Context,
	req dto.PaginationRequest,
) (dto.UserPaginationResponse, error) {
	dataWithPaginate, err := s.userRepo.GetAllUserWithPagination(ctx, nil, req)
	if err != nil {
		return dto.UserPaginationResponse{}, err
	}

	var users []dto.UserResponse
	for _, user := range dataWithPaginate.Users {
		data := dto.UserResponse{
			ID:          user.ID.String(),
			Name:        user.Name,
			Email:       user.Email,
			Role:        user.Role,
			PhoneNumber: user.PhoneNumber,
			ImageUrl:    user.ImageUrl,
			IsVerified:  user.IsVerified,
		}

		users = append(users, data)
	}

	return dto.UserPaginationResponse{
		Data: users,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

// GetUserById retrieves a user by their ID and returns a UserResponse or an error if the operation fails.
func (s *userService) GetUserById(ctx context.Context, userId string) (dto.UserResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, nil, userId)
	if err != nil {
		return dto.UserResponse{}, dto.ErrGetUserById
	}

	return dto.UserResponse{
		ID:          user.ID.String(),
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		Email:       user.Email,
		ImageUrl:    user.ImageUrl,
		IsVerified:  user.IsVerified,
	}, nil
}

// GetUserByEmail retrieves a user by their email address and returns a UserResponse or an error if the operation fails.
func (s *userService) GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, error) {
	emails, err := s.userRepo.GetUserByEmail(ctx, nil, email)
	if err != nil {
		return dto.UserResponse{}, dto.ErrGetUserByEmail
	}

	return dto.UserResponse{
		ID:          emails.ID.String(),
		Name:        emails.Name,
		PhoneNumber: emails.PhoneNumber,
		Role:        emails.Role,
		Email:       emails.Email,
		ImageUrl:    emails.ImageUrl,
		IsVerified:  emails.IsVerified,
	}, nil
}

// Update modifies user details based on the provided request data and user ID, returning the updated user or an error.
func (s *userService) Update(ctx context.Context, req dto.UserUpdateRequest, userId string) (
	dto.UserUpdateResponse,
	error,
) {
	user, err := s.userRepo.GetUserById(ctx, nil, userId)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUserNotFound
	}

	data := entity.User{
		ID:          user.ID,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Role:        user.Role,
		Email:       req.Email,
	}

	userUpdate, err := s.userRepo.Update(ctx, nil, data)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUpdateUser
	}

	return dto.UserUpdateResponse{
		ID:          userUpdate.ID.String(),
		Name:        userUpdate.Name,
		PhoneNumber: userUpdate.PhoneNumber,
		Role:        userUpdate.Role,
		Email:       userUpdate.Email,
		IsVerified:  user.IsVerified,
	}, nil
}

// Delete removes a user and their associated refresh tokens from the database, ensuring proper error handling and rollback.
func (s *userService) Delete(ctx context.Context, userId string) error {
	tx := s.db.Begin()
	defer SafeRollback(tx)

	user, err := s.userRepo.GetUserById(ctx, tx, userId)
	if err != nil {
		tx.Rollback()
		return dto.ErrUserNotFound
	}

	if err := s.refreshTokenRepo.DeleteByUserID(ctx, tx, user.ID.String()); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete refresh tokens: %w", err)
	}

	err = s.userRepo.Delete(ctx, tx, user.ID.String())
	if err != nil {
		tx.Rollback()
		return dto.ErrDeleteUser
	}

	return tx.Commit().Error
}

// Verify authenticates a user by validating their credentials, generating tokens, and saving the refresh token to the database.
func (s *userService) Verify(ctx context.Context, req dto.UserLoginRequest) (dto.TokenResponse, error) {
	tx := s.db.Begin()
	defer SafeRollback(tx)

	user, err := s.userRepo.GetUserByEmail(ctx, tx, req.Email)
	if err != nil {
		tx.Rollback()
		return dto.TokenResponse{}, errors.New("invalid email or password")
	}

	checkPassword, err := helpers.CheckPassword(user.Password, []byte(req.Password))
	if err != nil || !checkPassword {
		tx.Rollback()
		return dto.TokenResponse{}, errors.New("invalid email or password")
	}

	accessToken := s.jwtService.GenerateAccessToken(user.ID.String(), user.Role)

	refreshTokenString, expiresAt := s.jwtService.GenerateRefreshToken()

	if err := s.refreshTokenRepo.DeleteByUserID(ctx, tx, user.ID.String()); err != nil {
		tx.Rollback()
		return dto.TokenResponse{}, err
	}

	refreshToken := entity.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: expiresAt,
	}

	if _, err := s.refreshTokenRepo.Create(ctx, tx, refreshToken); err != nil {
		tx.Rollback()
		return dto.TokenResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		Role:         user.Role,
	}, nil
}

// RefreshToken refreshes the user's access and refresh tokens using a valid refresh token, ensuring proper validation and handling.
func (s *userService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (dto.TokenResponse, error) {
	tx := s.db.Begin()
	defer SafeRollback(tx)

	dbToken, err := s.refreshTokenRepo.FindByToken(ctx, tx, req.RefreshToken)
	if err != nil {
		tx.Rollback()
		return dto.TokenResponse{}, errors.New(dto.MESSAGE_FAILED_INVALID_REFRESH_TOKEN)
	}

	if time.Now().After(dbToken.ExpiresAt) {
		tx.Rollback()
		return dto.TokenResponse{}, errors.New(dto.MESSAGE_FAILED_EXPIRED_REFRESH_TOKEN)
	}

	user, err := s.userRepo.GetUserById(ctx, tx, dbToken.UserID.String())
	if err != nil {
		tx.Rollback()
		return dto.TokenResponse{}, dto.ErrUserNotFound
	}

	accessToken := s.jwtService.GenerateAccessToken(user.ID.String(), user.Role)

	refreshTokenString, expiresAt := s.jwtService.GenerateRefreshToken()

	if err := s.refreshTokenRepo.DeleteByUserID(ctx, tx, user.ID.String()); err != nil {
		tx.Rollback()
		return dto.TokenResponse{}, err
	}

	refreshToken := entity.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: expiresAt,
	}

	if _, err := s.refreshTokenRepo.Create(ctx, tx, refreshToken); err != nil {
		tx.Rollback()
		return dto.TokenResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		Role:         user.Role,
	}, nil
}

// RevokeRefreshToken revokes all refresh tokens associated with a user by their ID, ensuring proper transaction handling.
func (s *userService) RevokeRefreshToken(ctx context.Context, userID string) error {
	tx := s.db.Begin()
	defer SafeRollback(tx)

	_, err := s.userRepo.GetUserById(ctx, tx, userID)
	if err != nil {
		tx.Rollback()
		return dto.ErrUserNotFound
	}

	if err := s.refreshTokenRepo.DeleteByUserID(ctx, tx, userID); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
