package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"user-service/constants"
	errConstants "user-service/constants/error"
	"user-service/helpers/configs"

	"github.com/gin-gonic/gin"
)

func ValidateAPIKey(c *gin.Context) error {
	cfg := configs.Get()

	apiKey := c.GetHeader(constants.XApiKey)
	requestAt := c.GetHeader(constants.XRequestAt)
	serviceName := c.GetHeader(constants.XServiceName)
	signatureKey := cfg.Service.SignatureKey

	validateKey := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)

	hash := sha256.New()
	hash.Write([]byte(validateKey))
	resultHash := hex.EncodeToString(hash.Sum(nil))

	if apiKey != resultHash {
		return errConstants.ErrUnautorized
	}

	return nil
}
