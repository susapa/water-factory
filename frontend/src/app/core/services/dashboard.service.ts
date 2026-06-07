import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map } from 'rxjs';
import { environment } from '../../../environments/environment';
import { DashboardSummary } from '../models/dashboard.model';

@Injectable({ providedIn: 'root' })
export class DashboardService {
  private http = inject(HttpClient);
  private base = `${environment.apiUrl}/dashboard`;

  getDashboardSummary(): Observable<DashboardSummary> {
    return this.http
      .get<{ data: DashboardSummary }>(`${this.base}/summary`)
      .pipe(map(r => r.data));
  }
}
