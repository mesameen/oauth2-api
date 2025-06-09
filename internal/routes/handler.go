package routes

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	"github.com/mesameen/oauth2-api/internal/config"
	"github.com/mesameen/oauth2-api/internal/logger"
)

type Handler struct {
	cookieStore *sessions.CookieStore
}

func NewHandler() (*Handler, error) {
	var store = sessions.NewCookieStore([]byte(config.OAuthConfig.SessionKey))
	return &Handler{
		cookieStore: store,
	}, nil
}

func (h *Handler) Home(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		logger.Errorf("Failed to pase login html file. Error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(c.Writer, gin.H{})
	if err != nil {
		logger.Errorf("Failed to execute template. Error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SignInWithProvider(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (h *Handler) AuthProviderCallback(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		logger.Errorf("Failed to complete user auth: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	logger.Infof("token :%+v", user.AccessToken)
	logger.Infof("secret :%+v", user.AccessTokenSecret)
	// setting session name to cookies
	session, err := h.cookieStore.Get(c.Request, "session-name")
	if err != nil {
		logger.Errorf("Failed to get the session. Error: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	session.Values["auth_token"] = user.AccessToken
	session.Save(c.Request, c.Writer)
	c.Redirect(http.StatusTemporaryRedirect, "/api/success")
}

func (h *Handler) Success(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf(`<div style="
          background-color: #fff;
          padding: 40px;
          border-radius: 8px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
          text-align: center;
      ">
          <h1 style="
              color: #333;
              margin-bottom: 20px;
          ">You have Successfull signed in!</h1>
        <a href="/api/logout" style="
        display: inline-flex; 
        align-items: center; 
        justify-content: center; 
        background-color: #4285f4; 
        color: #fff; 
        text-decoration: none; 
        padding: 12px 20px; 
        border-radius: 4px; 
        transition: background-color 0.3s ease;">
            <span style="font-size: 16px; font-weight: bold;">Logout</span>
        </a>
        <a href="/api/protected" style="
        display: inline-flex; 
        align-items: center; 
        justify-content: center; 
        background-color: #4285f4; 
        color: #fff; 
        text-decoration: none; 
        padding: 12px 20px; 
        border-radius: 4px; 
        transition: background-color 0.3s ease;">
            <span style="font-size: 16px; font-weight: bold;">Protected</span>
        </a>
          </div>
          </div>
		  </div>`)))
}

func (h *Handler) Logout(c *gin.Context) {
	session, _ := h.cookieStore.Get(c.Request, "session-name")
	session.Options.MaxAge = -1 // deleting the sesion
	session.Save(c.Request, c.Writer)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (h *Handler) Protected(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "protected data"})
}
