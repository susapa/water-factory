package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/water-factory/api/internal/dto"
	"github.com/water-factory/api/internal/service"
)

type ProductionHandler struct {
	svc *service.ProductionService
}

func NewProductionHandler(svc *service.ProductionService) *ProductionHandler {
	return &ProductionHandler{svc: svc}
}

func isProdStatusConflict(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "is not in") ||
		strings.Contains(msg, "cannot be cancelled") ||
		strings.Contains(msg, "status")
}

// ── Orders ────────────────────────────────────────────────────────────────────

// ListOrders godoc
// @Summary      รายการใบสั่งผลิต
// @Tags         production
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]domain.ProductionOrder
// @Router       /production/orders [get]
func (h *ProductionHandler) ListOrders(c *gin.Context) {
	items, err := h.svc.ListOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// CreateOrder godoc
// @Summary      สร้างใบสั่งผลิต
// @Tags         production
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CreateProductionOrderRequest  true  "Production order"
// @Success      201   {object}  domain.ProductionOrder
// @Router       /production/orders [post]
func (h *ProductionHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateProductionOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	o, err := h.svc.CreateOrder(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, o)
}

// GetOrder godoc
// @Summary      ดูใบสั่งผลิต
// @Tags         production
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Production order ID"
// @Success      200  {object}  domain.ProductionOrder
// @Router       /production/orders/{id} [get]
func (h *ProductionHandler) GetOrder(c *gin.Context) {
	o, err := h.svc.GetOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

// UpdateOrder godoc
// @Summary      แก้ไขใบสั่งผลิต (draft เท่านั้น)
// @Tags         production
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                            true  "Production order ID"
// @Param        body  body      dto.UpdateProductionOrderRequest  true  "Update data"
// @Success      200   {object}  domain.ProductionOrder
// @Router       /production/orders/{id} [put]
func (h *ProductionHandler) UpdateOrder(c *gin.Context) {
	var req dto.UpdateProductionOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	o, err := h.svc.UpdateOrder(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		if isProdStatusConflict(err) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

// ConfirmOrder godoc
// @Summary      ยืนยันใบสั่งผลิต (BOM explosion)
// @Tags         production
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Production order ID"
// @Success      200  {object}  domain.ProductionOrder
// @Router       /production/orders/{id}/confirm [post]
func (h *ProductionHandler) ConfirmOrder(c *gin.Context) {
	o, err := h.svc.ConfirmOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		if isProdStatusConflict(err) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

// IssueRM godoc
// @Summary      เบิกวัตถุดิบ (FIFO)
// @Tags         production
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string              true  "Production order ID"
// @Param        body  body      dto.IssueRMRequest  true  "Issue lines"
// @Success      200   {object}  domain.ProductionOrder
// @Router       /production/orders/{id}/issue-rm [post]
func (h *ProductionHandler) IssueRM(c *gin.Context) {
	var req dto.IssueRMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	o, err := h.svc.IssueRM(c.Request.Context(), c.Param("id"), userID, &req)
	if err != nil {
		if isProdStatusConflict(err) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

// StartProduction godoc
// @Summary      เริ่มผลิต
// @Tags         production
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Production order ID"
// @Success      200  {object}  domain.ProductionOrder
// @Router       /production/orders/{id}/start [post]
func (h *ProductionHandler) StartProduction(c *gin.Context) {
	o, err := h.svc.StartProduction(c.Request.Context(), c.Param("id"))
	if err != nil {
		if isProdStatusConflict(err) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

// RecordYield godoc
// @Summary      บันทึกผลผลิต
// @Tags         production
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                   true  "Production order ID"
// @Param        body  body      dto.RecordYieldRequest   true  "Yield data"
// @Success      201   {object}  domain.ProductionYield
// @Router       /production/orders/{id}/record-yield [post]
func (h *ProductionHandler) RecordYield(c *gin.Context) {
	var req dto.RecordYieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	y, err := h.svc.RecordYield(c.Request.Context(), c.Param("id"), userID, &req)
	if err != nil {
		if isProdStatusConflict(err) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, y)
}

// CompleteOrder godoc
// @Summary      เสร็จสิ้นการผลิต
// @Tags         production
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Production order ID"
// @Success      200  {object}  domain.ProductionOrder
// @Router       /production/orders/{id}/complete [post]
func (h *ProductionHandler) CompleteOrder(c *gin.Context) {
	o, err := h.svc.CompleteOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		if isProdStatusConflict(err) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

// CancelOrder godoc
// @Summary      ยกเลิกใบสั่งผลิต
// @Tags         production
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Production order ID"
// @Success      200  {object}  map[string]string
// @Router       /production/orders/{id}/cancel [post]
func (h *ProductionHandler) CancelOrder(c *gin.Context) {
	if err := h.svc.CancelOrder(c.Request.Context(), c.Param("id")); err != nil {
		if isProdStatusConflict(err) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}

// ── Yields ────────────────────────────────────────────────────────────────────

// ListYields godoc
// @Summary      รายการบันทึกผลผลิต
// @Tags         production
// @Produce      json
// @Security     BearerAuth
// @Param        order_id  query     string  false  "Filter by production order ID"
// @Success      200       {object}  map[string][]domain.ProductionYield
// @Router       /production/yields [get]
func (h *ProductionHandler) ListYields(c *gin.Context) {
	orderID := c.Query("order_id")
	items, err := h.svc.ListYields(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}
