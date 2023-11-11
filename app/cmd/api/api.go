package api

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/arek-e/D7024E/app/internal"
	"github.com/gin-gonic/gin"
)

var PORT = 2337

type API struct {
	Node *internal.Kademlia
	Net  *internal.Network
}

type GetResponse struct {
	Data    string           `json:"data"`
	Contact internal.Contact `json:"contact"`
}

func (api *API) StartAPI(address string, exitCh chan<- struct{}) {
	fmt.Println("\n======Kadlab node API========")

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	objectsGroup := router.Group("/objects")
	{
		objectsGroup.GET("/:hash", api.GetData)
		objectsGroup.POST("", api.StoreData)
	}

	// apiPort := os.Getenv("API_PORT")
	// if apiPort == "" {
	// 	// Handle the case where the environment variable is not set
	// 	fmt.Println("API_PORT environment variable is not set")
	// 	return
	// }
	apiPort := "2337"

	ip := fmt.Sprintf("%s:%s", address, apiPort)
	fmt.Printf("Server is running at: %s\n", ip)
	err := router.Run(ip)
	if err != nil {
		fmt.Println("error when listening to the http server " + err.Error())
	}
}

func (api *API) StoreData(ctx *gin.Context) {
	// Read data from the request body
	var requestBody struct {
		Data string `json:"data"`
	}

	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Store the data
	hash := api.Net.Node.Store([]byte(requestBody.Data))

	// Set the Location header
	locationHeader := "/objects/" + hash
	ctx.Header("Location", locationHeader)

	// Respond with 201 CREATED
	ctx.IndentedJSON(http.StatusCreated, gin.H{"key": hash})
}

func (api *API) GetData(ctx *gin.Context) {
	hash := ctx.Param("hash")

	// Lookup the data and contact based on the hash
	_, data, contact := api.Net.Node.Lookup(hash)

	// If data is not found, return a 404 Not Found response
	if data == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Object not found"})
		return
	}

	// Convert the bytes to a string
	dataStr := base64.StdEncoding.EncodeToString(data)

	// Create the response structure
	res := GetResponse{
		Data:    dataStr,
		Contact: contact,
	}

	// Respond with the contents of the object and contact information
	ctx.JSON(http.StatusOK, res)
}
