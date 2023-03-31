package payment

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/grvsahil/golang-kafka/kafka-consumer/internal/log"
	"github.com/grvsahil/golang-kafka/kafka-consumer/internal/service/payment/model"
)

func RegisterHandlers(router *gin.Engine, svc model.Service, logger log.Logger) {
	res := resource{svc, logger}
	router.GET("/records", res.Records)
}

type resource struct {
	service model.Service
	log     log.Logger
}

func (res resource) Records(c *gin.Context) {
	records, err := res.service.Get(c.Request.Context(), res.log)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%v", err)})
		return
	}

	c.JSON(http.StatusOK, records)
}
