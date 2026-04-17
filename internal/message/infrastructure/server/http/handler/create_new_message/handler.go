package createnewmessage

import (
	"github.com/Arheon/markus-backend/pkg/outbox"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	outbox *outbox.Publisher
}

func NewHandler(publisher *outbox.Publisher) *Handler {
	return &Handler{
		outbox: publisher,
	}
}

func (h *Handler) Handle(ctx *gin.Context) {

}
