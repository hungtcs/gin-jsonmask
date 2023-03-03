package jsonmask

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	jsonmask "github.com/teambition/json-mask-go"
)

type responseWriter struct {
	gin.ResponseWriter
	buffer *bytes.Buffer
}

func (writer responseWriter) Write(bytes []byte) (int, error) {
	return writer.buffer.Write(bytes)
}

type Options struct {
	Getter       func(ctx *gin.Context) (string, error)
	ErrorHandler func(ctx *gin.Context, err error)
}

func Middleware(options Options) gin.HandlerFunc {
	if options.Getter == nil {
		options.Getter = func(ctx *gin.Context) (string, error) {
			return ctx.Query("jsonmask"), nil
		}
	}

	if options.ErrorHandler == nil {
		options.ErrorHandler = func(ctx *gin.Context, err error) {
			ctx.JSON(http.StatusBadRequest, err.Error())
		}
	}

	return func(ctx *gin.Context) {
		var err error
		var jsonMask string

		if jsonMask, err = options.Getter(ctx); err != nil {
			options.ErrorHandler(ctx, err)
			ctx.Abort()
			return
		}

		if jsonMask != "" {
			var selection jsonmask.Selection

			if selection, err = jsonmask.Compile(jsonMask); err != nil {
				options.ErrorHandler(ctx, err)
				ctx.Abort()
				return
			}

			writer := responseWriter{
				buffer:         &bytes.Buffer{},
				ResponseWriter: ctx.Writer,
			}
			oldWriter := ctx.Writer
			ctx.Writer = writer

			ctx.Next()

			if ctx.IsAborted() {
				ctx.Writer = oldWriter
				return
			}

			response, err := selection.Mask(writer.buffer.Bytes())
			if err != nil {
				options.ErrorHandler(ctx, err)
				ctx.Abort()
			} else {
				writer.ResponseWriter.Write(response)
			}
		} else {
			ctx.Next()
		}
	}
}
