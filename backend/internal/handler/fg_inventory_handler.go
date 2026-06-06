package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/water-factory/api/internal/dto"
	"github.com/water-factory/api/internal/service"
)

type FGInventoryHandler struct {
	svc *service.FGInventoryService
}

func NewFGInventoryHandler(svc *service.FGInventoryService) *FGInventoryHandler {
	return &FGInventoryHandler{svc: svc}
}

// ── Stock ─────────────────────────────────────────────────────────────────────

// ListFGStockSummary godoc
// @Summary      สรุปสต็อกสินค้าสำเร็จรูป (grouped by FG)
// @Tags         inventory-fg
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]domain.FGStockSummary
// @Router       /inventory/finished-goods/stock/summary [get]
func (h *FGInventoryHandler) ListStockSummary(c *gin.Context) {
	items, err := h.svc.ListStockSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// ListFGStockLots godoc
// @Summary      รายการ lots สินค้าสำเร็จรูป (filter ตาม finished_good_id ได้)
// @Tags         inventory-fg
// @Produce      json
// @Security     BearerAuth
// @Param        finished_good_id  query  string  false  "Filter by finished good ID"
// @Success      200  {object}  map[string][]domain.FGStockLot
// @Router       /inventory/finished-goods/stock/lots [get]
func (h *FGInventoryHandler) ListStockLots(c *gin.Context) {
	fgID := c.Query("finished_good_id")
	items, err := h.svc.ListStockLots(c.Request.Context(), fgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetFGStockLotDetail godoc
// @Summary      ดูรายละเอียด lot สินค้าสำเร็จรูป พร้อม movement history
// @Tags         inventory-fg
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "FG Stock Lot ID"
// @Success      200  {object}  domain.FGStockLotDetail
// @Router       /inventory/finished-goods/stock/lots/{id} [get]
func (h *FGInventoryHandler) GetStockLotDetail(c *gin.Context) {
	detail, err := h.svc.GetStockLotDetail(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "stock lot not found"})
		return
	}
	c.JSON(http.StatusOK, detail)
}

// ── Adjustments ───────────────────────────────────────────────────────────────

// ListFGAdjustments godoc
// @Summary      รายการปรับปรุงสต็อกสินค้าสำเร็จรูป
// @Tags         inventory-fg
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]domain.FGAdjustment
// @Router       /inventory/finished-goods/adjustments [get]
func (h *FGInventoryHandler) ListAdjustments(c *gin.Context) {
	items, err := h.svc.ListAdjustments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetFGAdjustment godoc
// @Summary      ดูใบปรับปรุงสต็อกสินค้าสำเร็จรูปตาม ID
// @Tags         inventory-fg
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "FG Adjustment ID"
// @Success      200  {object}  domain.FGAdjustment
// @Router       /inventory/finished-goods/adjustments/{id} [get]
func (h *FGInventoryHandler) GetAdjustment(c *gin.Context) {
	adj, err := h.svc.GetAdjustment(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "adjustment not found"})
		return
	}
	c.JSON(http.StatusOK, adj)
}

// CreateFGAdjustment godoc
// @Summary      สร้างใบปรับปรุงสต็อกสินค้าสำเร็จรูป (pending)
// @Tags         inventory-fg
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  dto.CreateFGAdjustmentRequest  true  "FG Adjustment request"
// @Success      201  {object}  domain.FGAdjustment
// @Router       /inventory/finished-goods/adjustments [post]
func (h *FGInventoryHandler) CreateAdjustment(c *gin.Context) {
	var req dto.CreateFGAdjustmentRequest
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

// ApproveFGAdjustment godoc
// @Summary      อนุมัติใบปรับปรุงสต็อกสินค้าสำเร็จรูป (อัปเดต stock lots)
// @Tags         inventory-fg
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "FG Adjustment ID"
// @Success      200  {object}  domain.FGAdjustment
// @Router       /inventory/finished-goods/adjustments/{id}/approve [post]
func (h *FGInventoryHandler) ApproveAdjustment(c *gin.Context) {
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

// CancelFGAdjustment godoc
// @Summary      ยกเลิกใบปรับปรุงสต็อกสินค้าสำเร็จรูป (pending เท่านั้น)
// @Tags         inventory-fg
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "FG Adjustment ID"
// @Success      200  {object}  map[string]string
// @Router       /inventory/finished-goods/adjustments/{id}/cancel [post]
func (h *FGInventoryHandler) CancelAdjustment(c *gin.Context) {
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
