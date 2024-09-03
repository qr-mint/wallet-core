package middlewares

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"nexus-wallet/api/response"
	"time"
)

type RateLimiter struct {
	responseFactory *response.ResponseFactory
	params          Params
}

type Params struct {
	MaxRequests float64
}

func NewRateLimiter(params Params, responseFactory *response.ResponseFactory) *RateLimiter {
	return &RateLimiter{
		params:          params,
		responseFactory: responseFactory,
	}
}

func (l RateLimiter) CheckLimit() gin.HandlerFunc {
	return func(context *gin.Context) {
		lmt := tollbooth.NewLimiter(l.params.MaxRequests, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
		lmt.SetBurst(int(l.params.MaxRequests))
		lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})

		err := tollbooth.LimitByRequest(lmt, context.Writer, context.Request)
		if err != nil {
			l.responseFactory.TooManyRequests(context)
			context.Abort()
			return

		}
		context.Next()
	}
}
