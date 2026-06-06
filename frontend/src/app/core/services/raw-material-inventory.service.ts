import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map } from 'rxjs';
import { environment } from '../../../environments/environment';
import {
  WarehouseLocation, GRN, StockLot, StockLotDetail,
  StockSummary, Adjustment,
} from '../models/raw-material-inventory.model';

@Injectable({ providedIn: 'root' })
export class RawMaterialInventoryService {
  private http = inject(HttpClient);
  private base = `${environment.apiUrl}/inventory/raw-material`;

  listWarehouseLocations(): Observable<WarehouseLocation[]> {
    return this.http.get<{ data: WarehouseLocation[] }>(`${this.base}/warehouse-locations`).pipe(map(r => r.data));
  }

  // GRN
  listGRNs(): Observable<GRN[]> {
    return this.http.get<{ data: GRN[] }>(`${this.base}/grn`).pipe(map(r => r.data));
  }

  getGRN(id: string): Observable<GRN> {
    return this.http.get<GRN>(`${this.base}/grn/${id}`);
  }

  createGRN(body: object): Observable<GRN> {
    return this.http.post<GRN>(`${this.base}/grn`, body);
  }

  updateGRN(id: string, body: object): Observable<GRN> {
    return this.http.put<GRN>(`${this.base}/grn/${id}`, body);
  }

  confirmGRN(id: string): Observable<GRN> {
    return this.http.post<GRN>(`${this.base}/grn/${id}/confirm`, {});
  }

  cancelGRN(id: string): Observable<void> {
    return this.http.post<void>(`${this.base}/grn/${id}/cancel`, {});
  }

  // Stock
  listStockSummary(): Observable<StockSummary[]> {
    return this.http.get<{ data: StockSummary[] }>(`${this.base}/stock/summary`).pipe(map(r => r.data));
  }

  listStockLots(rawMaterialId?: string): Observable<StockLot[]> {
    const params = rawMaterialId ? `?raw_material_id=${rawMaterialId}` : '';
    return this.http.get<{ data: StockLot[] }>(`${this.base}/stock/lots${params}`).pipe(map(r => r.data));
  }

  getStockLotDetail(id: string): Observable<StockLotDetail> {
    return this.http.get<StockLotDetail>(`${this.base}/stock/lots/${id}`);
  }

  // Adjustments
  listAdjustments(): Observable<Adjustment[]> {
    return this.http.get<{ data: Adjustment[] }>(`${this.base}/adjustments`).pipe(map(r => r.data));
  }

  getAdjustment(id: string): Observable<Adjustment> {
    return this.http.get<Adjustment>(`${this.base}/adjustments/${id}`);
  }

  createAdjustment(body: object): Observable<Adjustment> {
    return this.http.post<Adjustment>(`${this.base}/adjustments`, body);
  }

  approveAdjustment(id: string): Observable<Adjustment> {
    return this.http.post<Adjustment>(`${this.base}/adjustments/${id}/approve`, {});
  }

  cancelAdjustment(id: string): Observable<void> {
    return this.http.post<void>(`${this.base}/adjustments/${id}/cancel`, {});
  }
}
