package car

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MarNawar/carZone/models"
	"github.com/MarNawar/carZone/service"
	"github.com/gin-gonic/gin"
)

type CarHandler struct {
	service service.CarServiceInterface
}

func NewCarHandler(service service.CarServiceInterface) *CarHandler {
	return &CarHandler{
		service: service,
	}
}

func (h *CarHandler) HandleGetCarByID(c *gin.Context) {
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

	res, err := h.service.GetCarById(ctx, id)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching the car item"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *CarHandler) HandleGetCarByBrand(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	brand := c.Query("brand")
	if brand == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"message": "please provide the valid brand",
		})
		return
	}

	isEngineStr := c.Query("isEngine")
	isEngine, err := strconv.ParseBool(isEngineStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"message": "please provide the valid isEngine",
		})
		return
	}

	res, err := h.service.GetCarsByBrand(ctx, brand, isEngine)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}


func (h *CarHandler) HandleCreateCar(c *gin.Context){
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var carReq *models.CarRequest
	if err := c.BindJSON(&carReq); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.CreateCar(ctx, carReq)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *CarHandler) HandleUpdateCar(c *gin.Context){
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var carReq *models.CarRequest
	if err := c.BindJSON(&carReq); err != nil{
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
	
	res, err := h.service.UpdateCar(ctx, id, carReq)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *CarHandler) HandleDeleteCar(c *gin.Context){
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
	
	res, err := h.service.DeleteCar(ctx, id)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
