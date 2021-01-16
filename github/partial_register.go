package github

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

type partialRegisterVariables struct {
	Company string
	Email   string
	Name    string
	Token   string
}

func PartialRegister(ctx *gin.Context) {
	company := ctx.Query("company")
	email := ctx.Query("email")
	name := ctx.Query("name")
	token := ctx.Query("token")

	parsedTemplate, err := template.ParseFiles("./partial_register.gohtml")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("error parsing template: %w", err),
		})
		return
	}

	variables := partialRegisterVariables{
		Company: company,
		Email:   email,
		Name:    name,
		Token:   token,
	}

	if err := parsedTemplate.Execute(ctx.Writer, variables); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("error rendering template: %w", err),
		})
	}
}
