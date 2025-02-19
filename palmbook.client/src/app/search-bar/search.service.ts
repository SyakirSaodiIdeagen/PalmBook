import { Injectable } from '@angular/core';
import {HttpClient, HttpHeaders, HttpParams} from "@angular/common/http";
import {Observable} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class SearchService {

  private apiUrl = 'http://localhost:5555/search';
  private user: any;  

  constructor(private http: HttpClient) {}

  search(query: string): Observable<any> {
    debugger
    const user:any = localStorage.getItem('user')
    this.user = JSON.parse(user);
    const headers = new HttpHeaders({
      'Authorization': `${(this.user.token)}`,
      'Content-Type': 'application/json'
    });
    console.log(query);
    const params = new HttpParams().set('query', query);  

    return this.http.get<any>(this.apiUrl, { params, headers });
  }
}
