package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/config"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/domain"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userRepo     UserRepo
	categoryRepo CategoryRepo
	walletRepo   WalletRepo
	log          *zap.Logger
}

func NewService(userRepo UserRepo, categoryRepo CategoryRepo, walletRepo WalletRepo, log *zap.Logger) Service {
	return &service{
		userRepo:     userRepo,
		categoryRepo: categoryRepo,
		walletRepo:   walletRepo,
		log:          log,
	}
}

func (s *service) Register(ctx context.Context, input RegisterInput) (*AuthResult, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	username := strings.TrimSpace(strings.ToLower(input.Username))
	name := strings.TrimSpace(input.Name)

	if email == "" || username == "" || name == "" || len(input.Password) < 6 {
		return nil, domain.ErrInvalidInput
	}

	existingEmail, _ := s.userRepo.FindByEmail(ctx, email)
	if existingEmail != nil {
		return nil, domain.ErrEmailAlreadyInUse
	}

	existingUsername, _ := s.userRepo.FindByUsername(ctx, username)
	if existingUsername != nil {
		return nil, domain.ErrUsernameAlreadyInUse
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to hash password", zap.Error(err))
		return nil, domain.ErrInternalServer
	}

	hashStr := string(hashedPassword)
	now := time.Now()
	u := &domain.User{
		ID:           uuid.New().String(),
		Email:        email,
		Username:     username,
		Password:     &hashStr,
		Name:         name,
		AuthProvider: domain.AuthProviderLocal,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		s.log.Error("failed to create user", zap.Error(err))
		return nil, err
	}

	// Create default categories and wallets for new user
	if s.categoryRepo != nil {
		if err := s.categoryRepo.CreateDefaultCategories(ctx, u.ID); err != nil {
			s.log.Warn("failed to create default categories for new user", zap.String("user_id", u.ID), zap.Error(err))
		}
	}
	if s.walletRepo != nil {
		if err := s.walletRepo.CreateDefaultWallets(ctx, u.ID); err != nil {
			s.log.Warn("failed to create default wallets for new user", zap.String("user_id", u.ID), zap.Error(err))
		}
	}

	token, err := s.generateToken(u.ID, u.Email, u.Username, u.Name, "access", config.GetConfig().JWT.AccessTokenDuration)
	if err != nil {
		s.log.Error("failed to generate access token", zap.Error(err))
		return nil, domain.ErrInternalServer
	}

	refreshToken, err := s.generateToken(u.ID, u.Email, u.Username, u.Name, "refresh", config.GetConfig().JWT.RefreshTokenDuration)
	if err != nil {
		s.log.Error("failed to generate refresh token", zap.Error(err))
		return nil, domain.ErrInternalServer
	}

	return &AuthResult{
		User:         u,
		AccessToken:  token,
		RefreshToken: refreshToken,
		IsNewUser:    true,
	}, nil
}

func (s *service) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || input.Password == "" {
		return nil, domain.ErrInvalidInput
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, domain.ErrInvalidCredentials
	}

	if user.AuthProvider != domain.AuthProviderLocal && (user.Password == nil || *user.Password == "") {
		return nil, fmt.Errorf("this account is registered with %s. Please use %s to log in", user.AuthProvider, user.AuthProvider)
	}

	if user.Password == nil {
		return nil, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(input.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := s.generateToken(user.ID, user.Email, user.Username, user.Name, "access", config.GetConfig().JWT.AccessTokenDuration)
	if err != nil {
		s.log.Error("failed to generate token", zap.Error(err))
		return nil, domain.ErrInternalServer
	}

	refreshToken, err := s.generateToken(user.ID, user.Email, user.Username, user.Name, "refresh", config.GetConfig().JWT.RefreshTokenDuration)
	if err != nil {
		s.log.Error("failed to generate refresh token", zap.Error(err))
		return nil, domain.ErrInternalServer
	}

	return &AuthResult{
		User:         user,
		AccessToken:  token,
		RefreshToken: refreshToken,
		IsNewUser:    false,
	}, nil
}

func (s *service) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, domain.ErrInvalidInput
	}
	return s.userRepo.FindByID(ctx, userID)
}

func (s *service) UpdateProfile(ctx context.Context, userID string, input UpdateProfileInput) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if input.Email != nil {
		newEmail := strings.TrimSpace(strings.ToLower(*input.Email))
		if newEmail != "" && newEmail != user.Email {
			existing, _ := s.userRepo.FindByEmail(ctx, newEmail)
			if existing != nil && existing.ID != userID {
				return nil, domain.ErrEmailAlreadyInUse
			}
			user.Email = newEmail
		}
	}

	if input.Username != nil {
		newUsername := strings.TrimSpace(strings.ToLower(*input.Username))
		if newUsername != "" && newUsername != strings.ToLower(user.Username) {
			existing, _ := s.userRepo.FindByUsername(ctx, newUsername)
			if existing != nil && existing.ID != userID {
				return nil, domain.ErrUsernameAlreadyInUse
			}
			user.Username = newUsername
		}
	}

	if input.Name != nil && strings.TrimSpace(*input.Name) != "" {
		user.Name = strings.TrimSpace(*input.Name)
	}

	if input.Avatar != nil {
		user.Avatar = input.Avatar
	}

	if input.Password != nil && len(*input.Password) >= 6 {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, domain.ErrInternalServer
		}
		hashStr := string(hashed)
		user.Password = &hashStr
	}

	user.UpdatedAt = time.Now()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) OAuthLogin(ctx context.Context, profile OAuthProfile) (*AuthResult, error) {
	email := strings.TrimSpace(strings.ToLower(profile.Email))
	if email == "" {
		return nil, domain.ErrInvalidInput
	}

	name := strings.TrimSpace(profile.Name)
	if name == "" {
		name = strings.Split(email, "@")[0]
	}

	var avatarPtr *string
	if profile.Avatar != "" {
		avatarPtr = &profile.Avatar
	}

	var googleIDPtr *string
	if profile.GoogleID != "" {
		googleIDPtr = &profile.GoogleID
	}

	// 1. Try to find by Google ID
	if profile.GoogleID != "" {
		user, err := s.userRepo.FindByGoogleID(ctx, profile.GoogleID)
		if err == nil && user != nil {
			if avatarPtr != nil && (user.Avatar == nil || *user.Avatar != *avatarPtr) {
				user.Avatar = avatarPtr
				user.UpdatedAt = time.Now()
				_ = s.userRepo.Update(ctx, user)
			}
			token, _ := s.generateToken(user.ID, user.Email, user.Username, user.Name, "access", config.GetConfig().JWT.AccessTokenDuration)
			refreshToken, _ := s.generateToken(user.ID, user.Email, user.Username, user.Name, "refresh", config.GetConfig().JWT.RefreshTokenDuration)
			return &AuthResult{
				User:         user,
				AccessToken:  token,
				RefreshToken: refreshToken,
				IsNewUser:    false,
			}, nil
		}
	}

	// 2. Try to find by Email (Link Google account to existing user)
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		existingUser.GoogleID = googleIDPtr
		if existingUser.Avatar == nil && avatarPtr != nil {
			existingUser.Avatar = avatarPtr
		}
		existingUser.AuthProvider = domain.AuthProviderGoogle
		existingUser.UpdatedAt = time.Now()
		_ = s.userRepo.Update(ctx, existingUser)

		token, _ := s.generateToken(existingUser.ID, existingUser.Email, existingUser.Username, existingUser.Name, "access", config.GetConfig().JWT.AccessTokenDuration)
		refreshToken, _ := s.generateToken(existingUser.ID, existingUser.Email, existingUser.Username, existingUser.Name, "refresh", config.GetConfig().JWT.RefreshTokenDuration)
		return &AuthResult{
			User:         existingUser,
			AccessToken:  token,
			RefreshToken: refreshToken,
			IsNewUser:    false,
		}, nil
	}

	// 3. Create new user with Google OAuth
	now := time.Now()
	username := strings.Split(email, "@")[0]
	newUser := &domain.User{
		ID:           uuid.New().String(),
		Email:        email,
		Username:     username,
		Name:         name,
		GoogleID:     googleIDPtr,
		AuthProvider: domain.AuthProviderGoogle,
		Avatar:       avatarPtr,
		Password:     nil,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		s.log.Error("failed to create oauth user", zap.Error(err))
		return nil, err
	}

	// Create default categories and wallets
	if s.categoryRepo != nil {
		_ = s.categoryRepo.CreateDefaultCategories(ctx, newUser.ID)
	}
	if s.walletRepo != nil {
		_ = s.walletRepo.CreateDefaultWallets(ctx, newUser.ID)
	}

	token, _ := s.generateToken(newUser.ID, newUser.Email, newUser.Username, newUser.Name, "access", config.GetConfig().JWT.AccessTokenDuration)
	refreshToken, _ := s.generateToken(newUser.ID, newUser.Email, newUser.Username, newUser.Name, "refresh", config.GetConfig().JWT.RefreshTokenDuration)

	return &AuthResult{
		User:         newUser,
		AccessToken:  token,
		RefreshToken: refreshToken,
		IsNewUser:    true,
	}, nil
}

func (s *service) ValidateToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	jwtSecret := config.GetConfig().JWT.Secret

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	userID, _ := claims["user_id"].(string)
	email, _ := claims["email"].(string)
	username, _ := claims["username"].(string)
	name, _ := claims["name"].(string)
	tokenType, _ := claims["token_type"].(string)

	if userID == "" {
		return nil, domain.ErrInvalidToken
	}

	return &TokenClaims{
		UserID:    userID,
		Email:     email,
		Username:  username,
		Name:      name,
		TokenType: tokenType,
	}, nil
}

func (s *service) RefreshToken(ctx context.Context, tokenString string) (*AuthResult, error) {
	claims, err := s.ValidateToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	token, err := s.generateToken(user.ID, user.Email, user.Username, user.Name, "access", config.GetConfig().JWT.AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user.ID, user.Email, user.Username, user.Name, "refresh", config.GetConfig().JWT.RefreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         user,
		AccessToken:  token,
		RefreshToken: refreshToken,
		IsNewUser:    false,
	}, nil
}

func (s *service) generateToken(userID, email, username, name, tokenType string, duration time.Duration) (string, error) {
	jwtSecret := config.GetConfig().JWT.Secret

	claims := jwt.MapClaims{
		"jti":        uuid.New().String(),
		"user_id":    userID,
		"userId":     userID, // Support both snake_case and camelCase for client compatibility
		"email":      email,
		"username":   username,
		"name":       name,
		"token_type": tokenType,
		"exp":        time.Now().Add(duration).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}


