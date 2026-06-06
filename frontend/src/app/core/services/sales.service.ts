import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import {
  SalesOrder, DeliveryOrder, Invoice, Vehicle,
  CreateSalesOrderRequest, CreateDeliveryOrderRequest,
  CreateInvoiceRequest, MarkPaidRequest,
} from '../models/sales.model';

@Injectable({ providedIn: 'root' })
export class SalesService {
  private http = inject(HttpClient);
  private base = '/api/v1/sales';

  // Vehicles
  listVehicles(): Observable<Vehicle[]> {
    return this.http.get<{ data: Vehicle[] }>(`${this.base}/vehicles`).pipe(map(r => r.data));
  }

  // Sales Orders
  listSalesOrders(): Observable<SalesOrder[]> {
    return this.http.get<{ data: SalesOrder[] }>(`${this.base}/orders`).pipe(map(r => r.data));
  }
  getSalesOrder(id: string): Observable<SalesOrder> {
    return this.http.get<{ data: SalesOrder }>(`${this.base}/orders/${id}`).pipe(map(r => r.data));
  }
  createSalesOrder(req: CreateSalesOrderRequest): Observable<SalesOrder> {
    return this.http.post<{ data: SalesOrder }>(`${this.base}/orders`, req).pipe(map(r => r.data));
  }
  updateSalesOrder(id: string, req: Partial<CreateSalesOrderRequest>): Observable<SalesOrder> {
    return this.http.put<{ data: SalesOrder }>(`${this.base}/orders/${id}`, req).pipe(map(r => r.data));
  }
  confirmSalesOrder(id: string): Observable<SalesOrder> {
    return this.http.post<{ data: SalesOrder }>(`${this.base}/orders/${id}/confirm`, {}).pipe(map(r => r.data));
  }
  cancelSalesOrder(id: string): Observable<SalesOrder> {
    return this.http.post<{ data: SalesOrder }>(`${this.base}/orders/${id}/cancel`, {}).pipe(map(r => r.data));
  }

  // Delivery Orders
  listDeliveryOrders(): Observable<DeliveryOrder[]> {
    return this.http.get<{ data: DeliveryOrder[] }>(`${this.base}/delivery-orders`).pipe(map(r => r.data));
  }
  getDeliveryOrder(id: string): Observable<DeliveryOrder> {
    return this.http.get<{ data: DeliveryOrder }>(`${this.base}/delivery-orders/${id}`).pipe(map(r => r.data));
  }
  createDeliveryOrder(req: CreateDeliveryOrderRequest): Observable<DeliveryOrder> {
    return this.http.post<{ data: DeliveryOrder }>(`${this.base}/delivery-orders`, req).pipe(map(r => r.data));
  }
  dispatchDeliveryOrder(id: string): Observable<DeliveryOrder> {
    return this.http.post<{ data: DeliveryOrder }>(`${this.base}/delivery-orders/${id}/dispatch`, {}).pipe(map(r => r.data));
  }
  deliverDeliveryOrder(id: string): Observable<DeliveryOrder> {
    return this.http.post<{ data: DeliveryOrder }>(`${this.base}/delivery-orders/${id}/deliver`, {}).pipe(map(r => r.data));
  }

  // Invoices
  listInvoices(): Observable<Invoice[]> {
    return this.http.get<{ data: Invoice[] }>(`${this.base}/invoices`).pipe(map(r => r.data));
  }
  getInvoice(id: string): Observable<Invoice> {
    return this.http.get<{ data: Invoice }>(`${this.base}/invoices/${id}`).pipe(map(r => r.data));
  }
  createInvoice(req: CreateInvoiceRequest): Observable<Invoice> {
    return this.http.post<{ data: Invoice }>(`${this.base}/invoices`, req).pipe(map(r => r.data));
  }
  markInvoicePaid(id: string, req: MarkPaidRequest): Observable<Invoice> {
    return this.http.post<{ data: Invoice }>(`${this.base}/invoices/${id}/pay`, req).pipe(map(r => r.data));
  }
}
