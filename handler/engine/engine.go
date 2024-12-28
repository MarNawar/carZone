package engine

import (
	"context"
	"net/http"
	"time"

	"github.com/MarNawar/carZone/models"
	"github.com/MarNawar/carZone/service"
	"github.com/gin-gonic/gin"
)

type EngineHandler struct {
	service service.EngineServiceInterface
}

func NewEngineHandler(service service.EngineServiceInterface) *EngineHandler {
	return &EngineHandler{
		service: service,
	}
}

func (h *EngineHandler) HandleGetEngineByID(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"message": "please provide the valid id",
		})
		return
	}

	res, err := h.service.GetEngineByID(ctx, id)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching the engine"})
		return
	}
	c.JSON(http.StatusOK, res)
}


func (h *EngineHandler) HandleCreateEngine(c *gin.Context){
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var engineRequest *models.EngineRequest
	if err := c.BindJSON(&engineRequest); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.CreateEngine(ctx, engineRequest)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while creating the engine"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *EngineHandler) HandleUpdateEngine(c *gin.Context){
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var engineRequest *models.EngineRequest
	if err := c.BindJSON(&engineRequest); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"message": "please provide the valid id",
		})
		return
	}
	
	res, err := h.service.UpdateEngine(ctx, engineRequest, id)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while updating the engine"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *EngineHandler) HandleDeleteEngine(c *gin.Context){
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"message": "please provide the valid id",
		})
		return
	}
	
	res, err := h.service.DeleteEngine(ctx, id)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while deleting the engine item"})
		return
	}
	c.JSON(http.StatusOK, res)
}
