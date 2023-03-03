package jsonmask

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(
		Middleware(Options{
			Getter: func(ctx *gin.Context) (string, error) {
				return ctx.Query("jsonmask"), nil
			},
			ErrorHandler: func(ctx *gin.Context, err error) {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			},
		}),
	)

	router.GET("/users", func(ctx *gin.Context) {
		ctx.JSON(
			200,
			gin.H{
				"data": []gin.H{
					{
						"name": "xiao ming",
						"age":  12,
					},
					{
						"name": "xiao hong",
						"age":  11,
					},
				},
				"count": 2,
			},
		)
	})
	return router
}

func TestUsersRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()

	{
		req, _ := http.NewRequest("GET", "/users?jsonmask=(", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, `{"message":"invalid char before '('"}`, w.Body.String())
	}
}
