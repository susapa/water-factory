import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map } from 'rxjs';
import { environment } from '../../../environments/environment';
import { ProductionOrder, ProductionYield } from '../models/production.model';

@Injectable({ providedIn: 'root' })
export class ProductionService {
  private http = inject(HttpClient);
  private base = `${environment.apiUrl}/production`;

  listOrders(): Observable<ProductionOrder[]> {
    return this.http.get<{ data: ProductionOrder[] }>(`${this.base}/orders`).pipe(map(r => r.data));
  }

  getOrder(id: string): Observable<ProductionOrder> {
    return this.http.get<ProductionOrder>(`${this.base}/orders/${id}`);
  }

  createOrder(body: object): Observable<ProductionOrder> {
    return this.http.post<ProductionOrder>(`${this.base}/orders`, body);
  }

  updateOrder(id: string, body: object): Observable<ProductionOrder> {
    return this.http.put<ProductionOrder>(`${this.base}/orders/${id}`, body);
  }

  confirmOrder(id: string): Observable<ProductionOrder> {
    return this.http.post<ProductionOrder>(`${this.base}/orders/${id}/confirm`, {});
  }

  issueRM(id: string, body: object): Observable<ProductionOrder> {
    return this.http.post<ProductionOrder>(`${this.base}/orders/${id}/issue-rm`, body);
  }

  startProduction(id: string): Observable<ProductionOrder> {
    return this.http.post<ProductionOrder>(`${this.base}/orders/${id}/start`, {});
  }

  recordYield(id: string, body: object): Observable<ProductionYield> {
    return this.http.post<ProductionYield>(`${this.base}/orders/${id}/record-yield`, body);
  }

  completeOrder(id: string): Observable<ProductionOrder> {
    return this.http.post<ProductionOrder>(`${this.base}/orders/${id}/complete`, {});
  }

  cancelOrder(id: string): Observable<void> {
    return this.http.post<void>(`${this.base}/orders/${id}/cancel`, {});
  }

  listYields(orderId?: string): Observable<ProductionYield[]> {
    const params = orderId ? `?order_id=${orderId}` : '';
    return this.http.get<{ data: ProductionYield[] }>(`${this.base}/yields${params}`).pipe(map(r => r.data));
  }
}
