# Gin Json Mask Middleware

## Example

```go
router := gin.Default()

router.Use(
  jsonmask.Middleware(jsonmask.Options{
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
```

```shell
curl 'http://localhost:3000/users?jsonmask=count,data(id)' | jq
```
