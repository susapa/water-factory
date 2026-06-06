package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/water-factory/api/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// List godoc
// @Summary      รายชื่อผู้ใช้ทั้งหมด
// @Description  ดึงรายชื่อ user ทุกคน (admin only)
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /users [get]
func (h *UserHandler) List(c *gin.Context) {
	users, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// Create godoc
// @Summary      สร้างผู้ใช้ใหม่
// @Description  สร้าง user account ใหม่ (admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  object{email=string,password=string,full_name=string,role_id=int}  true  "ข้อมูล user"
// @Success      201   {object}  domain.User
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Router       /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"full_name" binding:"required"`
		RoleID   int    `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := h.svc.Create(c.Request.Context(), req.Email, req.Password, req.FullName, req.RoleID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// SetActive godoc
// @Summary      เปิด/ปิดการใช้งาน user
// @Description  Toggle is_active ของ user (admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                    true  "User ID (UUID)"
// @Param        body  body  object{is_active=boolean}  true  "สถานะ"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Router       /users/{id}/active [put]
func (h *UserHandler) SetActive(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Active bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := h.svc.SetActive(c.Request.Context(), id, req.Active); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// ListRoles godoc
// @Summary      รายชื่อ roles ทั้งหมด
// @Description  ดึง roles สำหรับ dropdown สร้าง user
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Router       /roles [get]
func (h *UserHandler) ListRoles(c *gin.Context) {
	_ = strconv.Itoa(0)
	roles, err := h.svc.ListRoles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": roles})
}
