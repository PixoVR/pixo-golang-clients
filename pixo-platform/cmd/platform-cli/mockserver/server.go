package mockserver

import (
	"encoding/json"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Run(mode string, endpoint string, mockResponse []byte) {
	gin.SetMode(mode)
	router := gin.Default()
	_ = router.SetTrustedProxies(nil)

	router.GET(endpoint, func(c *gin.Context) {
		requestHandler(c.Writer, c.Request, mockResponse)
	})

	port := input.GetConfigValue("server-port", "SERVER_PORT")
	log.Info().Msgf("Starting mock server serving endpoint %s on port %s", endpoint, port)

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

	var receivedData Request
	if err = conn.ReadJSON(&receivedData); err != nil {
		log.Error().Err(err).Msg("unable to read message")
		return
	}

	log.Info().Msg("Received request:")
	log.Info().Object("received", receivedData)

	var responseJSON Request
	_ = json.Unmarshal(response, &responseJSON)

	log.Info().Msg("Sending response:")
	log.Info().Object("response", responseJSON)

	if err = conn.WriteMessage(websocket.TextMessage, response); err != nil {
		log.Error().Err(err).Msg("unable to write message")
		return
	}
}

type Request map[string]interface{}

func (r Request) MarshalZerologObject(e *zerolog.Event) {
	printMap("", r, 0)
}

func printMap(key string, m map[string]interface{}, depth int) {
	bracketIndent := strings.Repeat("\t", depth)
	itemIndent := strings.Repeat("\t", depth+1)

	if key != "" {
		log.Info().Msgf("%s%s: {", bracketIndent, key)
	} else {
		log.Info().Msgf("%s{", bracketIndent)
	}

	for k, v := range m {
		if val, ok := v.(map[string]interface{}); ok {
			printMap(k, val, depth+1)
		} else {
			log.Info().Msgf("%s%s: %v", itemIndent, k, v)
		}
	}
	log.Info().Msgf("%s}", bracketIndent)
}
