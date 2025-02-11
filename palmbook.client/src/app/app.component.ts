import { Component} from '@angular/core';
import { Router } from '@angular/router';
import { MsalService } from '@azure/msal-angular';
import { AuthenticationResult } from '@azure/msal-browser';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  standalone: false,
  styleUrl: './app.component.css'
})
export class AppComponent {
  constructor(private authService: MsalService, private router: Router) { }

  ngOnInit() {
    this.authService.instance.handleRedirectPromise().then((result: AuthenticationResult | null) => {
      if (result !== null && result.account !== null) {
        this.authService.instance.setActiveAccount(result.account);
        this.router.navigate(['/search-bar']);

      }
    });
  }

 
}
