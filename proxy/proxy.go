package proxy

import (
	"fmt"
	"gateway/common"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

var jsonContentType = []string{"application/json; charset=utf-8"}

type resp struct {
	data []byte
}

func (r *resp) Render(w http.ResponseWriter) error {
	writeContentType(w, jsonContentType)
	_, err := w.Write(r.data)
	return err
}

func (r *resp) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

func Start() {
	eng := gin.Default()
	grpcCli := NewClient()
	eng.POST("/", func(context *gin.Context) {
		path := fmt.Sprintf(
			"%s.%s/%s",
			context.Request.Header.Get("Package"),
			context.Request.Header.Get("Service"),
			context.Request.Header.Get("Method"),
		)

		data, err := ioutil.ReadAll(context.Request.Body)
		if err != nil {
			fmt.Println(err)
			context.String(http.StatusBadGateway, err.Error())
			return
		}

		r, err := grpcCli.Post(common.GetGrpcAddress(), path, data)
		if err != nil {
			fmt.Println(err)
			context.String(http.StatusBadGateway, err.Error())
			return
		}

		context.Render(http.StatusOK, &resp{data: r})
	})

	eng.Run(":8080")
}
