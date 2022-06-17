package server

import (
	"context"
	"lab-assignment-system-backend/lib"
	"lab-assignment-system-backend/repository"
	"log"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
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
}

type SigninForm struct {
	IdToken string `json:"idToken,omitempty"`
}

func (srv *Server) AuthRouter() {
	gradesRouter := srv.r.Group("/auth")
	{
		gradesRouter.POST("/signin", srv.HandleSignin())
		gradesRouter.POST("/signup", srv.HandleSignup())
		gradesRouter.POST("/signout", srv.HandleSignout()).Use(srv.Authentication())
	}
}

func (srv *Server) HandleSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var signupForm SignupForm
		if err := c.BindJSON(&signupForm); err != nil {
			srv.logger.Printf("%+v\n", err)
			return
		}
		if !validateEmail(signupForm.Email) || len(signupForm.Password) < 8 {
			lib.AbortWithErrorJSON(c, lib.NewError(http.StatusBadRequest, "invalid email or password"))
			return
		}
		token, err := srv.auth.VerifyIDToken(ctx, signupForm.IdToken)
		if err != nil {
			srv.logger.Printf("%+v\n", err)
			lib.AbortWithErrorJSON(c, lib.NewError(http.StatusUnauthorized, "not logged in"))
			return
		}
		userdata, err := srv.auth.GetUser(ctx, token.UID)
		if err != nil {
			srv.logger.Printf("%+v\n", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		user, userKey := repository.NewUser(
			userdata.UID,
			userdata.Email,
			signupForm.StudentNumber,
			signupForm.Name,
			signupForm.Lab1,
			signupForm.Lab2,
			signupForm.Lab3,
		)
		labIds := []string{signupForm.Lab1, signupForm.Lab2, signupForm.Lab3}
		if _, err := srv.dc.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
			mutations := make([]*datastore.Mutation, 0, 4)
			mutations = append(mutations, datastore.NewInsert(userKey, user))
			for i, labId := range labIds {
				var lab repository.Lab
				key := repository.NewLabKey(labId)
				if err := tx.Get(key, &lab); err != nil {
					return err
				}
				if i == 0 {
					lab.FirstChoice++
				} else if i == 1 {
					lab.SecondChoice++
				} else {
					lab.ThirdChice++
				}
				mutations = append(mutations, datastore.NewUpdate(key, &lab))
			}
			if _, err := tx.Mutate(mutations...); err != nil {
				return err
			}
			return nil
		}); err != nil {
			srv.logger.Printf("%+v\n", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		sessionCookie, err := makeSessionCookie(ctx, srv.auth, signupForm.IdToken)
		if err != nil {
			srv.logger.Printf("%+v\n", err)
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
		SameSite: http.SameSiteNoneMode,
	}, nil
}
