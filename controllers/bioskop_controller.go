package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"bioskop-api/database"
	"bioskop-api/models"
)

type createBioskopReq struct {
	Nama   string  `json:"nama" binding:"required"`
	Lokasi string  `json:"lokasi" binding:"required"`
	Rating float32 `json:"rating"`
}

func CreateBioskop(c *gin.Context) {
	var in createBioskopReq
	if err := c.ShouldBindJSON(&in); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validasi Nama & Lokasi tidak boleh kosong/whitespace
	in.Nama = strings.TrimSpace(in.Nama)
	in.Lokasi = strings.TrimSpace(in.Lokasi)
	if in.Nama == "" || in.Lokasi == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "nama dan lokasi tidak boleh kosong",
		})
		return
	}

	const q = `
		INSERT INTO bioskop (nama, lokasi, rating)
		VALUES ($1, $2, $3)
		RETURNING id, nama, lokasi, rating;
	`

	var out models.Bioskop
	err := database.DB.QueryRow(q, in.Nama, in.Lokasi, in.Rating).
		Scan(&out.ID, &out.Nama, &out.Lokasi, &out.Rating)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "database error",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, out)
}
