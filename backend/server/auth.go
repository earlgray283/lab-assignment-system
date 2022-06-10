package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

type SignupForm struct {
	Email         string `json:"email,omitempty"`
	Password      string `json:"password,omitempty"`
	StudentNumber int    `json:"studentNumber,omitempty"`
	Name          string `json:"name,omitempty"`
	IdToken       string `json:"idToken,omitempty"`
}

type SigninForm struct {
	IdToken string `json:"idToken,omitempty"`
}

func (srv *Server) AuthRouter() {
	gradesRouter := srv.r.Group("/auth")
	{
		gradesRouter.POST("/signin", srv.HandleSignin())
		gradesRouter.POST("/signup", srv.HandleSignup())
	}
}

func (srv *Server) HandleSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var signupForm SignupForm
		if err := c.BindJSON(&signupForm); err != nil {
			srv.logger.Println(err)
			return
		}
		if !validateEmail(signupForm.Email) || len(signupForm.Password) < 8 {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("invalid email or password"))
			return
		}

		token, err := srv.auth.VerifyIDToken(ctx, signupForm.IdToken)
		if err != nil {
			srv.logger.Println(err)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Invalid IdToken"))
			return
		}
		if err := srv.auth.SetCustomUserClaims(ctx, token.UID, map[string]interface{}{
			"studentNumber": strconv.Itoa(signupForm.StudentNumber),
			"name":          signupForm.Name,
		}); err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		sessionCookie, err := makeSessionCookie(ctx, srv.auth, signupForm.IdToken)
		if err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		http.SetCookie(c.Writer, sessionCookie)
	}
}

func (srv *Server) HandleSignin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var signinForm SigninForm
		if err := c.BindJSON(&signinForm); err != nil {
			log.Println(err)
			return
		}
		sessionCookie, err := makeSessionCookie(ctx, srv.auth, signinForm.IdToken)
		if err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		http.SetCookie(c.Writer, sessionCookie)
	}
}

var domainWhitelist = map[string]struct{}{
	"shizuoka.ac.jp":     {},
	"inf.shizuoka.ac.jp": {},
}

func validateEmail(email string) bool {
	tokens := strings.Split(email, "@")
	if len(tokens) != 2 {
		return false
	}
	_, ok := domainWhitelist[tokens[1]]
	return ok
}

func makeSessionCookie(ctx context.Context, auth *auth.Client, idToken string) (*http.Cookie, error) {
	const expiresIn = time.Hour * 24 * 7
	sessionCookie, err := auth.SessionCookie(ctx, idToken, expiresIn)
	if err != nil {
		return nil, err
	}
	return &http.Cookie{
		Name:     "session",
		Value:    sessionCookie,
		MaxAge:   int(expiresIn),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}, nil
}
