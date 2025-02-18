import { Component, ElementRef, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { MsalService } from '@azure/msal-angular';
import { AuthenticationResult } from '@azure/msal-browser';
import { AuthService } from '../services/auth/auth.service';
import {SearchService} from "./search.service";

@Component({
  selector: 'app-search-bar',
  standalone: false,
  templateUrl: './search-bar.component.html',
  styleUrl: './search-bar.component.css'
})
export class SearchBarComponent {
  query: string = '';
  results: { name: string; type: string; source: string; downloadUrl: string }[] = [];
  loading = false;
  searched = false;
  loggedIn = false;

  @ViewChild('topElement') topElement!: ElementRef;

  constructor(private authService: AuthService, private router: Router,
              private _searchService: SearchService) { }
  ngOnInit() {
    // // Check if the user is logged in and update the loggedIn status
    // var user = this.authService.getUserDetails();
    // console.log(user);
    // this.loggedIn = !!this.authService.getUserDetails();
    //
    // if (!this.loggedIn) {
    //   // Redirect to login if not logged in
    //   this.router.navigate(['/login']);
    // }  
  }
  search() {
    if (!this.query.trim()) return;

    this.loading = true;
    this.results = [];
    this.searched = true;

    this._searchService.search(this.query).subscribe(result => {
      console.log('resultsss',result);
    })

    // setTimeout(() => {
    //   this.loading = false;
    //   this.results =
    //     this.query.toLowerCase() === 'test'
    //       ? [
    //         { name: 'Report.pdf', type: 'PDF', source: 'SharePoint', downloadUrl: '#' },
    //         { name: 'Data.xlsx', type: 'Excel', source: 'Teams', downloadUrl: '#' },
    //         { name: 'Presentation.pptx', type: 'PowerPoint', source: 'OneDrive', downloadUrl: '#' }
    //       ]
    //       : [];
    //
    //   this.scrollToTop();
    // }, 1500);
  }

  scrollToTop() {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }

  logout() {
    this.authService.microsoftLogout();
  }
}
