import { Inject, Injectable } from '@angular/core';
import { MSAL_GUARD_CONFIG, MsalGuardConfiguration, MsalService } from '@azure/msal-angular';
import { AuthenticationResult, PopupRequest, RedirectRequest } from '@azure/msal-browser';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  profile: any;
  constructor(@Inject(MSAL_GUARD_CONFIG) private msalGuardConfig: MsalGuardConfiguration, private msalService: MsalService) { }

  getUserDetails() {
    const profileData = this.msalService.instance.getAllAccounts()[0];
    this.profile = { name: profileData.name, mail: profileData.username, token: profileData.idToken }
    return this.profile;
  }

  microsoftLogin() {
    if (this.msalGuardConfig.authRequest) {
      this.msalService.loginPopup({ ...this.msalGuardConfig.authRequest } as PopupRequest)
        .subscribe((response: AuthenticationResult) => {
          this.msalService.instance.setActiveAccount(response.account);
        });
    } else {
      this.msalService.loginPopup()
        .subscribe((response: AuthenticationResult) => {
          this.msalService.instance.setActiveAccount(response.account);
        });
    }
  }

  microsoftLogout() {
    this.msalService.logoutRedirect({
      postLogoutRedirectUri: location.origin
    });
  }
}
