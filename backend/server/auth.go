package server

import (
	"context"
	"lab-assignment-system-backend/lib"
	"lab-assignment-system-backend/repository"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SignupForm struct {
	Email         string `json:"email,omitempty"`
	Password      string `json:"password,omitempty"`
	StudentNumber string `json:"studentNumber,omitempty"`
	Name          string `json:"name,omitempty"`
	IdToken       string `json:"idToken,omitempty"`
	Lab1          string `json:"lab1,omitempty"`
	Lab2          string `json:"lab2,omitempty"`
	Lab3          string `json:"lab3,omitempty"`
	Token         string `json:"token"`
}

type SigninForm struct {
	UID      string `json:"uid,omitempty"`
	Password string `json:"password,omitempty"`
}

const sessionExpiresIn = time.Hour * 24 * 7

func (srv *Server) AuthRouter() {
	r := srv.r.Group("/auth")
	{
		r.POST("/signin", srv.HandleSignin())
		r.POST("/signout", srv.HandleSignout()).Use(srv.Authentication())
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

		var user repository.User
		if err := srv.dc.Get(c.Request.Context(), repository.NewUserKey(signinForm.UID), &user); err != nil {
			srv.logger.Printf("%+v\n", err)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(signinForm.Password)); err != nil {
			srv.logger.Printf("%+v\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		now := time.Now()
		sessionValue := lib.MakeRandomString(32)
		repository.NewSession(user.UID, sessionValue, now, now.Add(sessionExpiresIn))
		sessionCookie, err := makeSessionCookie(ctx, sessionValue)
		if err != nil {
			srv.logger.Printf("%+v\n", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		http.SetCookie(c.Writer, sessionCookie)
	}
}

func (srv *Server) HandleSignout() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionCookie, err := c.Request.Cookie("session")
		if err != nil {
			lib.AbortWithErrorJSON(c, lib.NewError(http.StatusBadRequest, "no session cookie"))
			return
		}
		sessionCookie.Value = ""
		sessionCookie.MaxAge = 0
		sessionCookie.Path = "/"
		http.SetCookie(c.Writer, sessionCookie)
	}
}

func makeSessionCookie(ctx context.Context, session string) (*http.Cookie, error) {
	return &http.Cookie{
		Name:     "session",
		Value:    session,
		MaxAge:   int(sessionExpiresIn),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}, nil
}
