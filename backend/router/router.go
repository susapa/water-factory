package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/water-factory/api/docs"
	"github.com/water-factory/api/internal/handler"
	"github.com/water-factory/api/internal/middleware"
	jwtpkg "github.com/water-factory/api/pkg/jwt"
)

func Setup(
	engine *gin.Engine,
	jwtMgr *jwtpkg.Manager,
	authH *handler.AuthHandler,
	userH *handler.UserHandler,
	mdH *handler.MasterDataHandler,
	rmInvH *handler.RMInventoryHandler,
	prodH *handler.ProductionHandler,
	fgInvH *handler.FGInventoryHandler,
	salesH *handler.SalesHandler,
	dashboardH *handler.DashboardHandler,
) {
	engine.Use(middleware.CORS(), middleware.Logger())

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Swagger UI — http://localhost:8080/swagger/index.html
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := engine.Group("/api/v1")

	// Public auth routes
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authH.Login)
		auth.POST("/refresh", authH.Refresh)
		auth.POST("/logout", authH.Logout)
	}

	// Protected routes
	protected := v1.Group("", middleware.Auth(jwtMgr))
	{
		protected.GET("/auth/me", authH.Me)

		// User management (admin only)
		users := protected.Group("/users", middleware.RequirePermission("users", "read"))
		{
			users.GET("", userH.List)
			users.POST("", middleware.RequirePermission("users", "create"), userH.Create)
			users.PUT("/:id/active", middleware.RequirePermission("users", "update"), userH.SetActive)
		}

		protected.GET("/roles", middleware.RequirePermission("users", "read"), userH.ListRoles)

		// Master Data
		md := protected.Group("/master-data", middleware.RequirePermission("master_data", "read"))
		{
			md.GET("/uom", mdH.ListUOM)
			md.GET("/raw-material-categories", mdH.ListCategories)

			// Raw Materials
			rm := md.Group("/raw-materials")
			{
				rm.GET("", mdH.ListRawMaterials)
				rm.POST("", middleware.RequirePermission("master_data", "create"), mdH.CreateRawMaterial)
				rm.GET("/:id", mdH.GetRawMaterial)
				rm.PUT("/:id", middleware.RequirePermission("master_data", "update"), mdH.UpdateRawMaterial)
				rm.PUT("/:id/active", middleware.RequirePermission("master_data", "update"), mdH.SetRawMaterialActive)
			}

			// Finished Goods
			fg := md.Group("/finished-goods")
			{
				fg.GET("", mdH.ListFinishedGoods)
				fg.POST("", middleware.RequirePermission("master_data", "create"), mdH.CreateFinishedGood)
				fg.GET("/:id", mdH.GetFinishedGood)
				fg.PUT("/:id", middleware.RequirePermission("master_data", "update"), mdH.UpdateFinishedGood)
				fg.PUT("/:id/active", middleware.RequirePermission("master_data", "update"), mdH.SetFinishedGoodActive)
			}

			// BOM
			bom := md.Group("/bom")
			{
				bom.GET("", mdH.ListBOMs)
				bom.POST("", middleware.RequirePermission("master_data", "create"), mdH.CreateBOM)
				bom.GET("/:id", mdH.GetBOM)
				bom.PUT("/:id", middleware.RequirePermission("master_data", "update"), mdH.UpdateBOM)
			}

			// Customers
			cust := md.Group("/customers")
			{
				cust.GET("", mdH.ListCustomers)
				cust.POST("", middleware.RequirePermission("master_data", "create"), mdH.CreateCustomer)
				cust.GET("/:id", mdH.GetCustomer)
				cust.PUT("/:id", middleware.RequirePermission("master_data", "update"), mdH.UpdateCustomer)
				cust.PUT("/:id/active", middleware.RequirePermission("master_data", "update"), mdH.SetCustomerActive)
			}

			// Suppliers
			supp := md.Group("/suppliers")
			{
				supp.GET("", mdH.ListSuppliers)
				supp.POST("", middleware.RequirePermission("master_data", "create"), mdH.CreateSupplier)
				supp.GET("/:id", mdH.GetSupplier)
				supp.PUT("/:id", middleware.RequirePermission("master_data", "update"), mdH.UpdateSupplier)
				supp.PUT("/:id/active", middleware.RequirePermission("master_data", "update"), mdH.SetSupplierActive)
			}
		}

		// Inventory — Raw Material
		inv := protected.Group("/inventory/raw-material", middleware.RequirePermission("raw_material", "read"))
		{
			inv.GET("/warehouse-locations", rmInvH.ListWarehouseLocations)

			grn := inv.Group("/grn")
			{
				grn.GET("", rmInvH.ListGRNs)
				grn.POST("", middleware.RequirePermission("raw_material", "create"), rmInvH.CreateGRN)
				grn.GET("/:id", rmInvH.GetGRN)
				grn.PUT("/:id", middleware.RequirePermission("raw_material", "update"), rmInvH.UpdateGRN)
				grn.POST("/:id/confirm", middleware.RequirePermission("raw_material", "update"), rmInvH.ConfirmGRN)
				grn.POST("/:id/cancel", middleware.RequirePermission("raw_material", "update"), rmInvH.CancelGRN)
			}

			stock := inv.Group("/stock")
			{
				stock.GET("/summary", rmInvH.ListStockSummary)
				stock.GET("/lots", rmInvH.ListStockLots)
				stock.GET("/lots/:id", rmInvH.GetStockLotDetail)
			}

			adj := inv.Group("/adjustments")
			{
				adj.GET("", rmInvH.ListAdjustments)
				adj.POST("", middleware.RequirePermission("raw_material", "create"), rmInvH.CreateAdjustment)
				adj.GET("/:id", rmInvH.GetAdjustment)
				adj.POST("/:id/approve", middleware.RequirePermission("raw_material", "update"), rmInvH.ApproveAdjustment)
				adj.POST("/:id/cancel", middleware.RequirePermission("raw_material", "update"), rmInvH.CancelAdjustment)
			}
		}

		// Inventory — Finished Goods
		fgInv := protected.Group("/inventory/finished-goods", middleware.RequirePermission("finished_goods", "read"))
		{
			fgStock := fgInv.Group("/stock")
			{
				fgStock.GET("/summary", fgInvH.ListStockSummary)
				fgStock.GET("/lots", fgInvH.ListStockLots)
				fgStock.GET("/lots/:id", fgInvH.GetStockLotDetail)
			}
			fgAdj := fgInv.Group("/adjustments")
			{
				fgAdj.GET("", fgInvH.ListAdjustments)
				fgAdj.POST("", middleware.RequirePermission("finished_goods", "create"), fgInvH.CreateAdjustment)
				fgAdj.GET("/:id", fgInvH.GetAdjustment)
				fgAdj.POST("/:id/approve", middleware.RequirePermission("finished_goods", "update"), fgInvH.ApproveAdjustment)
				fgAdj.POST("/:id/cancel", middleware.RequirePermission("finished_goods", "update"), fgInvH.CancelAdjustment)
			}
		}

		// Production
		prod := protected.Group("/production", middleware.RequirePermission("production", "read"))
		{
			orders := prod.Group("/orders")
			{
				orders.GET("", prodH.ListOrders)
				orders.POST("", middleware.RequirePermission("production", "create"), prodH.CreateOrder)
				orders.GET("/:id", prodH.GetOrder)
				orders.PUT("/:id", middleware.RequirePermission("production", "update"), prodH.UpdateOrder)
				orders.POST("/:id/confirm", middleware.RequirePermission("production", "update"), prodH.ConfirmOrder)
				orders.POST("/:id/issue-rm", middleware.RequirePermission("production", "update"), prodH.IssueRM)
				orders.POST("/:id/start", middleware.RequirePermission("production", "update"), prodH.StartProduction)
				orders.POST("/:id/record-yield", middleware.RequirePermission("production", "update"), prodH.RecordYield)
				orders.POST("/:id/complete", middleware.RequirePermission("production", "update"), prodH.CompleteOrder)
				orders.POST("/:id/cancel", middleware.RequirePermission("production", "update"), prodH.CancelOrder)
			}
			prod.GET("/yields", prodH.ListYields)
		}

		// Dashboard
		dash := protected.Group("/dashboard", middleware.RequirePermission("dashboard", "read"))
		{
			dash.GET("/summary", dashboardH.GetSummary)
		}

		// Sales & Delivery
		sales := protected.Group("/sales", middleware.RequirePermission("sales", "read"))
		{
			sales.GET("/vehicles", salesH.ListVehicles)

			so := sales.Group("/orders")
			{
				so.GET("", salesH.ListSalesOrders)
				so.POST("", middleware.RequirePermission("sales", "create"), salesH.CreateSalesOrder)
				so.GET("/:id", salesH.GetSalesOrder)
				so.PUT("/:id", middleware.RequirePermission("sales", "update"), salesH.UpdateSalesOrder)
				so.POST("/:id/confirm", middleware.RequirePermission("sales", "update"), salesH.ConfirmSalesOrder)
				so.POST("/:id/cancel", middleware.RequirePermission("sales", "update"), salesH.CancelSalesOrder)
			}

			doGrp := sales.Group("/delivery-orders")
			{
				doGrp.GET("", salesH.ListDeliveryOrders)
				doGrp.POST("", middleware.RequirePermission("sales", "create"), salesH.CreateDeliveryOrder)
				doGrp.GET("/:id", salesH.GetDeliveryOrder)
				doGrp.POST("/:id/dispatch", middleware.RequirePermission("sales", "update"), salesH.DispatchDeliveryOrder)
				doGrp.POST("/:id/deliver", middleware.RequirePermission("sales", "update"), salesH.DeliverDeliveryOrder)
			}

			inv := sales.Group("/invoices")
			{
				inv.GET("", salesH.ListInvoices)
				inv.POST("", middleware.RequirePermission("sales", "create"), salesH.CreateInvoice)
				inv.GET("/:id", salesH.GetInvoice)
				inv.POST("/:id/pay", middleware.RequirePermission("sales", "update"), salesH.MarkInvoicePaid)
			}
		}
	}
}
