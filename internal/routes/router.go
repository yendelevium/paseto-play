package routes

import (
	"github.com/gin-gonic/gin"
	pasetotokens "github.com/yendelevium/paseto-play/internal/auth"
)

func AddRoutes(superRoute *gin.RouterGroup, makerPublic *pasetotokens.PasetoPublicKeyPair, makerLocal *pasetotokens.PasetoLocalKey) {
	AsymmetricRoutes(superRoute, makerPublic)
	SymmetricRoutes(superRoute, makerLocal)
}
