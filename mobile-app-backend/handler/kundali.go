package handler

import (
	"net/http"

	"mobile-app-backend/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SaveKundali godoc
// @Summary      Save kundali chart
// @Description  Save a kundali chart for the authenticated user
// @Tags         kundali
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      model.Kundli  true  "Kundali birth details"
// @Success      201      {object}  model.KundaliSave
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /save-kundali [post]
func (h *Handler) SaveKundali(c *gin.Context) {
	rawID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Missing user context",
		})
		return
	}

	userID := rawID.(uuid.UUID)

	var kundali model.Kundli
	if err := c.ShouldBindJSON(&kundali); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "Invalid kundli JSON: " + err.Error(),
		})
		return
	}

	meta := &model.KundaliSave{
		Name:       kundali.Name,
		Day:        kundali.Day,
		Month:      kundali.Month,
		Year:       kundali.Year,
		Hour:       kundali.Hour,
		Minute:     kundali.Minute,
		BirthPlace: kundali.BirthPlace,
		Latitude:   kundali.Latitude,
		Longitude:  kundali.Longitude,
		TimeZone:   kundali.TimeZone,
		IsFemale:   kundali.IsFemale,
	}

	saved := h.store.CreateKundaliForUser(userID, meta)
	if saved == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "store_error",
			"message": "Failed to create kundali record",
		})
		return
	}

	c.JSON(http.StatusCreated, saved)
}

// ListKundali godoc
// @Summary      List user kundalis
// @Description  Get all kundali charts for the authenticated user
// @Tags         kundali
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Router       /list-kundali [get]
func (h *Handler) ListKundali(c *gin.Context) {
	rawID, _ := c.Get("user_id")
	userID := rawID.(uuid.UUID)

	records := h.store.FindKundaliByUserID(userID)
	c.JSON(http.StatusOK, gin.H{
		"data":  records,
		"count": len(records),
	})
}

// ListAllKundali godoc
// @Summary      List all kundalis
// @Description  Get all kundali charts across all users
// @Tags         kundali
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Router       /all-kundali [get]
func (h *Handler) ListAllKundali(c *gin.Context) {
	records := h.store.AllKundali()
	c.JSON(http.StatusOK, gin.H{
		"data":  records,
		"count": len(records),
	})
}
