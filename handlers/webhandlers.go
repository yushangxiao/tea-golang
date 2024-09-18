package handlers

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func WebHandler(c *gin.Context) {
	if c.Request.URL.Path == "/" {
		c.File("./web/index.html")
		return
	}
	// 如果没有后缀，就默认加上.html
	if !strings.Contains(c.Request.URL.Path, ".") {
		c.File("./web" + c.Request.URL.Path + ".html")
		return
	}

	if strings.HasSuffix(c.Request.URL.Path, ".jpg") {
		c.Writer.Header().Set("Content-Encoding", "gzip")
		c.Writer.Header().Set("Content-Type", "image/jpeg")
		c.Writer.Header().Set("Cache-Control", "public, max-age=3600")
		c.Status(200)
		file, err := os.Open("./web" + c.Request.URL.Path)
		if err != nil {
			fmt.Println("open file error: ", err.Error())
			c.String(500, err.Error())
			return
		}
		defer file.Close()

		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()

		_, err = io.Copy(gz, file)
		if err != nil {
			fmt.Println("copy file error: ", err.Error())
			c.String(500, err.Error())
			return
		}
		// fmt.Println("copy file success")
		return
	}

	c.File("./web" + c.Request.URL.Path)
}
