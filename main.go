package main

import (
	"fmt"
	"gee"
	"net/http"
	"text/template"
	"time"
)

type student struct {
	Name string
	Age  int
}

func FormatAsData(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

// run后在另一个console查看
func main() {
	r := gee.New()
	r.Use(gee.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsData,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")
	stu1 := &student{Name: "suisbuds", Age: 20}
	stu2 := &student{Name: "sulis", Age: 22}
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *gee.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gee.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	r.GET("/date", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now":   time.Date(2020, 11, 10, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":9999")
}
