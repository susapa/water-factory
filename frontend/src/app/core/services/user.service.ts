import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map } from 'rxjs';
import { environment } from '../../../environments/environment';
import { User, Role } from '../models/user.model';

@Injectable({ providedIn: 'root' })
export class UserService {
  private http = inject(HttpClient);
  private base = `${environment.apiUrl}`;

  listUsers(): Observable<User[]> {
    return this.http.get<{ data: User[] }>(`${this.base}/users`).pipe(map(r => r.data));
  }

  createUser(body: { email: string; password: string; full_name: string; role_id: number }): Observable<User> {
    return this.http.post<User>(`${this.base}/users`, body);
  }

  setUserActive(id: string, isActive: boolean): Observable<void> {
    return this.http.put<void>(`${this.base}/users/${id}/active`, { is_active: isActive });
  }

  listRoles(): Observable<Role[]> {
    return this.http.get<{ data: Role[] }>(`${this.base}/roles`).pipe(map(r => r.data));
  }
}
