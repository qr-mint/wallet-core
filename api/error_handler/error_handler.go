package error_handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/golib4/logger/logger"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"nexus-wallet/api/response"
	"nexus-wallet/internal/app_error"
	"nexus-wallet/pkg/sentry"
	"strings"
)

type Params struct {
	IsDebugMode bool
}

type HttpErrorHandler struct {
	responseFactory *response.ResponseFactory
	sentry          *sentry.Sentry
	logger          logger.Logger
	params          Params
}

func NewHttpErrorHandler(
	responseFactory *response.ResponseFactory,
	sentry *sentry.Sentry,
	logger logger.Logger,
	params Params,
) *HttpErrorHandler {
	return &HttpErrorHandler{
		responseFactory: responseFactory,
		sentry:          sentry,
		logger:          logger,
		params:          params,
	}
}

func (h *HttpErrorHandler) Handle(context *gin.Context, err *app_error.AppError) {
	switch err.Code {
	case app_error.InvalidData:
		h.responseFactory.BadRequest(context, h.formatErrorText(err.Error.Error()))
	case app_error.IllegalOperation:
		h.responseFactory.UnprocessableEntity(context, h.formatErrorText(err.Error.Error()))
	case app_error.ResourceNotFound:
		h.responseFactory.NotFound(context, h.formatErrorText(err.Error.Error()))
	case app_error.Unauthorized:
		h.responseFactory.Unauthorized(context)
	default:
		message := "Internal server error."

		logMessage := h.formatErrorText(fmt.Sprintf("%s at endpoint `%s`. %s", message, context.Request.URL.Path, err.Error))
		h.logger.Error(logMessage)

		if h.params.IsDebugMode {
			message = logMessage
		}

		h.sentry.HandleError(err.Error)
		h.responseFactory.InternalServerError(context, message)
	}
}

func (h *HttpErrorHandler) formatErrorText(message string) string {
	sentences := strings.Split(message, ". ")
	for i := range sentences {
		words := strings.Fields(sentences[i])
		words[0] = cases.Title(language.English, cases.Compact).String(words[0])

		sentences[i] = strings.Join(words, " ")
	}
	message = strings.Join(sentences, ". ")
	if !strings.HasSuffix(message, ".") {
		message += "."
	}

	return message
}
