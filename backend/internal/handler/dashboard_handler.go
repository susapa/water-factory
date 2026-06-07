package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/water-factory/api/internal/service"
)

type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

// GetSummary godoc
// @Summary      Dashboard KPI summary
// @Description  Returns aggregated KPI counts and chart data for the dashboard
// @Tags         dashboard
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /dashboard/summary [get]
func (h *DashboardHandler) GetSummary(c *gin.Context) {
	data, err := h.svc.GetSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
