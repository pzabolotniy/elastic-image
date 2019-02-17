package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

// GetInputData duplicates body stream
// it is used to log input body
func GetInputData(ctx *gin.Context) *string {
	var bodyContent string

	if ctx.Request.Body != nil {
		buf, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			// some error handling should be done here
		}
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

		bodyBuffer := new(bytes.Buffer)
		bodyBuffer.ReadFrom(rdr1)
		bodyContent = bodyBuffer.String()

		ctx.Request.Body = rdr2
	}

	return &bodyContent
}
