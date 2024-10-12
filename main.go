package main

import (
	"fmt"
	"net/http"
	"time"
	"tinyGin"
	"html/template"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year,month,day := t.Date()
	return  fmt.Sprintf("%d-%02d-%02d",year,month,day)
}
// run后在另一个console查看
func main() {
	// r := tinyGin.Default()
	// r.GET("/", func(c *tinyGin.Context) {
	// 	c.String(http.StatusOK, "Hello suisbuds\n")
	// })
	// r.GET("/panic", func(c *tinyGin.Context) {
	// 	names := []string{"suisbuds"}
	// 	c.String(http.StatusOK, names[100])
	// })
	// r.Run(":9999")
	r := tinyGin.New()
	r.Use(tinyGin.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "tinyGinktutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/", func(c *tinyGin.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *tinyGin.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", tinyGin.H{
			"title":  "tinyGin",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *tinyGin.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", tinyGin.H{
			"title": "tinyGin",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":9999")
}
