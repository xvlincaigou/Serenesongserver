package controllers

import (
	"Serenesongserver/services"
	"Serenesongserver/utils"

	"net/http"

	"github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func SearchRouter(c *gin.Context) {
	keywords := c.QueryArray("keyword")
	option   := c.Query("option")
	if len(keywords) == 0 {
		utils.HandleError(c, http.StatusBadRequest, "Keyword is required", nil)
		return
	}
	if option == "" {
		utils.HandleError(c, http.StatusBadRequest, "Option is required", nil)
		return
	}
	// fmt.Println("Option:", option)
	switch {
		case option == "cipai":  services.CipaiConsult(c, keywords)
		case option == "ci": 	 services.CiConsult(c, keywords)
		case option == "author": services.AuthorConsult(c, keywords)
		default:
			utils.HandleError(c, http.StatusBadRequest, "Invalid option", nil)
	}
}
