package main

import (
	"tinyGin"
	"net/http"
)

// run后在另一个console查看
func main() {
	r := tinyGin.Default()
	r.GET("/", func(c *tinyGin.Context) {
		c.String(http.StatusOK, "Hello suisbuds\n")
	})
	r.GET("/panic", func(c *tinyGin.Context) {
		names := []string{"suisbuds"}
		c.String(http.StatusOK, names[100])
	})
	r.Run(":9999")
}
