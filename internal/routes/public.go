package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	pasetotokens "github.com/yendelevium/paseto-play/internal/auth"
)

func AsymmetricRoutes(superRoute *gin.RouterGroup, maker *pasetotokens.PasetoPublicKeyPair) {
	r := superRoute.Group("/public")
	{
		// Define a simple GET endpoint to sign a constant payload
		r.GET("/sign", func(c *gin.Context) {
			// Return JSON response
			generated_paseto, err := maker.CreateToken(
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
			c.SetCookie("paseto", generated_paseto, 3600, "/", "localhost", false, true)
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// Define another GET endpoint to verify and get the payload in the paseto token
		r.GET("/verify", func(c *gin.Context) {
			paseto_token, err := c.Cookie("paseto")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "couldn't get paseto token",
				})
				return
			}

			payload, err := maker.VerifyToken(paseto_token)
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
