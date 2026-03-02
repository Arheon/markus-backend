package health

import (
	"github.com/Arheon/markus-backend/internal/shared/application/query/health"
	"github.com/Arheon/markus-backend/internal/shared/domain/errors"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/helpers"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/sirupsen/logrus"
)

type HealthHandler struct {
	injector *do.Injector
}

func NewHandler(injector *do.Injector) *HealthHandler {
	return &HealthHandler{
		injector: injector,
	}
}

func (h *HealthHandler) Handle(context *gin.Context) {
	logger, loggerOk := do.InvokeAs[*logrus.Logger](*h.injector)

	var bind urlbind
	if err := context.ShouldBindUri(&bind); err != nil {
		if loggerOk == nil {
			logger.WithError(err).Error("Unexpected url binding error")
		}

		helpers.AbortWithBadRequestErrorJSON(context)
		return
	}

	query := health.NewQuery(h.injector)
	result, err := query.Handle(bind.Type)
	if err != nil {
		if loggerOk == nil {
			logger.WithError(err).Error("Undefined query error")
		}

		helpers.AbortWithDomainErrorJSON(context, errors.NewError(500, "Internal error", err.Error()))
		return
	}

	helpers.SuccessJSONV2(context, result)
}
