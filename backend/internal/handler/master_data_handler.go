package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/water-factory/api/internal/dto"
	"github.com/water-factory/api/internal/service"
)

type MasterDataHandler struct {
	svc *service.MasterDataService
}

func NewMasterDataHandler(svc *service.MasterDataService) *MasterDataHandler {
	return &MasterDataHandler{svc: svc}
}

func isDuplicateErr(err error) bool {
	return strings.Contains(err.Error(), "already exists")
}

// ── Lookups ───────────────────────────────────────────────────────────────────

// ListUOM godoc
// @Summary      รายชื่อหน่วยนับ (UOM)
// @Description  ดึง UOM ทั้งหมดสำหรับ dropdown
// @Tags         master-data
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]domain.UnitOfMeasure
// @Router       /master-data/uom [get]
func (h *MasterDataHandler) ListUOM(c *gin.Context) {
	items, err := h.svc.ListUOM(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// ListCategories godoc
// @Summary      รายชื่อหมวดหมู่วัตถุดิบ
// @Description  ดึงหมวดหมู่ทั้งหมดสำหรับ dropdown
// @Tags         master-data
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]domain.RawMaterialCategory
// @Router       /master-data/raw-material-categories [get]
func (h *MasterDataHandler) ListCategories(c *gin.Context) {
	items, err := h.svc.ListCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// ── Raw Materials ─────────────────────────────────────────────────────────────

// ListRawMaterials godoc
// @Summary      รายชื่อวัตถุดิบทั้งหมด
// @Tags         master-data
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]domain.RawMaterial
// @Router       /master-data/raw-materials [get]
func (h *MasterDataHandler) ListRawMaterials(c *gin.Context) {
	items, err := h.svc.ListRawMaterials(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetRawMaterial godoc
// @Summary      ดึงวัตถุดิบตาม ID
// @Tags         master-data
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Raw Material ID (UUID)"
// @Success      200  {object}  domain.RawMaterial
// @Failure      404  {object}  map[string]string
// @Router       /master-data/raw-materials/{id} [get]
func (h *MasterDataHandler) GetRawMaterial(c *gin.Context) {
	item, err := h.svc.GetRawMaterial(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateRawMaterial godoc
// @Summary      สร้างวัตถุดิบใหม่
// @Tags         master-data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CreateRawMaterialRequest  true  "ข้อมูลวัตถุดิบ"
// @Success      201   {object}  domain.RawMaterial
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Router       /master-data/raw-materials [post]
func (h *MasterDataHandler) CreateRawMaterial(c *gin.Context) {
	var req dto.CreateRawMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	item, err := h.svc.CreateRawMaterial(c.Request.Context(), &req)
	if err != nil {
		if isDuplicateErr(err) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateRawMaterial godoc
// @Summary      แก้ไขวัตถุดิบ
// @Tags         master-data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                        true  "Raw Material ID"
// @Param        body  body      dto.UpdateRawMaterialRequest  true  "ข้อมูลที่แก้ไข"
// @Success      200   {object}  domain.RawMaterial
// @Failure      400   {object}  map[string]string
// @Router       /master-data/raw-materials/{id} [put]
func (h *MasterDataHandler) UpdateRawMaterial(c *gin.Context) {
	var req dto.UpdateRawMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	item, err := h.svc.UpdateRawMaterial(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// SetRawMaterialActive godoc
// @Summary      เปิด/ปิดวัตถุดิบ
// @Tags         master-data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                     true  "Raw Material ID"
// @Param        body  body      object{is_active=boolean}  true  "สถานะ"
// @Success      200   {object}  map[string]string
// @Router       /master-data/raw-materials/{id}/active [put]
func (h *MasterDataHandler) SetRawMaterialActive(c *gin.Context) {
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := h.svc.SetRawMaterialActive(c.Request.Context(), c.Param("id"), req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// ── Finished Goods ────────────────────────────────────────────────────────────

// ListFinishedGoods godoc
// @Summary      รายชื่อสินค้าสำเร็จรูปทั้งหมด
// @Tags         master-data
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]domain.FinishedGood
// @Router       /master-data/finished-goods [get]
func (h *MasterDataHandler) ListFinishedGoods(c *gin.Context) {
	items, err := h.svc.ListFinishedGoods(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetFinishedGood godoc
// @Summary      ดึงสินค้าสำเร็จรูปตาม ID
// @Tags         master-data
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Finished Good ID (UUID)"
// @Success      200  {object}  domain.FinishedGood
// @Failure      404  {object}  map[string]string
// @Router       /master-data/finished-goods/{id} [get]
func (h *MasterDataHandler) GetFinishedGood(c *gin.Context) {
	item, err := h.svc.GetFinishedGood(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateFinishedGood godoc
// @Summary      สร้างสินค้าสำเร็จรูปใหม่
// @Tags         master-data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CreateFinishedGoodRequest  true  "ข้อมูลสินค้า"
// @Success      201   {object}  domain.FinishedGood
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Router       /master-data/finished-goods [post]
func (h *MasterDataHandler) CreateFinishedGood(c *gin.Context) {
	var req dto.CreateFinishedGoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	item, err := h.svc.CreateFinishedGood(c.Request.Context(), &req)
	if err != nil {
		if isDuplicateErr(err) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateFinishedGood godoc
// @Summary      แก้ไขสินค้าสำเร็จรูป
// @Tags         master-data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                         true  "Finished Good ID"
// @Param        body  body      dto.UpdateFinishedGoodRequest  true  "ข้อมูลที่แก้ไข"
// @Success      200   {object}  domain.FinishedGood
// @Failure      400   {object}  map[string]string
// @Router       /master-data/finished-goods/{id} [put]
func (h *MasterDataHandler) UpdateFinishedGood(c *gin.Context) {
	var req dto.UpdateFinishedGoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	item, err := h.svc.UpdateFinishedGood(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// SetFinishedGoodActive godoc
// @Summary      เปิด/ปิดสินค้าสำเร็จรูป
// @Tags         master-data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                     true  "Finished Good ID"
// @Param        body  body      object{is_active=boolean}  true  "สถานะ"
// @Success      200   {object}  map[string]string
// @Router       /master-data/finished-goods/{id}/active [put]
func (h *MasterDataHandler) SetFinishedGoodActive(c *gin.Context) {
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := h.svc.SetFinishedGoodActive(c.Request.Context(), c.Param("id"), req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// ── BOM ───────────────────────────────────────────────────────────────────────

// ListBOMs godoc
// @Summary      รายชื่อ BOM ทั้งหมด
// @Tags         master-data
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]domain.BOM
// @Router       /master-data/bom [get]
func (h *MasterDataHandler) ListBOMs(c *gin.Context) {
	items, err := h.svc.ListBOMs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetBOM godoc
// @Summary      ดึง BOM ตาม ID (พร้อม lines)
// @Tags         master-data
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "BOM ID (UUID)"
// @Success      200  {object}  domain.BOM
// @Failure      404  {object}  map[string]string
// @Router       /master-data/bom/{id} [get]
func (h *MasterDataHandler) GetBOM(c *gin.Context) {
	item, err := h.svc.GetBOM(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateBOM godoc
// @Summary      สร้าง BOM ใหม่ (พร้อม lines)
// @Tags         master-data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CreateBOMRequest  true  "ข้อมูล BOM + lines"
// @Success      201   {object}  domain.BOM
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Router       /master-data/bom [post]
func (h *MasterDataHandler) CreateBOM(c *gin.Context) {
	var req dto.CreateBOMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	item, err := h.svc.CreateBOM(c.Request.Context(), &req)
	if err != nil {
		if isDuplicateErr(err) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateBOM godoc
// @Summary      แก้ไข BOM (replace lines ทั้งหมด)
// @Tags         master-data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string              true  "BOM ID"
// @Param        body  body      dto.UpdateBOMRequest  true  "ข้อมูล BOM + lines ใหม่"
// @Success      200   {object}  domain.BOM
// @Failure      400   {object}  map[string]string
// @Router       /master-data/bom/{id} [put]
func (h *MasterDataHandler) UpdateBOM(c *gin.Context) {
	var req dto.UpdateBOMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	item, err := h.svc.UpdateBOM(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		if isDuplicateErr(err) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// ── Customers ─────────────────────────────────────────────────────────────────

// ListCustomers godoc
// @Summary      รายชื่อลูกค้าทั้งหมด
// @Tags         master-data
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]domain.Customer
// @Router       /master-data/customers [get]
func (h *MasterDataHandler) ListCustomers(c *gin.Context) {
	items, err := h.svc.ListCustomers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetCustomer godoc
// @Summary      ดึงลูกค้าตาม ID
// @Tags         master-data
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Customer ID (UUID)"
// @Success      200  {object}  domain.Customer
// @Failure      404  {object}  map[string]string
// @Router       /master-data/customers/{id} [get]
func (h *MasterDataHandler) GetCustomer(c *gin.Context) {
	item, err := h.svc.GetCustomer(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateCustomer godoc
// @Summary      สร้างลูกค้าใหม่
// @Tags         master-data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CreateCustomerRequest  true  "ข้อมูลลูกค้า"
// @Success      201   {object}  domain.Customer
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Router       /master-data/customers [post]
func (h *MasterDataHandler) CreateCustomer(c *gin.Context) {
	var req dto.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	item, err := h.svc.CreateCustomer(c.Request.Context(), &req)
	if err != nil {
		if isDuplicateErr(err) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateCustomer godoc
// @Summary      แก้ไขลูกค้า
// @Tags         master-data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                     true  "Customer ID"
// @Param        body  body      dto.UpdateCustomerRequest  true  "ข้อมูลที่แก้ไข"
// @Success      200   {object}  domain.Customer
// @Failure      400   {object}  map[string]string
// @Router       /master-data/customers/{id} [put]
func (h *MasterDataHandler) UpdateCustomer(c *gin.Context) {
	var req dto.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	item, err := h.svc.UpdateCustomer(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// SetCustomerActive godoc
// @Summary      เปิด/ปิดลูกค้า
// @Tags         master-data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                     true  "Customer ID"
// @Param        body  body      object{is_active=boolean}  true  "สถานะ"
// @Success      200   {object}  map[string]string
// @Router       /master-data/customers/{id}/active [put]
func (h *MasterDataHandler) SetCustomerActive(c *gin.Context) {
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := h.svc.SetCustomerActive(c.Request.Context(), c.Param("id"), req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// ── Suppliers ─────────────────────────────────────────────────────────────────

// ListSuppliers godoc
// @Summary      รายชื่อซัพพลายเออร์ทั้งหมด
// @Tags         master-data
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]domain.Supplier
// @Router       /master-data/suppliers [get]
func (h *MasterDataHandler) ListSuppliers(c *gin.Context) {
	items, err := h.svc.ListSuppliers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetSupplier godoc
// @Summary      ดึงซัพพลายเออร์ตาม ID
// @Tags         master-data
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Supplier ID (UUID)"
// @Success      200  {object}  domain.Supplier
// @Failure      404  {object}  map[string]string
// @Router       /master-data/suppliers/{id} [get]
func (h *MasterDataHandler) GetSupplier(c *gin.Context) {
	item, err := h.svc.GetSupplier(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateSupplier godoc
// @Summary      สร้างซัพพลายเออร์ใหม่
// @Tags         master-data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CreateSupplierRequest  true  "ข้อมูลซัพพลายเออร์"
// @Success      201   {object}  domain.Supplier
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Router       /master-data/suppliers [post]
func (h *MasterDataHandler) CreateSupplier(c *gin.Context) {
	var req dto.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	item, err := h.svc.CreateSupplier(c.Request.Context(), &req)
	if err != nil {
		if isDuplicateErr(err) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateSupplier godoc
// @Summary      แก้ไขซัพพลายเออร์
// @Tags         master-data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                     true  "Supplier ID"
// @Param        body  body      dto.UpdateSupplierRequest  true  "ข้อมูลที่แก้ไข"
// @Success      200   {object}  domain.Supplier
// @Failure      400   {object}  map[string]string
// @Router       /master-data/suppliers/{id} [put]
func (h *MasterDataHandler) UpdateSupplier(c *gin.Context) {
	var req dto.UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	item, err := h.svc.UpdateSupplier(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// SetSupplierActive godoc
// @Summary      เปิด/ปิดซัพพลายเออร์
// @Tags         master-data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                     true  "Supplier ID"
// @Param        body  body      object{is_active=boolean}  true  "สถานะ"
// @Success      200   {object}  map[string]string
// @Router       /master-data/suppliers/{id}/active [put]
func (h *MasterDataHandler) SetSupplierActive(c *gin.Context) {
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := h.svc.SetSupplierActive(c.Request.Context(), c.Param("id"), req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}
