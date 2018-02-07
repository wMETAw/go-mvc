package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"html/template"
	"io"
	"models"
	"net/http"
)

var db *gorm.DB

// レイアウト適用済みのテンプレートを保存するmap
var templates map[string]*template.Template

// TemplateはHTMLテンプレートを利用するためのRenderer Interface
type Template struct {
}

// RenderはHTMLテンプレートにデータを埋め込んだ結果をWriterに書き込む
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return templates[name].ExecuteTemplate(w, "layout.html", data)
}

func main() {

	// Echo instance
	e := echo.New()
	t := &Template{}
	e.Renderer = t

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 各ルーティングに対するハンドラを設定
	e.GET("/users", HandleUsers)
	e.GET("/api/users", HandleAPIUsers)

	// start server
	e.Logger.Fatal(e.Start(":1324"))
}

func init() {

	var err error
	db, err = gorm.Open("mysql", "root@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	loadTemplates()
}

func loadTemplates() {
	var baseTemplate = "templates/layout.html"
	templates = make(map[string]*template.Template)
	templates["index"] = template.Must(template.ParseFiles(baseTemplate, "templates/hello.html"))
}

func HandleUsers(c echo.Context) error {

	// usersインスタンス化
	users := []models.User{}
	db.Find(&users)
	data := map[string]interface{}{"users": users}
	return c.Render(http.StatusOK, "index", data)
}

func HandleAPIUsers(c echo.Context) error {

	// usersインスタンス化
	users := []models.User{}
	db.Find(&users)
	data := map[string]interface{}{"users": users}
	return c.JSON(http.StatusOK, data)
}
