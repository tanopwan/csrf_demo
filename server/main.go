package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

	echo "github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"
)

var sessionID int
var sessions map[int]string
var myBalance int

func main() {
	corsPtr := flag.Bool("cors", false, "true to enable CORS")
	credPtr := flag.Bool("cred", false, "true to enable credentials with")
	flag.Parse()
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	myBalance = 50000
	sessionID = 0
	sessions = make(map[int]string)
	e := echo.New()
	e.Renderer = t
	if *corsPtr && *credPtr {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowCredentials: true,
		}))
	} else {
		e.Use(middleware.CORS())
	}

	e.GET("/", hello)
	e.GET("/login", loginPage)
	e.GET("/transfer", transferPage)
	e.POST("/api/login", login)
	e.POST("/api/logout", logout)
	e.GET("/api/transfer", transferGet)
	e.POST("/api/transfer2", transferPost)
	e.POST("/api/transfer3", transferPostJSON)
	e.PUT("/api/transfer2", transferPost)
	e.Logger.Fatal(e.Start(":1323"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
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
	session.Domain = "domain1.com"
	c.SetCookie(session)
	return c.Redirect(http.StatusFound, "/login")
}

func login(c echo.Context) error {
	username := c.Request().PostFormValue("username")
	password := c.Request().PostFormValue("password")

	if username == "golf" && password == "golf" {
		sessions[sessionID] = "Tanopwan"
		session := http.Cookie{
			Name:   "session",
			Value:  strconv.Itoa(sessionID),
			Path:   "/",
			Domain: "domain1.com",
		}
		c.SetCookie(&session)
		sessionID++
		return c.Redirect(http.StatusFound, "/transfer")
	}
	return c.Redirect(http.StatusFound, "/login")
}

func loginPage(c echo.Context) error {
	username := validateSession(c)
	if username != "" {
		return c.Redirect(http.StatusFound, "/transfer")
	}

	return c.Render(http.StatusOK, "login", "Please login")
}

func validateSession(c echo.Context) string {
	session, err := c.Cookie("session")
	if err != nil {
		fmt.Println("cookie is not found,", err)
		return ""
	}
	value, err := strconv.Atoi(session.Value)
	if err != nil {
		fmt.Println("session is not number,", err)
		return ""
	}
	fmt.Printf("sessions: %+v\n", sessions)
	if username, ok := sessions[value]; ok {
		return username
	}

	fmt.Println("session is not found with id ", value)
	return ""
}

func transferPage(c echo.Context) error {
	username := validateSession(c)
	if username != "" {
		return c.Render(http.StatusOK, "transfer", username)
	}

	return c.Redirect(http.StatusFound, "/login")
}

func transferGet(c echo.Context) error {
	username := validateSession(c)
	if username == "" {
		return c.Redirect(http.StatusFound, "/login")
	}
	toParam := c.QueryParam("to")
	amountParam := c.QueryParam("amount")
	amount, _ := strconv.Atoi(amountParam)

	myBalance = myBalance - amount
	return c.String(http.StatusOK, "successfully transfer "+amountParam+" baht to "+toParam+" You have "+strconv.Itoa(myBalance)+" baht left")
}

func transferPost(c echo.Context) error {
	username := validateSession(c)
	if username == "" {
		return c.Redirect(http.StatusFound, "/login")
	}
	toParam := c.FormValue("to")
	amountParam := c.FormValue("amount")
	amount, _ := strconv.Atoi(amountParam)

	myBalance = myBalance - amount
	return c.String(http.StatusOK, "successfully transfer "+amountParam+" baht to "+toParam+" You have "+strconv.Itoa(myBalance)+" baht left")
}

func transferPostJSON(c echo.Context) error {
	username := validateSession(c)
	if username == "" {
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

	myBalance = myBalance - body.Amount
	return c.String(http.StatusOK, "successfully transfer "+strconv.Itoa(body.Amount)+" baht to "+body.To+" You have "+strconv.Itoa(myBalance)+" baht left")
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
