package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/water-factory/api/internal/dto"
	"github.com/water-factory/api/internal/service"
)

type RMInventoryHandler struct {
	svc *service.RMInventoryService
}

func NewRMInventoryHandler(svc *service.RMInventoryService) *RMInventoryHandler {
	return &RMInventoryHandler{svc: svc}
}

func isStatusConflict(err error) bool {
	return strings.Contains(err.Error(), "is not in") ||
		strings.Contains(err.Error(), "can be cancelled")
}

// ── Warehouse Locations ───────────────────────────────────────────────────────

// ListWarehouseLocations godoc
// @Summary      รายชื่อตำแหน่งคลัง
// @Tags         inventory-raw-material
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]domain.WarehouseLocation
// @Router       /inventory/raw-material/warehouse-locations [get]
func (h *RMInventoryHandler) ListWarehouseLocations(c *gin.Context) {
	items, err := h.svc.ListWarehouseLocations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// ── GRN ──────────────────────────────────────────────────────────────────────

// ListGRNs godoc
// @Summary      รายการใบรับวัตถุดิบ (GRN)
// @Tags         inventory-raw-material
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]domain.GRN
// @Router       /inventory/raw-material/grn [get]
func (h *RMInventoryHandler) ListGRNs(c *gin.Context) {
	items, err := h.svc.ListGRNs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetGRN godoc
// @Summary      ดึงใบรับวัตถุดิบตาม ID
// @Tags         inventory-raw-material
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "GRN ID"
// @Success      200  {object}  domain.GRN
// @Router       /inventory/raw-material/grn/{id} [get]
func (h *RMInventoryHandler) GetGRN(c *gin.Context) {
	grn, err := h.svc.GetGRN(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "GRN not found"})
		return
	}
	c.JSON(http.StatusOK, grn)
}

// CreateGRN godoc
// @Summary      สร้างใบรับวัตถุดิบ (draft)
// @Tags         inventory-raw-material
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  dto.CreateGRNRequest  true  "GRN request"
// @Success      201  {object}  domain.GRN
// @Router       /inventory/raw-material/grn [post]
func (h *RMInventoryHandler) CreateGRN(c *gin.Context) {
	var req dto.CreateGRNRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	grn, err := h.svc.CreateGRN(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, grn)
}

// UpdateGRN godoc
// @Summary      แก้ไขใบรับวัตถุดิบ (draft เท่านั้น)
// @Tags         inventory-raw-material
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string              true  "GRN ID"
// @Param        body  body  dto.UpdateGRNRequest  true  "GRN request"
// @Success      200  {object}  domain.GRN
// @Router       /inventory/raw-material/grn/{id} [put]
func (h *RMInventoryHandler) UpdateGRN(c *gin.Context) {
	var req dto.UpdateGRNRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	grn, err := h.svc.UpdateGRN(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		if isStatusConflict(err) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, grn)
}

// ConfirmGRN godoc
// @Summary      ยืนยันการรับวัตถุดิบ (สร้าง stock lots)
// @Tags         inventory-raw-material
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "GRN ID"
// @Success      200  {object}  domain.GRN
// @Router       /inventory/raw-material/grn/{id}/confirm [post]
func (h *RMInventoryHandler) ConfirmGRN(c *gin.Context) {
	userID := c.GetString("user_id")
	grn, err := h.svc.ConfirmGRN(c.Request.Context(), c.Param("id"), userID)
	if err != nil {
		if isStatusConflict(err) || strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, grn)
}

// CancelGRN godoc
// @Summary      ยกเลิกใบรับวัตถุดิบ (draft เท่านั้น)
// @Tags         inventory-raw-material
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "GRN ID"
// @Success      200  {object}  map[string]string
// @Router       /inventory/raw-material/grn/{id}/cancel [post]
func (h *RMInventoryHandler) CancelGRN(c *gin.Context) {
	err := h.svc.CancelGRN(c.Request.Context(), c.Param("id"))
	if err != nil {
		if isStatusConflict(err) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}

// ── Stock ─────────────────────────────────────────────────────────────────────

// ListStockSummary godoc
// @Summary      สรุปสต็อกวัตถุดิบปัจจุบัน (grouped by RM)
// @Tags         inventory-raw-material
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]domain.StockSummary
// @Router       /inventory/raw-material/stock/summary [get]
func (h *RMInventoryHandler) ListStockSummary(c *gin.Context) {
	items, err := h.svc.ListStockSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// ListStockLots godoc
// @Summary      รายการ stock lots (ทั้งหมด หรือ filter ตาม raw_material_id)
// @Tags         inventory-raw-material
// @Produce      json
// @Security     BearerAuth
// @Param        raw_material_id  query  string  false  "Filter by raw material ID"
// @Success      200  {object}  map[string][]domain.StockLot
// @Router       /inventory/raw-material/stock/lots [get]
func (h *RMInventoryHandler) ListStockLots(c *gin.Context) {
	rmID := c.Query("raw_material_id")
	items, err := h.svc.ListStockLots(c.Request.Context(), rmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetStockLotDetail godoc
// @Summary      ดูรายละเอียด stock lot พร้อม movement history
// @Tags         inventory-raw-material
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Stock Lot ID"
// @Success      200  {object}  domain.StockLotDetail
// @Router       /inventory/raw-material/stock/lots/{id} [get]
func (h *RMInventoryHandler) GetStockLotDetail(c *gin.Context) {
	detail, err := h.svc.GetStockLotDetail(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "stock lot not found"})
		return
	}
	c.JSON(http.StatusOK, detail)
}

// ── Adjustments ───────────────────────────────────────────────────────────────

// ListAdjustments godoc
// @Summary      รายการปรับปรุงสต็อก
// @Tags         inventory-raw-material
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]domain.Adjustment
// @Router       /inventory/raw-material/adjustments [get]
func (h *RMInventoryHandler) ListAdjustments(c *gin.Context) {
	items, err := h.svc.ListAdjustments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetAdjustment godoc
// @Summary      ดูใบปรับปรุงสต็อกตาม ID
// @Tags         inventory-raw-material
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Adjustment ID"
// @Success      200  {object}  domain.Adjustment
// @Router       /inventory/raw-material/adjustments/{id} [get]
func (h *RMInventoryHandler) GetAdjustment(c *gin.Context) {
	adj, err := h.svc.GetAdjustment(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "adjustment not found"})
		return
	}
	c.JSON(http.StatusOK, adj)
}

// CreateAdjustment godoc
// @Summary      สร้างใบปรับปรุงสต็อก (pending)
// @Tags         inventory-raw-material
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  dto.CreateAdjustmentRequest  true  "Adjustment request"
// @Success      201  {object}  domain.Adjustment
// @Router       /inventory/raw-material/adjustments [post]
func (h *RMInventoryHandler) CreateAdjustment(c *gin.Context) {
	var req dto.CreateAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	adj, err := h.svc.CreateAdjustment(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, adj)
}

// ApproveAdjustment godoc
// @Summary      อนุมัติใบปรับปรุงสต็อก (อัปเดต stock lots)
// @Tags         inventory-raw-material
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Adjustment ID"
// @Success      200  {object}  domain.Adjustment
// @Router       /inventory/raw-material/adjustments/{id}/approve [post]
func (h *RMInventoryHandler) ApproveAdjustment(c *gin.Context) {
	approverID := c.GetString("user_id")
	adj, err := h.svc.ApproveAdjustment(c.Request.Context(), c.Param("id"), approverID)
	if err != nil {
		if isStatusConflict(err) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, adj)
}

// CancelAdjustment godoc
// @Summary      ยกเลิกใบปรับปรุงสต็อก (pending เท่านั้น)
// @Tags         inventory-raw-material
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Adjustment ID"
// @Success      200  {object}  map[string]string
// @Router       /inventory/raw-material/adjustments/{id}/cancel [post]
func (h *RMInventoryHandler) CancelAdjustment(c *gin.Context) {
	err := h.svc.CancelAdjustment(c.Request.Context(), c.Param("id"))
	if err != nil {
		if isStatusConflict(err) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}
