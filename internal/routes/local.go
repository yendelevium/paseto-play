package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	pasetotokens "github.com/yendelevium/paseto-play/internal/auth"
)

func SymmetricRoutes(superRoute *gin.RouterGroup, maker *pasetotokens.PasetoLocalKey) {
	r := superRoute.Group("/local")
	{
		// Define a simple GET endpoint to encrypt a constant payload
		r.GET("/encrypt", func(c *gin.Context) {
			// Return JSON response
			generated_paseto, err := maker.CreateTokenEncrypted(
				struct {
					Name       string
					Enrollment int
				}{
					Name:       "A",
					Enrollment: 12345,
				})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "couldn't create paseto token",
				})
				return
			}
			c.SetCookie("paseto-local", generated_paseto, 3600, "/", "localhost", false, true)
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// Define another GET endpoint to decrypt and get the payload in the paseto token
		r.GET("/decrypt", func(c *gin.Context) {
			paseto_token, err := c.Cookie("paseto-local")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "couldn't get paseto token",
				})
				return
			}

			payload, err := maker.DecryptToken(paseto_token)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "couldn't parse paseto token",
				})
				return
			}

			var s any
			err = json.Unmarshal(payload, &s)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "couldn't parse paseto token",
				})
				return
			}
			log.Printf("Token verified, parsed token: %v \n", s)

			c.JSON(http.StatusOK, gin.H{
				"payload": s,
			})
		})

	}
}
