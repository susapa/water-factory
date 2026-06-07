# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# Water Factory Management System

**GitHub:** https://github.com/susapa/water-factory

## Stack
- **Frontend:** Angular 21, PrimeNG 21 (Aura theme), SCSS แยกไฟล์, standalone components, port 4200
- **Backend:** Go 1.25 + Gin, port 8080
- **Database:** PostgreSQL 17 (local), database = `water_factory`, user = `water_user`
- **Auth:** JWT (access 15m, refresh 7d), bcrypt passwords
- **Docs:** Swagger UI → http://localhost:8080/swagger/index.html

## การรัน (Development)

### Database
PostgreSQL 17 รันอยู่แล้ว (Windows service `postgresql-x64-17`)
```
host: localhost:5432
database: water_factory
user: water_user
password: WaterFactory2026!
```

### Backend
```powershell
cd C:\AI\water\backend
go run ./cmd/server
# .env อยู่ใน backend/ แล้ว
```

### Frontend
```powershell
cd C:\AI\water\frontend
npm run start
```

### Swagger (regenerate หลัง update annotations)
```powershell
cd C:\AI\water\backend
& "$env:USERPROFILE\go\bin\swag.exe" init -g cmd/server/main.go --output docs
```

## Admin User (seed แล้ว)
- Email: `admin@water.local`
- Password: `Admin@1234`

## โครงสร้าง backend
```
backend/
├── cmd/server/main.go          ← entry point + wire dependencies
├── cmd/hashgen/main.go         ← utility: gen bcrypt hash
├── config/config.go            ← Viper env loader (.env)
├── docs/                       ← Swagger generated (swag init)
├── internal/
│   ├── database/
│   │   ├── postgres.go         ← pgxpool connection
│   │   ├── migrate.go          ← simple SQL migration runner
│   │   └── migrations/         ← 001..007 SQL files (ทั้งหมด applied แล้ว)
│   ├── middleware/
│   │   ├── auth.go             ← JWT validation middleware
│   │   ├── rbac.go             ← RequirePermission(module, action)
│   │   ├── cors.go
│   │   └── logger.go
│   ├── domain/
│   │   ├── user.go
│   │   ├── master_data.go      ← UOM, RawMaterial, FinishedGood, BOM, Customer, Supplier
│   │   ├── raw_material_inventory.go  ← GRN, StockLot, StockMovement, Adjustment ฯลฯ
│   │   ├── production.go       ← ProductionOrder, ProductionRequirement, ProductionYield, DefectDetail
│   │   ├── fg_inventory.go     ← FGStockLot, FGStockMovement, FGStockSummary, FGAdjustment ฯลฯ
│   │   ├── sales.go            ← SalesOrder, SalesOrderLine, Vehicle, DeliveryOrder, DeliveryOrderLine, Invoice
│   │   └── dashboard.go        ← DashboardSummary, DailySales, FGExpiringSoonItem
│   ├── dto/
│   │   ├── auth_dto.go
│   │   ├── master_data_dto.go  ← Create/Update requests for all master data
│   │   ├── raw_material_inventory_dto.go  ← GRN/Adjustment request DTOs
│   │   ├── production_dto.go   ← CreateOrder, UpdateOrder, IssueRMRequest, RecordYieldRequest
│   │   ├── fg_inventory_dto.go ← CreateFGAdjustmentRequest, FGAdjustmentLineRequest
│   │   └── sales_dto.go        ← CreateSalesOrderRequest, CreateDeliveryOrderRequest, CreateInvoiceRequest, MarkPaidRequest
│   ├── repository/
│   │   ├── user_repo.go
│   │   ├── master_data_repo.go        ← pgxpool queries, BOM ใช้ transaction
│   │   ├── raw_material_inventory_repo.go  ← GRN confirm + Adjustment approve transactions
│   │   ├── production_repo.go         ← ConfirmOrder (BOM explosion), IssueRM (FIFO), RecordYield transactions
│   │   │                                   RecordYield ยัง insert fg_stock_lot + PRODUCTION_RECEIPT movement ด้วย
│   │   ├── fg_inventory_repo.go       ← ListStockSummary, ListStockLots, GetStockLotDetail, Adjustment CRUD
│   │   ├── sales_repo.go              ← SO CRUD+confirm/cancel, DO create(FEFO)+dispatch+deliver, Invoice create+pay
│   │   └── dashboard_repo.go          ← GetSummary() — 4 aggregate queries, read-only
│   ├── service/
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── master_data_service.go
│   │   ├── raw_material_inventory_service.go
│   │   ├── production_service.go
│   │   ├── fg_inventory_service.go
│   │   ├── sales_service.go
│   │   └── dashboard_service.go
│   └── handler/
│       ├── auth_handler.go
│       ├── user_handler.go
│       ├── master_data_handler.go         ← 22 handlers พร้อม Swagger annotations
│       ├── raw_material_inventory_handler.go  ← 15 handlers พร้อม Swagger annotations
│       ├── production_handler.go          ← 11 handlers พร้อม Swagger annotations
│       ├── fg_inventory_handler.go        ← 8 handlers พร้อม Swagger annotations
│       ├── sales_handler.go               ← 13 handlers พร้อม Swagger annotations
│       └── dashboard_handler.go           ← 1 handler: GetSummary
├── pkg/
│   ├── jwt/jwt.go
│   ├── number/generator.go     ← doc number: GRN-YYYYMMDD-00001, PO-YYYYMMDD-00001 ฯลฯ
│   ├── errs/errors.go
│   └── validator/validator.go
└── router/router.go            ← ทุก route + RBAC + Swagger UI
```

## โครงสร้าง frontend
```
frontend/src/app/
├── core/
│   ├── auth/
│   │   ├── auth.service.ts     ← signal-based, localStorage tokens
│   │   ├── auth.guard.ts       ← CanActivateFn
│   │   ├── role.guard.ts       ← เช็ค route data.roles
│   │   └── jwt.interceptor.ts  ← auto refresh on 401
│   ├── models/
│   │   ├── user.model.ts
│   │   ├── api-response.model.ts
│   │   ├── master-data.model.ts                ← UOM, RawMaterial, FinishedGood, BOM, Customer, Supplier
│   │   ├── raw-material-inventory.model.ts     ← GRN, StockLot, StockMovement, Adjustment ฯลฯ
│   │   ├── production.model.ts                 ← ProductionOrder, ProductionRequirement, ProductionYield, DefectDetail
│   │   ├── finished-goods-inventory.model.ts   ← FGStockLot, FGStockSummary, FGAdjustment ฯลฯ
│   │   ├── sales.model.ts                      ← SalesOrder, DeliveryOrder, Invoice, Vehicle, DTOs
│   │   └── dashboard.model.ts                  ← DashboardSummary, DailySales, FGExpiringSoonItem
│   └── services/
│       ├── master-data.service.ts              ← HTTP service สำหรับ master data
│       ├── raw-material-inventory.service.ts   ← HTTP service สำหรับ Phase 2
│       ├── production.service.ts               ← HTTP service สำหรับ Phase 3
│       ├── finished-goods-inventory.service.ts ← HTTP service สำหรับ Phase 4
│       ├── sales.service.ts                    ← HTTP service สำหรับ Phase 5
│       ├── dashboard.service.ts                ← getDashboardSummary() สำหรับ Phase 6
│       └── user.service.ts                     ← listUsers, createUser, setUserActive, listRoles
├── layout/
│   ├── layout.component.ts/html/scss
│   ├── sidebar/                ← role-aware nav, computed() signal
│   └── topbar/                 ← user info + logout
└── features/
    ├── auth/login/             ← login page (Reactive Forms)
    ├── dashboard/              ← ✅ Phase 6 เสร็จแล้ว — KPI tiles + bar/donut charts + expiry alerts
    ├── master-data/            ← ✅ Phase 1 เสร็จแล้ว
    │   ├── master-data-shell.component.*   ← tab navigation
    │   ├── master-data.routes.ts
    │   ├── _list-shared.scss               ← shared styles (reuse ด้วย @use '../../master-data/list-shared' as *)
    │   ├── raw-materials/raw-material-list.component.*
    │   ├── finished-goods/finished-good-list.component.*
    │   ├── bom/bom-list.component.*        ← FormArray lines pattern
    │   ├── customers/customer-list.component.*
    │   └── suppliers/supplier-list.component.*
    ├── raw-material/           ← ✅ Phase 2 เสร็จแล้ว
    │   ├── raw-material-shell.component.*  ← tab navigation
    │   ├── raw-material.routes.ts
    │   ├── grn/grn-list.component.*        ← GRN create/confirm/cancel + FormArray lines
    │   ├── stock/stock-list.component.*    ← stock summary + lot detail + movement history
    │   └── adjustments/adjustment-list.component.*  ← adjustment create/approve/cancel
    ├── production/             ← ✅ Phase 3 เสร็จแล้ว
    │   ├── production-shell.component.*    ← tab navigation
    │   ├── production.routes.ts
    │   ├── orders/order-list.component.*   ← create/confirm/issue-rm/start/yield/complete/cancel + 4 dialogs
    │   └── yields/yield-list.component.*   ← yield records + row expand แสดง defect details
    ├── finished-goods/         ← ✅ Phase 4 เสร็จแล้ว
    │   ├── finished-goods-shell.component.*  ← tab navigation
    │   ├── finished-goods.routes.ts
    │   ├── stock/fg-stock-list.component.*   ← stock summary + lot detail (FEFO) + expiry alerts + movement history
    │   └── adjustments/fg-adjustment-list.component.*  ← adjustment create/approve/cancel
    ├── sales/                  ← ✅ Phase 5 เสร็จแล้ว
    │   ├── sales-shell.component.*      ← tab navigation (role-aware: delivery เห็นแค่จัดส่ง)
    │   ├── sales.routes.ts
    │   ├── orders/so-list.component.*   ← SO create/confirm/cancel + FormArray lines + totals
    │   ├── deliveries/do-list.component.*  ← DO create(FEFO auto-pick)+dispatch+deliver
    │   └── invoices/invoice-list.component.*  ← invoice create+pay
    ├── reports/                ← ✅ Phase 6 เสร็จแล้ว
    │   ├── reports-shell.component.*  ← tab navigation (3 tabs)
    │   ├── reports.routes.ts
    │   ├── stock/stock-report.component.*      ← RM + FG stock tables
    │   ├── sales/sales-report.component.*      ← SO list + date filter + total footer
    │   └── production/production-report.component.*  ← PO list + summary tiles + date filter
    └── admin/                  ← ✅ Phase 0 เสร็จแล้ว
        └── user-list.component.*  ← ตาราง users + dialog เพิ่มผู้ใช้ + toggle active
```

## RBAC Roles
| Role | Module Access |
|------|--------------|
| admin | ทุกอย่าง + user management |
| warehouse_manager | raw_material (full), finished_goods (full), production (read) |
| production_manager | production (full), raw_material (read) |
| sales_admin | sales (full), finished_goods (read) |
| delivery | sales (dispatch menu เท่านั้น) |

## API Endpoints ที่สร้างแล้ว

### Phase 0 — Auth & Users
```
GET  /health
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
GET  /api/v1/auth/me           (auth required)
GET  /api/v1/users             (admin only)
POST /api/v1/users             (admin only)
PUT  /api/v1/users/:id/active  (admin only)
GET  /api/v1/roles             (admin only)
GET  /swagger/*                (Swagger UI)
```

### Phase 1 — Master Data (admin only, master_data permission)
```
GET  /api/v1/master-data/uom
GET  /api/v1/master-data/raw-material-categories
GET/POST        /api/v1/master-data/raw-materials
GET/PUT         /api/v1/master-data/raw-materials/:id
PUT             /api/v1/master-data/raw-materials/:id/active
GET/POST        /api/v1/master-data/finished-goods
GET/PUT         /api/v1/master-data/finished-goods/:id
PUT             /api/v1/master-data/finished-goods/:id/active
GET/POST        /api/v1/master-data/bom
GET/PUT         /api/v1/master-data/bom/:id
GET/POST        /api/v1/master-data/customers
GET/PUT         /api/v1/master-data/customers/:id
PUT             /api/v1/master-data/customers/:id/active
GET/POST        /api/v1/master-data/suppliers
GET/PUT         /api/v1/master-data/suppliers/:id
PUT             /api/v1/master-data/suppliers/:id/active
```

### Phase 2 — Raw Material Inventory (raw_material permission)
```
GET  /api/v1/inventory/raw-material/warehouse-locations

# GRN (Goods Receipt Note)
GET  /api/v1/inventory/raw-material/grn                     (read)
POST /api/v1/inventory/raw-material/grn                     (create)
GET  /api/v1/inventory/raw-material/grn/:id                 (read)
PUT  /api/v1/inventory/raw-material/grn/:id                 (update, draft only)
POST /api/v1/inventory/raw-material/grn/:id/confirm         (update) → สร้าง stock lots + movements
POST /api/v1/inventory/raw-material/grn/:id/cancel          (update, draft only)

# Stock
GET  /api/v1/inventory/raw-material/stock/summary           ← GROUP BY RM, status=available
GET  /api/v1/inventory/raw-material/stock/lots              ← ?raw_material_id= optional
GET  /api/v1/inventory/raw-material/stock/lots/:id          ← lot + movement history

# Adjustments
GET  /api/v1/inventory/raw-material/adjustments
POST /api/v1/inventory/raw-material/adjustments             (create)
GET  /api/v1/inventory/raw-material/adjustments/:id
POST /api/v1/inventory/raw-material/adjustments/:id/approve (update) → อัปเดต stock lots
POST /api/v1/inventory/raw-material/adjustments/:id/cancel  (update, pending only)
```

### Phase 3 — Production (production permission)
```
GET  /api/v1/production/orders
POST /api/v1/production/orders
GET  /api/v1/production/orders/:id
PUT  /api/v1/production/orders/:id              (update, draft only)
POST /api/v1/production/orders/:id/confirm      → BOM explosion → snapshot rm_requirements
POST /api/v1/production/orders/:id/issue-rm     → FIFO lot deduction → rm_stock_movements (ISSUE_TO_PROD)
POST /api/v1/production/orders/:id/start        → status=in_progress, actual_start_date=NOW()
POST /api/v1/production/orders/:id/record-yield → insert production_yields + defect_details + fg_stock_lot + PRODUCTION_RECEIPT
POST /api/v1/production/orders/:id/complete     → status=completed, actual_end_date=NOW()
POST /api/v1/production/orders/:id/cancel       (draft/confirmed only)
GET  /api/v1/production/yields                  ← ?order_id= optional
```

### Phase 4 — Finished Goods Inventory (finished_goods permission)
```
# Stock
GET  /api/v1/inventory/finished-goods/stock/summary    ← GROUP BY FG, status=available + nearest_expiry
GET  /api/v1/inventory/finished-goods/stock/lots       ← ?finished_good_id= optional, sorted FEFO
GET  /api/v1/inventory/finished-goods/stock/lots/:id   ← lot + movement history

# Adjustments (doc number: FGADJ-YYYYMMDD-00001)
GET  /api/v1/inventory/finished-goods/adjustments
POST /api/v1/inventory/finished-goods/adjustments             (create, pending)
GET  /api/v1/inventory/finished-goods/adjustments/:id
POST /api/v1/inventory/finished-goods/adjustments/:id/approve → อัปเดต fg_stock_lots + insert movements
POST /api/v1/inventory/finished-goods/adjustments/:id/cancel  (pending only)
```

### Phase 5 — Sales (sales permission)
```
GET  /api/v1/sales/vehicles

# Sales Orders (SO-YYYYMMDD-00001)
GET  /api/v1/sales/orders
POST /api/v1/sales/orders                  (create, draft)
GET  /api/v1/sales/orders/:id
PUT  /api/v1/sales/orders/:id              (update, draft only)
POST /api/v1/sales/orders/:id/confirm      → status=confirmed
POST /api/v1/sales/orders/:id/cancel       (draft/confirmed only)

# Delivery Orders (DO-YYYYMMDD-00001)
GET  /api/v1/sales/delivery-orders
POST /api/v1/sales/delivery-orders         (FEFO auto-pick from confirmed SO) → status=pending, SO→picking
GET  /api/v1/sales/delivery-orders/:id
POST /api/v1/sales/delivery-orders/:id/dispatch  → deduct fg_stock_lots + SALES_DISPATCH movements, SO→dispatched
POST /api/v1/sales/delivery-orders/:id/deliver   → status=delivered

# Invoices (INV-YYYYMMDD-00001)
GET  /api/v1/sales/invoices
POST /api/v1/sales/invoices                (from dispatched SO) → status=issued, SO→invoiced
GET  /api/v1/sales/invoices/:id
POST /api/v1/sales/invoices/:id/pay        → status=paid
```

### Phase 6 — Dashboard (dashboard:read permission — admin, warehouse_manager, production_manager, sales_admin)
```
GET  /api/v1/dashboard/summary   ← KPI counts + sales_last_7_days[] + po_status_counts{} + fg_expiring_soon[]
```

## หน้า Frontend ที่เสร็จแล้ว
| หน้า | Route | Roles |
|------|-------|-------|
| Login | /login | - |
| จัดการผู้ใช้ | /admin | admin |
| ข้อมูลหลัก (5 tabs) | /master-data | admin |
| คลังวัตถุดิบ (3 tabs) | /raw-material | admin, warehouse_manager, production_manager |
| การผลิต (2 tabs) | /production | admin, warehouse_manager, production_manager |
| คลังสินค้าสำเร็จรูป (2 tabs) | /finished-goods | admin, warehouse_manager, sales_admin |
| ขายและจัดส่ง (3 tabs / 1 tab delivery) | /sales | admin, sales_admin, delivery |
| แดชบอร์ด | /dashboard | admin, warehouse_manager, production_manager, sales_admin |
| รายงาน (3 tabs: สต็อก/ยอดขาย/การผลิต) | /reports | admin, warehouse_manager, production_manager, sales_admin |

## Database Tables (ทั้งหมด applied แล้ว)
Migration 001: roles, users, refresh_tokens, sequences
Migration 002: units_of_measure, raw_material_categories, raw_materials, finished_goods, bill_of_materials, bom_lines, customers, suppliers, warehouse_locations
Migration 003: rm_goods_receipts, rm_goods_receipt_lines, rm_stock_lots, rm_stock_movements, rm_adjustments, rm_adjustment_lines
Migration 004: production_orders, production_order_rm_requirements, production_rm_issues, production_yields, production_defect_details
Migration 005: fg_stock_lots, fg_stock_movements, notifications
Migration 006: sales_orders, sales_order_lines, vehicles, delivery_orders, delivery_order_lines, invoices
Migration 007: fg_adjustments, fg_adjustment_lines

## หมายเหตุสำคัญ

### Angular / Frontend
- Angular 21 ใช้ `inject()` แทน constructor injection สำหรับ class field init
- ทุก component แยกไฟล์: `.ts` + `.html` + `.scss`
- SCSS ใช้ nesting syntax ได้, shared styles อยู่ที่ `features/master-data/_list-shared.scss`
  - reuse ด้วย `@use '../../master-data/list-shared' as *;` (ปรับ path ตาม depth)
- Angular build ผ่าน `npm run build` (ไม่ใช้ global `ng`)
- **State management ใช้ Signals เสมอ** — `signal()`, `computed()` สำหรับ `items`, `loading`, `searchText`, dropdown options ทุกตัว ห้ามใช้ plain property กับ async data เพราะ Angular change detection อาจไม่ trigger
- Template เรียก signal ด้วย `()` เสมอ: `[value]="items()"`, `[loading]="loading()"`
- Search input ใช้ `[value]="searchText()" (input)="searchText.set($any($event.target).value)"` ไม่ใช้ `[(ngModel)]`
- `filtered` list ทำเป็น `computed()` signal เสมอ
- **PrimeNG `[(selection)]` ต้องใช้ plain property** ไม่ใช้ signal — two-way binding จะพัง
- **PrimeNG `[severity]` ต้องคืน union type** `'success' | 'secondary' | 'info' | 'warn' | 'danger'` ไม่ใช่ `string`
- **Pipes ต้อง import ใน component** — `DatePipe`, `DecimalPipe` จาก `@angular/common` (standalone component ไม่มี CommonModule auto)
- PrimeNG 21.x ต้องมี `providePrimeNG({ theme: { preset: Aura } })` ใน `app.config.ts` — ถ้าไม่มี: dialog mask ขาว, table loading ค้าง
- PrimeNG theme package: `@primeng/themes` (install แล้ว), preset ที่ใช้: **Aura**
- `MessageService` ต้อง provide ใน component `providers: [MessageService]` ไม่ใช้ root
- `ConfirmationService` ต้อง provide ใน component `providers: [ConfirmationService]` พร้อมเพิ่ม `<p-confirmDialog />` ใน template
- **`Object.keys()` ใช้ใน Angular template ไม่ได้** — ถ้าต้องเช็ค empty object ให้ทำเป็น `computed()` signal แทน เช่น `hasCounts = computed(() => Object.keys(this.data()?.map ?? {}).length > 0)`
- **`chart.js` ต้อง install แยก** — `npm install chart.js` (PrimeNG `p-chart` ไม่ bundle มาด้วย)

### Go / Backend
- Go list functions ต้อง `make([]*T, 0)` ไม่ใช้ `var items []*T` — Go nil slice serialize เป็น JSON `null` ทำให้ Angular พัง
- Go deps: gin v1.12, swaggo/swag v1.16.6 (swag CLI อยู่ที่ `%USERPROFILE%\go\bin\swag.exe`)
- Swagger annotations อยู่ใน handler files, regenerate ด้วย `swag init -g cmd/server/main.go --output docs`
  - **ใช้ `main.go` ไม่ใช้ `docs.go`** เป็น entry point ของ swag
- `.env` อยู่ที่ `backend/.env` (ไม่ commit) — ต้องมีก่อนรัน server
- Status conflict (confirm ซ้ำ, approve ซ้ำ) → return HTTP 409 พร้อม message
- `number.Next(ctx, pool, "GRN")` → `GRN-YYYYMMDD-00001`, `number.Next(ctx, pool, "PO")` → `PO-YYYYMMDD-00001` (ใช้ sequences table)
- `rm_stock_movements.movement_type` CHECK constraint: `'GRN','ISSUE_TO_PROD','RETURN_FROM_PROD','ADJUSTMENT_IN','ADJUSTMENT_OUT','CYCLE_COUNT'` — ต้องใช้ `ISSUE_TO_PROD` ไม่ใช่ `production_issue`
- Production order status flow: `draft` → `confirmed` (BOM explosion) → `rm_issued` (all reqs fully_issued) → `in_progress` → `completed`; cancel ได้จาก `draft`/`confirmed` เท่านั้น
- IssueRM ใช้ FIFO: `ORDER BY expiry_date ASC NULLS LAST, received_date ASC`; สามารถ issue บางส่วนได้ (partial) — order status กลายเป็น `rm_issued` เมื่อทุก requirement เป็น `fully_issued`
- **computed() signal ที่ขึ้นกับ reactive form value** ต้องใช้ `valueChanges.subscribe(v => signal.set(v))` ไม่ใช่อ่าน `form.get('field')?.value` ตรงใน computed เพราะไม่ reactive
- **Expiry alert threshold (Phase 4):** `nearest_expiry` ≤ 7 วัน → severity `'warn'`; ผ่านแล้ว → severity `'danger'`; เกิน 7 วัน → แสดง text ปกติ
- **fg_stock_movements.movement_type** CHECK constraint: `'PRODUCTION_RECEIPT','SALES_DISPATCH','RETURN','ADJUSTMENT_IN','ADJUSTMENT_OUT'`
- `number.Next(ctx, pool, "FGADJ")` → `FGADJ-YYYYMMDD-00001`
- RecordYield (production_repo.go) ทำงาน atomic: insert production_yield → insert fg_stock_lot → insert fg_stock_movement(PRODUCTION_RECEIPT) ในตัว transaction เดียวกัน
- **Sales DO dispatch** transaction: update fg_stock_lots.current_qty → insert fg_stock_movements(SALES_DISPATCH) → update SO status=dispatched
- `number.Next(ctx, pool, "SO")` → `SO-YYYYMMDD-00001`, `"DO"` → `DO-...`, `"INV"` → `INV-...` (auto-insert ลงใน sequences table)
- **delivery role** เห็นแค่ tab "จัดส่ง" (deliveries) ใน sales shell — ซ่อน orders + invoices ด้วย `isDelivery()` computed signal
- **dashboard_repo.go** ใช้ 4 queries แยกกัน (ไม่ใช้ transaction เพราะ read-only): COUNT pending SO/PO/DO, SUM invoices, generate_series สำหรับ 7-day sales, GROUP BY status สำหรับ PO chart, และ JOIN fg_stock_lots+finished_goods สำหรับ expiry alerts
- **`dashboard.service.ts` ต้องใช้ `environment.apiUrl`** (absolute URL `http://localhost:8080/api/v1`) — ห้ามใช้ relative path `/api/v1/...` เพราะ Angular dev server ไม่มี proxy config จะ return HTML แทน JSON
- **`dashboard:read` permission** อยู่ใน `rbac.go` แล้วสำหรับ 4 roles: admin, warehouse_manager, production_manager, sales_admin — ไม่ต้องแก้ rbac.go
