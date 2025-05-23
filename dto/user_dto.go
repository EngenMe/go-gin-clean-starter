package dto

import (
	"errors"
	"mime/multipart"

	"github.com/Caknoooo/go-gin-clean-starter/entity"
)

const (
	// MESSAGE_FAILED_GET_DATA_FROM_BODY indicates an error occurred while attempting to retrieve data from the request body.
	MESSAGE_FAILED_GET_DATA_FROM_BODY = "failed get data from body"

	// MESSAGE_FAILED_REGISTER_USER represents the failure message for unsuccessful user registration attempts.
	MESSAGE_FAILED_REGISTER_USER = "failed create user"

	// MESSAGE_FAILED_GET_LIST_USER represents an error message for a failure in retrieving the list of users.
	MESSAGE_FAILED_GET_LIST_USER = "failed get list user"

	// MESSAGE_FAILED_TOKEN_NOT_VALID indicates that the provided token is invalid or does not meet the required format.
	MESSAGE_FAILED_TOKEN_NOT_VALID = "token not valid"

	// MESSAGE_FAILED_TOKEN_NOT_FOUND indicates that the required token was not found in the request or header.
	MESSAGE_FAILED_TOKEN_NOT_FOUND = "token not found"

	// MESSAGE_FAILED_GET_USER represents an error message indicating the failure to retrieve user information.
	MESSAGE_FAILED_GET_USER = "failed get user"

	// MESSAGE_FAILED_LOGIN indicates that a login attempt has failed.
	MESSAGE_FAILED_LOGIN = "failed login"

	// MESSAGE_FAILED_UPDATE_USER is a constant string returned when updating a user in the system fails.
	MESSAGE_FAILED_UPDATE_USER = "failed update user"

	// MESSAGE_FAILED_DELETE_USER indicates that the user deletion process has failed.
	MESSAGE_FAILED_DELETE_USER = "failed delete user"

	// MESSAGE_FAILED_PROSES_REQUEST indicates a failure occurred while processing a request.
	MESSAGE_FAILED_PROSES_REQUEST = "failed proses request"

	// MESSAGE_FAILED_DENIED_ACCESS indicates that access to the requested resource has been denied.
	MESSAGE_FAILED_DENIED_ACCESS = "denied access"

	// MESSAGE_FAILED_VERIFY_EMAIL indicates an error occurred during the email verification process.
	MESSAGE_FAILED_VERIFY_EMAIL = "failed verify email"

	// MESSAGE_SUCCESS_REGISTER_USER represents a message indicating the successful registration of a user.
	MESSAGE_SUCCESS_REGISTER_USER = "success create user"

	// MESSAGE_SUCCESS_GET_LIST_USER indicates a successful retrieval of a list of users from the system.
	MESSAGE_SUCCESS_GET_LIST_USER = "success get list user"

	// MESSAGE_SUCCESS_GET_USER represents the success message returned when a user's details are successfully retrieved.
	MESSAGE_SUCCESS_GET_USER = "success get user"

	// MESSAGE_SUCCESS_LOGIN is a constant string indicating a successful login operation.
	MESSAGE_SUCCESS_LOGIN = "success login"

	// MESSAGE_SUCCESS_UPDATE_USER represents a success message for updating user information.
	MESSAGE_SUCCESS_UPDATE_USER = "success update user"

	// MESSAGE_SUCCESS_DELETE_USER indicates a successful user deletion operation.
	MESSAGE_SUCCESS_DELETE_USER = "success delete user"

	// MESSAGE_SEND_VERIFICATION_EMAIL_SUCCESS indicates a successful operation of sending a verification email.
	MESSAGE_SEND_VERIFICATION_EMAIL_SUCCESS = "success send verification email"

	// MESSAGE_SUCCESS_VERIFY_EMAIL represents a success message indicating that email verification has been successfully completed.
	MESSAGE_SUCCESS_VERIFY_EMAIL = "success verify email"
)

var (
	// ErrCreateUser represents an error when the creation of a user record fails.
	ErrCreateUser = errors.New("failed to create user")

	// ErrGetUserById represents an error when failing to retrieve a user by their ID.
	ErrGetUserById = errors.New("failed to get user by id")

	// ErrGetUserByEmail represents an error when failing to retrieve a user by their email address.
	ErrGetUserByEmail = errors.New("failed to get user by email")

	// ErrEmailAlreadyExists indicates that an attempt was made to register with an email that is already in use.
	ErrEmailAlreadyExists = errors.New("email already exist")

	// ErrUpdateUser represents an error that occurs when updating a user record fails.
	ErrUpdateUser = errors.New("failed to update user")

	// ErrUserNotFound indicates that the requested user could not be found in the system.
	ErrUserNotFound = errors.New("user not found")

	// ErrEmailNotFound indicates that the specified email address was not found in the system.
	ErrEmailNotFound = errors.New("email not found")

	// ErrDeleteUser represents an error that occurs when deleting a user record fails.
	ErrDeleteUser = errors.New("failed to delete user")

	// ErrTokenInvalid indicates that the provided token is invalid or cannot be processed.
	ErrTokenInvalid = errors.New("token invalid")

	// ErrTokenExpired indicates that the provided token has expired and is no longer valid.
	ErrTokenExpired = errors.New("token expired")

	// ErrAccountAlreadyVerified indicates that the account has already been marked as verified.
	ErrAccountAlreadyVerified = errors.New("account already verified")
)

type (
	// UserCreateRequest is used to encapsulate data for creating a new user with optional image upload.
	UserCreateRequest struct {
		Name        string                `json:"name" form:"name" binding:"required,min=2,max=100"`
		PhoneNumber string                `json:"phone_number" form:"phone_number" binding:"omitempty,min=8,max=20"`
		Email       string                `json:"email" form:"email" binding:"required,email"`
		Password    string                `json:"password" form:"password" binding:"required,min=8"`
		Image       *multipart.FileHeader `json:"-" form:"image" swaggerignore:"true"`
	}

	// UserResponse represents the structure for user data returned in API responses.
	UserResponse struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Role        string `json:"role"`
		ImageUrl    string `json:"image_url"`
		IsVerified  bool   `json:"is_verified"`
	}

	// UserPaginationResponse represents paginated response data for a list of users including metadata and user details.
	UserPaginationResponse struct {
		Data []UserResponse `json:"data"`
		PaginationResponse
	}

	// GetAllUserRepositoryResponse represents the response structure containing a list of users and pagination metadata.
	GetAllUserRepositoryResponse struct {
		Users []entity.User `json:"users"`
		PaginationResponse
	}

	// UserUpdateRequest represents a request to update user details such as name, telephone number, and email.
	UserUpdateRequest struct {
		Name        string `json:"name" form:"name" binding:"omitempty,min=2,max=100"`
		PhoneNumber string `json:"phone_number" form:"phone_number" binding:"omitempty,min=8,max=20"`
		Email       string `json:"email" form:"email" binding:"omitempty,email"`
	}

	// UserUpdateResponse represents the response type returned after updating user details in the system.
	UserUpdateResponse struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
		Role        string `json:"role"`
		Email       string `json:"email"`
		IsVerified  bool   `json:"is_verified"`
	}

	// SendVerificationEmailRequest represents a request to send a verification email to the user.
	SendVerificationEmailRequest struct {
		Email string `json:"email" form:"email" binding:"required"`
	}

	// VerifyEmailRequest represents the request structure for verifying user email using a token.
	VerifyEmailRequest struct {
		Token string `json:"token" form:"token" binding:"required"`
	}

	// VerifyEmailResponse represents the response structure for email verification containing email and verification status.
	VerifyEmailResponse struct {
		Email      string `json:"email"`
		IsVerified bool   `json:"is_verified"`
	}

	// UserLoginRequest represents the required fields for a user login request.
	// Email is the user's email address, required for authentication.
	// Password is the user's password, also required for authentication.
	UserLoginRequest struct {
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}
)
