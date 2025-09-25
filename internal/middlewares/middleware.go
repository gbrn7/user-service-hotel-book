package middlewares

import (
	"context"
	"net/http"
	"user-service/constants"
	errConstants "user-service/constants/error"
	"user-service/helpers"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandlePanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("recovered from panic: %v", r)

				c.JSON(http.StatusInternalServerError, helpers.Response{
					Status:  constants.Error,
					Message: errConstants.ErrInternalServerError.Error(),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if err != nil {
			logrus.Errorf("rate limit exceeded: %v", err)
			c.JSON(http.StatusTooManyRequests, helpers.Response{
				Status:  constants.Error,
				Message: errConstants.ErrTooManyRequests.Error(),
			})
			c.Abort()
		}
		c.Next()
	}
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		token := c.GetHeader(constants.Authorization)

		if token == "" {
			logrus.Errorf("missing authorization token")
			helpers.ResponseUnauthorized(c, errConstants.ErrUnautorized.Error())
			return
		}

		claimToken, err := helpers.ValidateBearerToken(c, token)
		if err != nil {
			logrus.Errorf("invalid token: %v", err)
			helpers.ResponseUnauthorized(c, err.Error())
			return
		}

		err = helpers.ValidateAPIKey(c)
		if err != nil {
			logrus.Errorf("invalid API key: %v", err)
			helpers.ResponseUnauthorized(c, err.Error())
			return
		}

		userLogin := c.Request.WithContext(context.WithValue(c.Request.Context(), constants.UserLogin, claimToken.User))
		c.Request = userLogin

		c.Set(constants.Token, token)

		c.Next()
	}
}
