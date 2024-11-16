package controllers

import (
	"Serenesongserver/services"
	"Serenesongserver/utils"

	"net/http"

	"github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func SearchRouter(c *gin.Context) {
	keyword := c.Query("keyword")
	option  := c.Query("option")
	if keyword == "" {
		utils.HandleError(c, http.StatusBadRequest, "Keyword is required", nil)
		return
	}
	if option == "" {
		utils.HandleError(c, http.StatusBadRequest, "Option is required", nil)
		return
	}
	// fmt.Println("Option:", option)
	switch {
		case option == "cipai":  services.CipaiConsult(c, keyword)
		case option == "ci": 	 services.CiConsult(c, keyword)
		case option == "author": services.AuthorConsult(c, keyword)
		default:
			utils.HandleError(c, http.StatusBadRequest, "Invalid option", nil)
	}
}
