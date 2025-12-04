package controllers

import (
	"database/sql"
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

// CREATE: Menambahkan bioskop baru
func CreateBioskop(c *gin.Context) {
	var in createBioskopReq
	if err := c.ShouldBindJSON(&in); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validasi input
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

// READ ALL: Mendapatkan semua data bioskop
func GetAllBioskop(c *gin.Context) {
	const q = `SELECT id, nama, lokasi, rating FROM bioskop`

	rows, err := database.DB.Query(q)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "database error",
			"details": err.Error(),
		})
		return
	}
	defer rows.Close()

	var result []models.Bioskop

	for rows.Next() {
		var b models.Bioskop
		err = rows.Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   "scanning error",
				"details": err.Error(),
			})
			return
		}
		result = append(result, b)
	}

	// Jika kosong, kembalikan array kosong, bukan null
	if result == nil {
		result = []models.Bioskop{}
	}

	c.JSON(http.StatusOK, result)
}

// READ ONE: Mendapatkan detail bioskop berdasarkan ID
func GetBioskopByID(c *gin.Context) {
	id := c.Param("id")
	const q = `SELECT id, nama, lokasi, rating FROM bioskop WHERE id = $1`

	var b models.Bioskop
	err := database.DB.QueryRow(q, id).Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)

	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "bioskop not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "database error",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, b)
}

// UPDATE: Memperbarui data bioskop
func UpdateBioskop(c *gin.Context) {
	id := c.Param("id")
	var in createBioskopReq

	// Bind JSON body
	if err := c.ShouldBindJSON(&in); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validasi input
	in.Nama = strings.TrimSpace(in.Nama)
	in.Lokasi = strings.TrimSpace(in.Lokasi)
	if in.Nama == "" || in.Lokasi == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "nama dan lokasi tidak boleh kosong",
		})
		return
	}

	const q = `
		UPDATE bioskop 
		SET nama = $1, lokasi = $2, rating = $3 
		WHERE id = $4
	`

	res, err := database.DB.Exec(q, in.Nama, in.Lokasi, in.Rating, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "database error",
			"details": err.Error(),
		})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "bioskop not found or no changes made",
		})
		return
	}

	// Mengembalikan data yang sudah diupdate (tanpa query ulang untuk efisiensi)
	c.JSON(http.StatusOK, gin.H{
		"message": "update success",
		"data": gin.H{
			"id":     id,
			"nama":   in.Nama,
			"lokasi": in.Lokasi,
			"rating": in.Rating,
		},
	})
}

// DELETE: Menghapus bioskop berdasarkan ID
func DeleteBioskop(c *gin.Context) {
	id := c.Param("id")
	const q = `DELETE FROM bioskop WHERE id = $1`

	res, err := database.DB.Exec(q, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "database error",
			"details": err.Error(),
		})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "bioskop not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "delete success",
		"id":      id,
	})
}
