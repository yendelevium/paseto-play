package routes

import (
	"github.com/gin-gonic/gin"
	pasetotokens "github.com/yendelevium/paseto-play/internal/auth"
)

func AddRoutes(superRoute *gin.RouterGroup, maker *pasetotokens.PasetoPublicKeyPair) {
	AsymmetricRoutes(superRoute, maker)
}
