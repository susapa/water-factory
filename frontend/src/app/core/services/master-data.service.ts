import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map } from 'rxjs';
import { environment } from '../../../environments/environment';
import {
  UOM, RawMaterialCategory,
  RawMaterial, FinishedGood, BOM, Customer, Supplier,
} from '../models/master-data.model';

@Injectable({ providedIn: 'root' })
export class MasterDataService {
  private http = inject(HttpClient);
  private base = `${environment.apiUrl}/master-data`;

  // Lookups
  listUOM(): Observable<UOM[]> {
    return this.http.get<{ data: UOM[] }>(`${this.base}/uom`).pipe(map(r => r.data));
  }

  listCategories(): Observable<RawMaterialCategory[]> {
    return this.http.get<{ data: RawMaterialCategory[] }>(`${this.base}/raw-material-categories`).pipe(map(r => r.data));
  }

  // Raw Materials
  listRawMaterials(): Observable<RawMaterial[]> {
    return this.http.get<{ data: RawMaterial[] }>(`${this.base}/raw-materials`).pipe(map(r => r.data));
  }

  getRawMaterial(id: string): Observable<RawMaterial> {
    return this.http.get<RawMaterial>(`${this.base}/raw-materials/${id}`);
  }

  createRawMaterial(body: object): Observable<RawMaterial> {
    return this.http.post<RawMaterial>(`${this.base}/raw-materials`, body);
  }

  updateRawMaterial(id: string, body: object): Observable<RawMaterial> {
    return this.http.put<RawMaterial>(`${this.base}/raw-materials/${id}`, body);
  }

  setRawMaterialActive(id: string, isActive: boolean): Observable<void> {
    return this.http.put<void>(`${this.base}/raw-materials/${id}/active`, { is_active: isActive });
  }

  // Finished Goods
  listFinishedGoods(): Observable<FinishedGood[]> {
    return this.http.get<{ data: FinishedGood[] }>(`${this.base}/finished-goods`).pipe(map(r => r.data));
  }

  getFinishedGood(id: string): Observable<FinishedGood> {
    return this.http.get<FinishedGood>(`${this.base}/finished-goods/${id}`);
  }

  createFinishedGood(body: object): Observable<FinishedGood> {
    return this.http.post<FinishedGood>(`${this.base}/finished-goods`, body);
  }

  updateFinishedGood(id: string, body: object): Observable<FinishedGood> {
    return this.http.put<FinishedGood>(`${this.base}/finished-goods/${id}`, body);
  }

  setFinishedGoodActive(id: string, isActive: boolean): Observable<void> {
    return this.http.put<void>(`${this.base}/finished-goods/${id}/active`, { is_active: isActive });
  }

  // BOM
  listBOMs(): Observable<BOM[]> {
    return this.http.get<{ data: BOM[] }>(`${this.base}/bom`).pipe(map(r => r.data));
  }

  getBOM(id: string): Observable<BOM> {
    return this.http.get<BOM>(`${this.base}/bom/${id}`);
  }

  createBOM(body: object): Observable<BOM> {
    return this.http.post<BOM>(`${this.base}/bom`, body);
  }

  updateBOM(id: string, body: object): Observable<BOM> {
    return this.http.put<BOM>(`${this.base}/bom/${id}`, body);
  }

  // Customers
  listCustomers(): Observable<Customer[]> {
    return this.http.get<{ data: Customer[] }>(`${this.base}/customers`).pipe(map(r => r.data));
  }

  getCustomer(id: string): Observable<Customer> {
    return this.http.get<Customer>(`${this.base}/customers/${id}`);
  }

  createCustomer(body: object): Observable<Customer> {
    return this.http.post<Customer>(`${this.base}/customers`, body);
  }

  updateCustomer(id: string, body: object): Observable<Customer> {
    return this.http.put<Customer>(`${this.base}/customers/${id}`, body);
  }

  setCustomerActive(id: string, isActive: boolean): Observable<void> {
    return this.http.put<void>(`${this.base}/customers/${id}/active`, { is_active: isActive });
  }

  // Suppliers
  listSuppliers(): Observable<Supplier[]> {
    return this.http.get<{ data: Supplier[] }>(`${this.base}/suppliers`).pipe(map(r => r.data));
  }

  getSupplier(id: string): Observable<Supplier> {
    return this.http.get<Supplier>(`${this.base}/suppliers/${id}`);
  }

  createSupplier(body: object): Observable<Supplier> {
    return this.http.post<Supplier>(`${this.base}/suppliers`, body);
  }

  updateSupplier(id: string, body: object): Observable<Supplier> {
    return this.http.put<Supplier>(`${this.base}/suppliers/${id}`, body);
  }

  setSupplierActive(id: string, isActive: boolean): Observable<void> {
    return this.http.put<void>(`${this.base}/suppliers/${id}/active`, { is_active: isActive });
  }
}
