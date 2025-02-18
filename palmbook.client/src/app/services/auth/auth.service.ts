import { Inject, Injectable } from '@angular/core';
import { MSAL_GUARD_CONFIG, MsalGuardConfiguration, MsalService } from '@azure/msal-angular';
import { RedirectRequest } from '@azure/msal-browser';

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
      this.msalService.loginRedirect({ ...this.msalGuardConfig.authRequest } as RedirectRequest);
    } else {
      this.msalService.loginRedirect();
    }
  }

  microsoftLogout() {
    this.msalService.logoutRedirect({
      postLogoutRedirectUri: location.origin
    });
  }
}
