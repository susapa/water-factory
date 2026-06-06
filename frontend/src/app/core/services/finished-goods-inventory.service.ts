import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map } from 'rxjs';
import { environment } from '../../../environments/environment';
import {
  FGStockLot, FGStockLotDetail, FGStockSummary, FGAdjustment,
} from '../models/finished-goods-inventory.model';

@Injectable({ providedIn: 'root' })
export class FinishedGoodsInventoryService {
  private http = inject(HttpClient);
  private base = `${environment.apiUrl}/inventory/finished-goods`;

  // Stock
  listStockSummary(): Observable<FGStockSummary[]> {
    return this.http.get<{ data: FGStockSummary[] }>(`${this.base}/stock/summary`).pipe(map(r => r.data));
  }

  listStockLots(finishedGoodId?: string): Observable<FGStockLot[]> {
    const params = finishedGoodId ? `?finished_good_id=${finishedGoodId}` : '';
    return this.http.get<{ data: FGStockLot[] }>(`${this.base}/stock/lots${params}`).pipe(map(r => r.data));
  }

  getStockLotDetail(id: string): Observable<FGStockLotDetail> {
    return this.http.get<FGStockLotDetail>(`${this.base}/stock/lots/${id}`);
  }

  // Adjustments
  listAdjustments(): Observable<FGAdjustment[]> {
    return this.http.get<{ data: FGAdjustment[] }>(`${this.base}/adjustments`).pipe(map(r => r.data));
  }

  getAdjustment(id: string): Observable<FGAdjustment> {
    return this.http.get<FGAdjustment>(`${this.base}/adjustments/${id}`);
  }

  createAdjustment(body: object): Observable<FGAdjustment> {
    return this.http.post<FGAdjustment>(`${this.base}/adjustments`, body);
  }

  approveAdjustment(id: string): Observable<FGAdjustment> {
    return this.http.post<FGAdjustment>(`${this.base}/adjustments/${id}/approve`, {});
  }

  cancelAdjustment(id: string): Observable<void> {
    return this.http.post<void>(`${this.base}/adjustments/${id}/cancel`, {});
  }
}
