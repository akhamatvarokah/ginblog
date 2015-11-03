package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginblog/helpers"
	"github.com/denisbakhtin/ginblog/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//SignInGet handles GET /signin route
func SignInGet(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "Basic GIN web-site signin form"
	h["Active"] = "signin"
	session := helpers.GetSession(c)
	h["Flash"] = session.Flashes()
	session.Save(c.Request, c.Writer)
	c.HTML(http.StatusOK, "auth/signin", h)
}

//SignInPost handles POST /signin route, authenticates user
func SignInPost(c *gin.Context) {
	session := helpers.GetSession(c)
	user := &models.User{}
	if err := c.Bind(user); err != nil {
		session.AddFlash("Please, fill out form correctly.")
		session.Save(c.Request, c.Writer)
		c.Redirect(http.StatusFound, "/signin")
		return
	}

	userDB, _ := models.GetUserByEmail(user.Email)
	if userDB.ID == 0 {
		logrus.Errorf("Login error, IP: %s, Email: %s", c.ClientIP(), user.Email)
		session.AddFlash("Email or password incorrect")
		session.Save(c.Request, c.Writer)
		c.Redirect(http.StatusFound, "/signin")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password)); err != nil {
		logrus.Errorf("Login error, IP: %s, Email: %s", c.ClientIP(), user.Email)
		session.AddFlash("Email or password incorrect")
		session.Save(c.Request, c.Writer)
		c.Redirect(http.StatusFound, "/signin")
		return
	}

	session.Values["UserID"] = userDB.ID
	session.Save(c.Request, c.Writer)
	c.Redirect(http.StatusFound, "/")
}

//SignUpGet handles GET /signup route
func SignUpGet(c *gin.Context) {
	h := helpers.DefaultH(c)
	h["Title"] = "Basic GIN web-site signup form"
	h["Active"] = "signup"
	session := helpers.GetSession(c)
	h["Flash"] = session.Flashes()
	session.Save(c.Request, c.Writer)
	c.HTML(http.StatusOK, "auth/signup", h)
}

//SignUpPost handles POST /signup route, creates new user
func SignUpPost(c *gin.Context) {
	session := helpers.GetSession(c)
	user := &models.User{}
	if err := c.Bind(user); err != nil {
		session.AddFlash(err.Error())
		session.Save(c.Request, c.Writer)
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	userDB, _ := models.GetUserByEmail(user.Email)
	if userDB.ID != 0 {
		session.AddFlash("User exists")
		session.Save(c.Request, c.Writer)
		c.Redirect(http.StatusFound, "/signup")
		return
	}
	//create user
	err := user.HashPassword()
	if err != nil {
		session.AddFlash("Error whilst registering user.")
		session.Save(c.Request, c.Writer)
		logrus.Errorf("Error whilst registering user: %v", err)
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	if err := user.Insert(); err != nil {
		session.AddFlash("Error whilst registering user.")
		session.Save(c.Request, c.Writer)
		logrus.Errorf("Error whilst registering user: %v", err)
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	session.Values["UserID"] = user.ID
	session.Save(c.Request, c.Writer)
	c.Redirect(http.StatusFound, "/")
	return
}

//LogoutGet handles GET /logout route
func LogoutGet(c *gin.Context) {
	session := helpers.GetSession(c)
	delete(session.Values, "UserID")
	session.Save(c.Request, c.Writer)
	c.Redirect(http.StatusSeeOther, "/")
}
