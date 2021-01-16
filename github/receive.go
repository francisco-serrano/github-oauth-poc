package github

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/udacity/graphb"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type githubResponse struct {
	Data struct {
		Viewer githubData `json:"viewer"`
	} `json:"data"`
}

func ReceiveAuthCode(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	stateExpirationTime, ok := loginAttempts[state]
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "you shall not pass!!!",
		})
		return
	}

	if time.Now().After(stateExpirationTime) {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "login attempt expired",
		})
		return
	}

	authToken, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("error while obtaining authToken: %w", err),
		})
		return
	}

	query := graphb.Query{
		Type: graphb.TypeQuery,
		Name: "sample_query",
		Fields: []*graphb.Field{
			{
				Name:   "viewer",
				Fields: graphb.Fields("id", "company", "email", "name"),
			},
		},
	}

	jsonQuery, err := query.JSON()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("error while building graphQL query: %w", err),
		})
		return
	}

	requestBody := strings.NewReader(jsonQuery)

	client := oauthConfig.Client(context.Background(), authToken)

	res, err := client.Post("https://api.github.com/graphql", "application/json", requestBody)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("error while obtaining github data: %w", err),
		})
		return
	}

	defer res.Body.Close()

	var githubResponse githubResponse
	if err := json.NewDecoder(res.Body).Decode(&githubResponse); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("error while reading response body: %w", err),
		})
		return
	}

	externalUserID := githubResponse.Data.Viewer.ID

	sessionID := uuid.New().String()

	internalUserID, ok := oauthConnections[externalUserID]
	if !ok {
		internalUserID = uuid.New().String()

		oauthConnections[externalUserID] = internalUserID

		sessions[sessionID] = internalUserID

		token, err := CreateToken(sessionID, internalUserID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Errorf("error while generating token: %w", err),
			})
		}

		q := url.Values{}
		q.Add("company", githubResponse.Data.Viewer.Company)
		q.Add("email", githubResponse.Data.Viewer.Email)
		q.Add("name", githubResponse.Data.Viewer.Name)
		q.Add("token", token)

		partialRegisterURL := &url.URL{
			Path:     "/partial-register",
			RawQuery: q.Encode(),
		}

		ctx.Redirect(http.StatusSeeOther, partialRegisterURL.RequestURI())
	}

	sessions[sessionID] = internalUserID

	ctx.Redirect(http.StatusSeeOther, "/")
}
