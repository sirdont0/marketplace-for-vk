package handlers

import (
	"net/http"
	"strconv"

	"marketplace/internal/models"
	"marketplace/internal/storage"

	"github.com/gin-gonic/gin"
)

func CreateAd(c *gin.Context) {
	var ad models.Ad
	if err := c.ShouldBindJSON(&ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.MustGet("userID").(string)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}
	ad.AuthorID = userID

	// Валидация данных
	if len(ad.Title) < 5 || len(ad.Title) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title must be between 5 and 100 characters"})
		return
	}

	if len(ad.Text) < 10 || len(ad.Text) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Text must be between 10 and 1000 characters"})
		return
	}

	if ad.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be positive"})
		return
	}

	if err := storage.CreateAd(&ad); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create ad"})
		return
	}

	c.JSON(http.StatusCreated, ad)
}

func GetAds(c *gin.Context) {
	// Параметры запроса
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sort", "created_at_desc")
	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

	// Определение сортировки
	var orderBy string
	switch sortBy {
	case "price_asc":
		orderBy = "price ASC"
	case "price_desc":
		orderBy = "price DESC"
	default:
		orderBy = "created_at DESC"
	}

	// Получение объявлений
	ads, err := storage.GetAds(page, limit, orderBy, minPrice, maxPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get ads"})
		return
	}

	// Добавление признака "моё" для авторизованных пользователей
	userID, exists := c.Get("userID")
	response := make([]models.AdResponse, len(ads))
	for i, ad := range ads {
		response[i] = models.AdResponse{
			ID:        ad.ID,
			Title:     ad.Title,
			Text:      ad.Text,
			ImageURL:  ad.ImageURL,
			Price:     ad.Price,
			Author:    ad.Author,
			CreatedAt: ad.CreatedAt,
			IsMine:    exists && ad.AuthorID == userID.(int),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": len(response),
		},
	})
}
