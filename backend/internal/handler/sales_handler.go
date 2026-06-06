package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/water-factory/api/internal/dto"
	"github.com/water-factory/api/internal/service"
)

type SalesHandler struct {
	svc *service.SalesService
}

func NewSalesHandler(svc *service.SalesService) *SalesHandler {
	return &SalesHandler{svc: svc}
}

func userIDFromCtx(c *gin.Context) string {
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// ── Vehicles ──────────────────────────────────────────────────────────────────

// ListVehicles godoc
// @Summary      List active vehicles
// @Tags         sales
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Router       /sales/vehicles [get]
func (h *SalesHandler) ListVehicles(c *gin.Context) {
	items, err := h.svc.ListVehicles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// ── Sales Orders ──────────────────────────────────────────────────────────────

// ListSalesOrders godoc
// @Summary      List sales orders
// @Tags         sales
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Router       /sales/orders [get]
func (h *SalesHandler) ListSalesOrders(c *gin.Context) {
	items, err := h.svc.ListSalesOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetSalesOrder godoc
// @Summary      Get sales order by ID
// @Tags         sales
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Sales Order ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /sales/orders/{id} [get]
func (h *SalesHandler) GetSalesOrder(c *gin.Context) {
	id := c.Param("id")
	item, err := h.svc.GetSalesOrder(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

// CreateSalesOrder godoc
// @Summary      Create sales order
// @Tags         sales
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  dto.CreateSalesOrderRequest  true  "Request body"
// @Success      201  {object}  map[string]interface{}
// @Router       /sales/orders [post]
func (h *SalesHandler) CreateSalesOrder(c *gin.Context) {
	var req dto.CreateSalesOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateSalesOrder(c.Request.Context(), userIDFromCtx(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": item})
}

// UpdateSalesOrder godoc
// @Summary      Update sales order (draft only)
// @Tags         sales
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                       true  "Sales Order ID"
// @Param        body  body  dto.UpdateSalesOrderRequest  true  "Request body"
// @Success      200  {object}  map[string]interface{}
// @Router       /sales/orders/{id} [put]
func (h *SalesHandler) UpdateSalesOrder(c *gin.Context) {
	var req dto.UpdateSalesOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.UpdateSalesOrder(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "only draft") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

// ConfirmSalesOrder godoc
// @Summary      Confirm sales order
// @Tags         sales
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Sales Order ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /sales/orders/{id}/confirm [post]
func (h *SalesHandler) ConfirmSalesOrder(c *gin.Context) {
	item, err := h.svc.ConfirmSalesOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "conflict") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

// CancelSalesOrder godoc
// @Summary      Cancel sales order
// @Tags         sales
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Sales Order ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /sales/orders/{id}/cancel [post]
func (h *SalesHandler) CancelSalesOrder(c *gin.Context) {
	item, err := h.svc.CancelSalesOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "conflict") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

// ── Delivery Orders ───────────────────────────────────────────────────────────

// ListDeliveryOrders godoc
// @Summary      List delivery orders
// @Tags         sales
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Router       /sales/delivery-orders [get]
func (h *SalesHandler) ListDeliveryOrders(c *gin.Context) {
	items, err := h.svc.ListDeliveryOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetDeliveryOrder godoc
// @Summary      Get delivery order by ID
// @Tags         sales
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Delivery Order ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /sales/delivery-orders/{id} [get]
func (h *SalesHandler) GetDeliveryOrder(c *gin.Context) {
	item, err := h.svc.GetDeliveryOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

// CreateDeliveryOrder godoc
// @Summary      Create delivery order (FEFO auto-pick)
// @Tags         sales
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  dto.CreateDeliveryOrderRequest  true  "Request body"
// @Success      201  {object}  map[string]interface{}
// @Router       /sales/delivery-orders [post]
func (h *SalesHandler) CreateDeliveryOrder(c *gin.Context) {
	var req dto.CreateDeliveryOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateDeliveryOrder(c.Request.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "conflict") || strings.Contains(err.Error(), "insufficient") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": item})
}

// DispatchDeliveryOrder godoc
// @Summary      Dispatch delivery order (deducts FG stock)
// @Tags         sales
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Delivery Order ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /sales/delivery-orders/{id}/dispatch [post]
func (h *SalesHandler) DispatchDeliveryOrder(c *gin.Context) {
	item, err := h.svc.DispatchDeliveryOrder(c.Request.Context(), c.Param("id"), userIDFromCtx(c))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "conflict") || strings.Contains(err.Error(), "insufficient") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

// DeliverDeliveryOrder godoc
// @Summary      Mark delivery order as delivered
// @Tags         sales
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Delivery Order ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /sales/delivery-orders/{id}/deliver [post]
func (h *SalesHandler) DeliverDeliveryOrder(c *gin.Context) {
	item, err := h.svc.DeliverDeliveryOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "conflict") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

// ── Invoices ──────────────────────────────────────────────────────────────────

// ListInvoices godoc
// @Summary      List invoices
// @Tags         sales
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Router       /sales/invoices [get]
func (h *SalesHandler) ListInvoices(c *gin.Context) {
	items, err := h.svc.ListInvoices(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetInvoice godoc
// @Summary      Get invoice by ID
// @Tags         sales
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Invoice ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /sales/invoices/{id} [get]
func (h *SalesHandler) GetInvoice(c *gin.Context) {
	item, err := h.svc.GetInvoice(c.Request.Context(), c.Param("id"))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

// CreateInvoice godoc
// @Summary      Create invoice from sales order
// @Tags         sales
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  dto.CreateInvoiceRequest  true  "Request body"
// @Success      201  {object}  map[string]interface{}
// @Router       /sales/invoices [post]
func (h *SalesHandler) CreateInvoice(c *gin.Context) {
	var req dto.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateInvoice(c.Request.Context(), userIDFromCtx(c), &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "conflict") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": item})
}

// MarkInvoicePaid godoc
// @Summary      Mark invoice as paid
// @Tags         sales
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string               true  "Invoice ID"
// @Param        body  body  dto.MarkPaidRequest  true  "Request body"
// @Success      200  {object}  map[string]interface{}
// @Router       /sales/invoices/{id}/pay [post]
func (h *SalesHandler) MarkInvoicePaid(c *gin.Context) {
	var req dto.MarkPaidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.MarkInvoicePaid(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "conflict") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}
