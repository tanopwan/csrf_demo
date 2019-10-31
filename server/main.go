package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"
)

var sessionID int
var sessions map[int]*UserProfile
var domain string

type UserProfile struct {
	Name    string
	Balance int
}

func getEnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	envCORS := getEnvOrDefault("ENABLE_CORS", "false")

	domainPtr := flag.String("domain", "localhost", "domain of the cookie")
	domain = *domainPtr

	flag.Parse()
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	sessionID = 0
	sessions = make(map[int]*UserProfile)
	e := echo.New()
	e.Renderer = t
	if envCORS == "true" {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowCredentials: true,
		}))
	}

	e.GET("/result", resultPage)
	e.GET("/login", loginPage)
	e.GET("/transfer/level1", transferLevel1Page)
	e.GET("/transfer/level2", transferLevel2Page)
	e.GET("/transfer/level2/1", transferLevel21Page)
	e.GET("/transfer/level2/2", transferLevel22Page)
	e.GET("/transfer/level3", transferLevel3Page)
	e.POST("/api/login", login)
	e.POST("/api/logout", logout)
	e.GET("/api/transfer", transferGet)
	e.POST("/api/transfer2", transferPost)
	e.POST("/api/transfer2/1", transferPostRedirect)
	e.PUT("/api/transfer2/2", transferPost)
	e.POST("/api/transfer3", transferPostJSON)
	e.GET("/", indexPage)
	e.Static("/", "public")

	port := getEnvOrDefault("PORT", "1323")
	e.Logger.Fatal(e.Start(":" + port))
}

func indexPage(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func logout(c echo.Context) error {
	session, err := c.Cookie("session")
	if err != nil {
		fmt.Println("cookie is not found,", err)
		return c.Redirect(http.StatusFound, "/login")
	}
	fmt.Printf("session: %+v\n", session)
	session.Expires = time.Now()
	session.Path = "/"
	session.Domain = domain
	c.SetCookie(session)
	return c.Redirect(http.StatusFound, "/login")
}

func login(c echo.Context) error {
	username := c.Request().PostFormValue("username")
	password := c.Request().PostFormValue("password")

	if username == "golf" && password == "golf" {
		sessions[sessionID] = &UserProfile{
			Name:    "Tanopwan",
			Balance: 50000,
		}
		session := http.Cookie{
			Name:   "session",
			Value:  strconv.Itoa(sessionID),
			Path:   "/",
			Domain: domain,
		}
		c.SetCookie(&session)
		sessionID++
		return c.Redirect(http.StatusFound, "/")
	}
	return c.Redirect(http.StatusFound, "/login")
}

func resultPage(c echo.Context) error {
	userProfile := validateSession(c)
	data := struct {
		Balance int
	}{
		Balance: userProfile.Balance,
	}
	return c.Render(http.StatusOK, "resultLv2", data)
}

func loginPage(c echo.Context) error {
	userProfile := validateSession(c)
	if userProfile.Name != "" {
		return c.Redirect(http.StatusFound, "/transfer/level1")
	}

	return c.Render(http.StatusOK, "login", "Please login")
}

func validateSession(c echo.Context) *UserProfile {
	session, err := c.Cookie("session")
	if err != nil {
		fmt.Println("cookie is not found,", err)
		return &UserProfile{}
	}
	value, err := strconv.Atoi(session.Value)
	if err != nil {
		fmt.Println("session is not number,", err)
		return &UserProfile{}
	}
	fmt.Printf("sessions: %+v\n", sessions)
	if userProfile, ok := sessions[value]; ok {
		return userProfile
	}

	fmt.Println("session is not found with id ", value)
	return &UserProfile{}
}

func transferLevel1Page(c echo.Context) error {
	userProfile := validateSession(c)
	if userProfile.Name != "" {
		return c.Render(http.StatusOK, "transferLv1", userProfile.Name)
	}

	return c.Redirect(http.StatusFound, "/login")
}

func transferLevel2Page(c echo.Context) error {
	userProfile := validateSession(c)
	if userProfile.Name != "" {
		return c.Render(http.StatusOK, "transferLv2", userProfile.Name)
	}

	return c.Redirect(http.StatusFound, "/login")
}

func transferLevel21Page(c echo.Context) error {
	userProfile := validateSession(c)
	if userProfile.Name != "" {
		return c.Render(http.StatusOK, "transferLv2_1", userProfile.Name)
	}

	return c.Redirect(http.StatusFound, "/login")
}

func transferLevel22Page(c echo.Context) error {
	userProfile := validateSession(c)
	if userProfile.Name != "" {
		return c.Render(http.StatusOK, "transferLv2_2", userProfile.Name)
	}

	return c.Redirect(http.StatusFound, "/login")
}

func transferLevel3Page(c echo.Context) error {
	userProfile := validateSession(c)
	if userProfile.Name != "" {
		return c.Render(http.StatusOK, "transferLv3", userProfile.Name)
	}

	return c.Redirect(http.StatusFound, "/login")
}

func transferGet(c echo.Context) error {
	userProfile := validateSession(c)
	if userProfile.Name == "" {
		return c.Redirect(http.StatusFound, "/login")
	}
	toParam := c.QueryParam("to")
	amountParam := c.QueryParam("amount")
	amount, _ := strconv.Atoi(amountParam)
	userProfile.Balance = userProfile.Balance - amount
	data := struct {
		To      string
		Amount  int
		Balance int
	}{
		To:      toParam,
		Amount:  amount,
		Balance: userProfile.Balance,
	}
	return c.Render(http.StatusOK, "result", data)
}

func transferPost(c echo.Context) error {
	userProfile := validateSession(c)
	if userProfile.Name == "" {
		return c.Redirect(http.StatusFound, "/login")
	}
	toParam := c.FormValue("to")
	amountParam := c.FormValue("amount")
	amount, _ := strconv.Atoi(amountParam)

	userProfile.Balance = userProfile.Balance - amount
	data := struct {
		To      string
		Amount  int
		Balance int
	}{
		To:      toParam,
		Amount:  amount,
		Balance: userProfile.Balance,
	}
	return c.Render(http.StatusOK, "result", data)
}

func transferPostRedirect(c echo.Context) error {
	userProfile := validateSession(c)
	if userProfile.Name == "" {
		return c.Redirect(http.StatusFound, "/login")
	}
	amountParam := c.FormValue("amount")
	amount, _ := strconv.Atoi(amountParam)

	userProfile.Balance = userProfile.Balance - amount
	return c.Redirect(http.StatusSeeOther, "/result")
}

func transferPostJSON(c echo.Context) error {
	userProfile := validateSession(c)
	if userProfile.Name == "" {
		return c.Redirect(http.StatusFound, "/login")
	}
	body := struct {
		To     string `json:"to"`
		Amount int    `json:"amount"`
	}{}
	err := c.Bind(&body)
	if err != nil {
		c.String(http.StatusBadRequest, "error reading body "+err.Error())
	}

	userProfile.Balance = userProfile.Balance - body.Amount
	return c.Redirect(http.StatusSeeOther, "/result")
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
