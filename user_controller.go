package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"html/template"
	"io"
	"models"
	"net/http"
	"strconv"
)

// gorm Instance
var db *gorm.DB

// レイアウト適用済みのテンプレートを保存するmap
var templates map[string]*template.Template

// TemplateはHTMLテンプレートを利用するためのRenderer Interface
type Template struct{}

// RenderはHTMLテンプレートにデータを埋め込んだ結果をWriterに書き込む
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return templates[name].ExecuteTemplate(w, "layout.html", data)
}

func loadTemplates() {
	templates = make(map[string]*template.Template)
	templates["index"] = template.Must(template.ParseFiles("templates/layout.html", "views/index.html"))
}

func SetTemplate(html string) {
	templates[html] = template.Must(template.ParseFiles("templates/layout.html", "views/"+html+".html"))
}

func init() {

	var err error
	db, err = gorm.Open("mysql", "root@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	loadTemplates()
}

func main() {

	// Echo instance
	e := echo.New()
	t := &Template{}
	e.Renderer = t

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 1}))

	// 各ルーティングに対するハンドラを設定
	e.GET("/users", Users)
	e.GET("/users/:id", UsersGet)
	e.POST("/users/:id", UsersUpdate)

	// API
	e.GET("/api/users", APIUsers)
	e.GET("/api/users/:id", APIUsersGet)

	// start server
	e.Logger.Fatal(e.Start(":1324"))
}

func Users(c echo.Context) error {

	// usersインスタンス化
	users := []models.User{}
	db.Find(&users)
	data := map[string]interface{}{"users": users}
	return c.Render(http.StatusOK, "index", data)
}

func UsersGet(c echo.Context) error {

	id := c.Param("id")
	if id == "" {
		panic("no value for param requested")
	}

	user := models.User{}
	db.Find(&user, id)
	data := map[string]interface{}{"user": user}

	// set html to Template
	SetTemplate("show")
	return c.Render(http.StatusOK, "show", data)
}

func UsersUpdate(c echo.Context) error {

	id := c.FormValue("id")
	name := c.FormValue("name")
	age, _ := strconv.Atoi(c.FormValue("age"))
	fmt.Println(name, age)

	tx := db.Begin()
	err := tx.Model(&models.User{}).Where("id = ?", 0).Updates(models.User{Name: name, Age: age}).Error
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()

	return c.Redirect(http.StatusMovedPermanently, "/users/"+id)
}

func APIUsers(c echo.Context) error {

	users := []models.User{}
	db.Find(&users)
	data := map[string]interface{}{"users": users}
	return c.JSON(http.StatusOK, data)
}

func APIUsersGet(c echo.Context) error {

	id := c.Param("id")
	if id == "" {
		panic("no value for param requested")
	}

	user := models.User{}
	db.Find(&user, id)
	data := map[string]interface{}{"user": user}

	return c.JSON(http.StatusOK, data)
}
