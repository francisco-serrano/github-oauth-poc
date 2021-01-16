package github

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func Index(ctx *gin.Context) {
	bs, err := ioutil.ReadFile("./sample.html")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Errorf("error while reading HTML: %w", err),
		})
		return
	}

	ctx.Data(http.StatusOK, "text/html; charset=utf-8", bs)
}
