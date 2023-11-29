package mockserver

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Run(endpoint string, mockResponse []byte) {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	//router.SetTrustedProxies(nil)

	router.GET(endpoint, func(c *gin.Context) {
		requestHandler(c.Writer, c.Request, mockResponse)
	})

	port := input.GetConfigValue("server-port", "SERVER_PORT")
	log.Info().Msgf("Starting mock server on port %s", port)

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Error().Err(err).Msg("failed to run mock server")
		return
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request, response []byte) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to upgrade connection")
		return
	}
	defer conn.Close()

	var receivedData map[string]interface{}
	if err = conn.ReadJSON(&receivedData); err != nil {
		log.Error().Err(err).Msg("unable to read message")
		return
	}

	log.Info().Msgf("Received data: %+v", receivedData)

	log.Info().Msgf("Response to be sent: %v", string(response))

	if err = conn.WriteMessage(websocket.TextMessage, response); err != nil {
		log.Error().Err(err).Msg("unable to write message")
		return
	}
}
