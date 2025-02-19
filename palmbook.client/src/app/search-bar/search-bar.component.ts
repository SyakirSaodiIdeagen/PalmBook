import { AfterViewInit, Component, ElementRef, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { MsalService } from '@azure/msal-angular';
import { AuthenticationResult } from '@azure/msal-browser';
import { AuthService } from '../services/auth/auth.service';
import {SearchService} from "./search.service";
import { MatTableDataSource } from '@angular/material/table';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { Subject, catchError, debounceTime, distinctUntilChanged, of, switchMap } from 'rxjs';

@Component({
  selector: 'app-search-bar',
  standalone: false,
  templateUrl: './search-bar.component.html',
  styleUrl: './search-bar.component.css'
})
export class SearchBarComponent implements AfterViewInit {
  query: string = '';
  results: { name: string; type: string; source: string; downloadUrl: string }[] = [];
  loading = false;
  searched = false;
  loggedIn = false;
  dataSource = new MatTableDataSource([]);
  displayedColumns: string[] = ['name', 'source', 'createDate', 'updateDate', 'action'];
  private searchSubject = new Subject<string>();

  @ViewChild(MatPaginator) paginator!: MatPaginator;
  @ViewChild(MatSort) sort!: MatSort;
  @ViewChild('topElement') topElement!: ElementRef;
  constructor(private authService: AuthService, private router: Router,
              private _searchService: SearchService) { }
  ngOnInit() {
    const user = this.authService.getUserDetails();
    console.log(user);
    this.loggedIn = !!user;

    if (!this.loggedIn) {
      this.router.navigate(['/login']);
    }

    this.searchSubject
      .pipe(
        debounceTime(500),
        distinctUntilChanged(),
        switchMap(query => {
          if (!query.trim()) {
            this.loading = false;
            this.searched = false;
            return of([]); // ✅ Fix: Return an observable
          }
          this.loading = true;
          this.searched = true;
          return this._searchService.search(query).pipe(
            catchError(() => {
              this.loading = false;
              return of([]); // ✅ Handle API error
            })
          );
        })
      )
      .subscribe(result => {
        console.log(result, "Result");
        this.results = result;
        this.dataSource.data = result;
        this.loading = false;

        setTimeout(() => this.scrollToTop(), 0);
      });
  }

  ngAfterViewInit() {
    this.dataSource.paginator = this.paginator;
    this.dataSource.sort = this.sort;
  }

  onSearchInput() {
    this.searchSubject.next(this.query);
  }

  scrollToTop() {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }

  openFile(url: string): void {
    window.open(url, '_blank', 'noopener,noreferrer');
  }

  logout() {
    this.authService.microsoftLogout();
  }
}
