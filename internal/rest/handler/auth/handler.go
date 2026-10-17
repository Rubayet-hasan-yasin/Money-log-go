package authhandler

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/auth"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/config"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/rest/middleware"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/rest/response"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Handler struct {
	svc         auth.Service
	log         *zap.Logger
	oauthConfig *oauth2.Config
}

func NewHandler(svc auth.Service, log *zap.Logger, oauthConfig *oauth2.Config) *Handler {
	if oauthConfig == nil {
		cfg := config.GetConfig()
		if cfg.OAuth.GoogleClientID != "" && cfg.OAuth.GoogleClientSecret != "" && cfg.OAuth.GoogleCallbackURL != "" {
			oauthConfig = &oauth2.Config{
				ClientID:     cfg.OAuth.GoogleClientID,
				ClientSecret: cfg.OAuth.GoogleClientSecret,
				RedirectURL:  cfg.OAuth.GoogleCallbackURL,
				Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.email",
					"https://www.googleapis.com/auth/userinfo.profile",
				},
				Endpoint: google.Endpoint,
			}
		}
	}

	return &Handler{
		svc:         svc,
		log:         log,
		oauthConfig: oauthConfig,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.svc.Register(c.Request.Context(), auth.RegisterInput{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "User registered successfully",
		"user":         ToUserDTO(result.User),
		"accessToken":  result.AccessToken,
		"refreshToken": result.RefreshToken,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.svc.Login(c.Request.Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"user":         ToUserDTO(result.User),
		"accessToken":  result.AccessToken,
		"refreshToken": result.RefreshToken,
	})
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Unauthorized(c, "unauthorized")
		return
	}

	user, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, ToUserDTO(user))
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Unauthorized(c, "unauthorized")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.svc.UpdateProfile(c.Request.Context(), userID, auth.UpdateProfileInput{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
		Avatar:   req.Avatar,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    ToUserDTO(user),
	})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.svc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Token refreshed successfully",
		"user":         ToUserDTO(result.User),
		"accessToken":  result.AccessToken,
		"refreshToken": result.RefreshToken,
	})
}

func (h *Handler) GoogleLogin(c *gin.Context) {
	if h.oauthConfig == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Google OAuth is not configured"})
		return
	}

	returnTo := c.Query("returnTo")
	if returnTo == "" {
		returnTo = c.Query("redirectUrl")
	}

	b := make([]byte, 16)
	_, _ = rand.Read(b)
	nonce := base64.RawURLEncoding.EncodeToString(b)

	statePayload := map[string]string{
		"nonce":    nonce,
		"returnTo": returnTo,
	}
	stateBytes, _ := json.Marshal(statePayload)
	state := base64.RawURLEncoding.EncodeToString(stateBytes)

	// Save state in HTTP-only cookie for CSRF validation
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	authURL := h.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *Handler) GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	redirectURL := extractRedirectURL(state)

	if errParam := c.Query("error"); errParam != "" {
		h.log.Error("google oauth returned error", zap.String("error", errParam))
		handleOAuthError(c, redirectURL, errParam, http.StatusBadRequest)
		return
	}

	if h.oauthConfig == nil {
		h.log.Error("google oauth is not configured")
		handleOAuthError(c, redirectURL, "Google OAuth not configured", http.StatusServiceUnavailable)
		return
	}

	cookieState, _ := c.Cookie("oauth_state")
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	if cookieState != "" && cookieState != state {
		h.log.Error("oauth state mismatch")
		handleOAuthError(c, redirectURL, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	code := c.Query("code")
	if code == "" {
		h.log.Error("missing authorization code from google")
		handleOAuthError(c, redirectURL, "Missing authorization code", http.StatusBadRequest)
		return
	}

	oauthToken, err := h.oauthConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		h.log.Error("google oauth code exchange failed", zap.Error(err))
		handleOAuthError(c, redirectURL, "Authentication failed", http.StatusBadRequest)
		return
	}

	client := h.oauthConfig.Client(c.Request.Context(), oauthToken)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		h.log.Error("failed to get user info from google", zap.Error(err))
		handleOAuthError(c, redirectURL, "Failed to retrieve user profile", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		h.log.Error("google userinfo returned non-200", zap.Int("status", resp.StatusCode))
		handleOAuthError(c, redirectURL, "Failed to retrieve user profile", http.StatusInternalServerError)
		return
	}

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		h.log.Error("failed to parse google userinfo response", zap.Error(err))
		handleOAuthError(c, redirectURL, "Failed to parse user profile", http.StatusInternalServerError)
		return
	}

	result, err := h.svc.OAuthLogin(c.Request.Context(), auth.OAuthProfile{
		Provider: "google",
		Email:    googleUser.Email,
		Name:     googleUser.Name,
		Avatar:   googleUser.Picture,
		GoogleID: googleUser.ID,
	})
	if err != nil {
		h.log.Error("oauth login failed", zap.Error(err))
		handleOAuthError(c, redirectURL, "Authentication failed", http.StatusInternalServerError)
		return
	}

	if redirectURL != "" {
		c.Redirect(http.StatusTemporaryRedirect, appendTokensToURL(redirectURL, result.AccessToken, result.RefreshToken))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "OAuth login successful",
		"user":         ToUserDTO(result.User),
		"accessToken":  result.AccessToken,
		"refreshToken": result.RefreshToken,
	})
}

func (h *Handler) GoogleFailure(c *gin.Context) {
	redirectURL := c.Query("redirectUrl")
	if redirectURL == "" {
		redirectURL = c.Query("returnTo")
	}
	handleOAuthError(c, redirectURL, "Authentication failed", http.StatusUnauthorized)
}

func extractRedirectURL(state string) string {
	if state == "" {
		return ""
	}
	var stateObj struct {
		ReturnTo string `json:"returnTo"`
	}
	decoded, decodeErr := base64.RawURLEncoding.DecodeString(state)
	if decodeErr != nil {
		decoded, decodeErr = base64.URLEncoding.DecodeString(state)
	}
	if decodeErr != nil {
		decoded, decodeErr = base64.StdEncoding.DecodeString(state)
	}
	if decodeErr == nil && json.Unmarshal(decoded, &stateObj) == nil {
		return stateObj.ReturnTo
	}
	return ""
}

func appendTokensToURL(targetURL, accessToken, refreshToken string) string {
	parsed, err := url.Parse(targetURL)
	if err != nil {
		separator := "?"
		if strings.Contains(targetURL, "?") {
			separator = "&"
		}
		return targetURL + separator + "accessToken=" + url.QueryEscape(accessToken) + "&refreshToken=" + url.QueryEscape(refreshToken)
	}
	q := parsed.Query()
	q.Set("accessToken", accessToken)
	q.Set("refreshToken", refreshToken)
	parsed.RawQuery = q.Encode()
	return parsed.String()
}

func appendErrorToURL(targetURL, errMsg string) string {
	parsed, err := url.Parse(targetURL)
	if err != nil {
		separator := "?"
		if strings.Contains(targetURL, "?") {
			separator = "&"
		}
		return targetURL + separator + "error=" + url.QueryEscape(errMsg)
	}
	q := parsed.Query()
	q.Set("error", errMsg)
	parsed.RawQuery = q.Encode()
	return parsed.String()
}

func handleOAuthError(c *gin.Context, redirectURL, errMsg string, statusCode int) {
	if redirectURL != "" {
		c.Redirect(http.StatusTemporaryRedirect, appendErrorToURL(redirectURL, errMsg))
		return
	}
	c.JSON(statusCode, gin.H{"error": errMsg})
}
