// src/app/auth-config.ts
import { Configuration, InteractionType, LogLevel } from '@azure/msal-browser';
import { MsalGuardConfiguration, MsalInterceptorConfiguration} from '@azure/msal-angular';

export const msalConfig: Configuration = {
  auth: {
    clientId: '1f7ce7f6-ac7c-43f5-889a-50960ebc58c7', // Application (client) ID from Azure AD
    authority: 'https://login.microsoftonline.com/3192a717-1c36-4a32-b40f-d91972b86f32', // Directory (tenant) ID
    redirectUri: 'https://localhost:53482/', // Your application's redirect URI
  },
  cache: {
    cacheLocation: 'localStorage', // This configures where your cache will be stored
    storeAuthStateInCookie: false, // Set this to true if you are having issues on IE11 or Edge
  },
  system: {
    loggerOptions: {
      loggerCallback: (level, message, containsPii) => {
        if (containsPii) {
          return;
        }
        switch (level) {
          case LogLevel.Error:
            console.error(message);
            return;
          case LogLevel.Info:
            console.info(message);
            return;
          case LogLevel.Verbose:
            console.debug(message);
            return;
          case LogLevel.Warning:
            console.warn(message);
            return;
        }
      },
    },
  },
};

export const msalGuardConfig: MsalGuardConfiguration = {
  interactionType: InteractionType.Redirect,
  authRequest: {
    scopes: ['user.read'],
  },
};

export const msalInterceptorConfig: MsalInterceptorConfiguration = {
  interactionType: InteractionType.Redirect,
  protectedResourceMap: new Map([
    ['https://graph.microsoft.com/v1.0/me', ['user.read']],
  ]),
};
