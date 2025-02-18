import { Component} from '@angular/core';
import { Router } from '@angular/router';
import { MsalBroadcastService, MsalService } from '@azure/msal-angular';
import { AuthenticationResult, EventMessage, EventType, InteractionStatus } from '@azure/msal-browser';
import { Subject, async, filter, takeUntil } from 'rxjs';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  standalone: false,
  styleUrl: './app.component.css'
})
export class AppComponent {
  private readonly _destroying$ = new Subject<void>();

  loginStatus: boolean = false;
  constructor(private broadcastService: MsalBroadcastService, private msalService: MsalService, private route: Router) {
    this.broadcastService.msalSubject$
      .pipe(
        filter((msg: EventMessage) => msg.eventType === EventType.ACCOUNT_ADDED || msg.eventType === EventType.ACCOUNT_REMOVED),
      )
      .subscribe((result: EventMessage) => {
        if (this.msalService.instance.getAllAccounts().length === 0) {
          // Redirect to Login page
          this.route.navigate(['/login']);
        } else {
          this.setLoginDisplay();
        }
      });

    this.broadcastService.inProgress$
      .pipe(
        filter((status: InteractionStatus) => status === InteractionStatus.None),
        takeUntil(this._destroying$)
      )
      .subscribe(() => {
        this.setLoginDisplay();
      })
  }

  setLoginDisplay() {
    this.loginStatus = this.msalService.instance.getAllAccounts().length > 0;

    if (this.loginStatus) {
      const profile_data = this.msalService.instance.getAllAccounts()[0];
      this.msalService.instance.setActiveAccount(profile_data);
      this.route.navigate(['/']);

    } else {
      console.log("Not logged-in");
      // Redirect to Login page
      this.route.navigate(['/login']);
    }

  }

  ngOnDestroy(): void {
    this._destroying$.next(undefined);
    this._destroying$.complete();
  }
 
}
