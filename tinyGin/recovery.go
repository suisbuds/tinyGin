package tinyGin

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// Recovery 中间件的实现
// defer挂载在函数调用链的最外层，可以捕获到panic异常，并且堆栈信息打印出来
func Recovery() HandleFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		c.Next()
	}
}

// trace 用来获取panic的堆栈信息
func trace(message string) string {
	var pcs [32]uintptr             // 堆栈信息
	n := runtime.Callers(3, pcs[:]) // trace从调用栈的第3层开始,第0层是Callers,第1层是trace,第2层是defer函数
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	// 遍历调用栈信息
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		//打印文件名和行号
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
