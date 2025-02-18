import { Injectable } from '@angular/core';
import { CanActivate, Router } from '@angular/router';
import { MsalService } from '@azure/msal-angular';

@Injectable({
  providedIn: 'root',
})
export class AuthGuard implements CanActivate {
  constructor(private msalService: MsalService, private router: Router) { }

  canActivate(): boolean {
    // Check if the user is logged in using MSAL
    const isLoggedIn = this.msalService.instance.getAllAccounts().length > 0;

    if (isLoggedIn) {
      // Redirect if the user is already logged in
      this.router.navigate(['/search-bar']);
      return false;
    }
    return true;
  }

}
