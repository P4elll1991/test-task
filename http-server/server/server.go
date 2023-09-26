package server

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type broker interface {
	SendChallenge(destination string) error
	SendMessage(destination string, data []byte) error
}

type server struct {
	port   string
	broker broker
}

func New(port string, broker broker) *server {
	return &server{port: port, broker: broker}
}

func (server *server) Run() {
	handler := gin.Default()

	handler.POST("/sendChallenge", server.sendChallenge)
	handler.POST("/sendMessage", server.sendMessage)

	httpServer := &http.Server{
		Addr:    ":8000",
		Handler: handler,
	}

	httpServer.ListenAndServe()
}

func (server *server) sendChallenge(c *gin.Context) {
	destination := c.Request.PostFormValue("destination")

	err := server.broker.SendChallenge(destination)
	if err != nil {
		log.Printf("sendChallenge: %s", err.Error())
		c.AbortWithStatusJSON(parseError(err))
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "OK"})
}

func (server *server) sendMessage(c *gin.Context) {
	data := []byte(c.Request.PostFormValue("data"))

	destination := c.Request.PostFormValue("destination")
	err := server.broker.SendMessage(destination, data)
	if err != nil {
		log.Printf("sendMessage: %s", err.Error())
		c.AbortWithStatusJSON(parseError(err))
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "OK"})

}

func parseError(err error) (int, map[string]string) {
	statusCode := http.StatusInternalServerError
	body := map[string]string{"details": err.Error(), "error": "Unknown error", "status": "ERROR"}
	if grpc.Code(err) == codes.Unavailable {
		statusCode = http.StatusBadRequest
		body["error"] = "there is no active gRPC server at destination"
	}
	return statusCode, body
}

func getData(c *gin.Context) ([]byte, error) {
	data, err := c.FormFile("data")
	if err != nil {
		return nil, err
	}

	file, err := data.Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return ioutil.ReadAll(file)
}
