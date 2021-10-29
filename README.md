# cam-backend-services
## Dev Environment Setup
This module relies on imports from another repository in this private organization, which requires additional config.

- Add an ssh key to your account under `Settings > SSH and GPG keys`
- Add the organization to GOPRIVATE to bypass go mod proxy

    ```go end -w GOPRIVATE=github.com/ASUIFT401ProjectGroup19```
- Update git config to redirect https to ssh

    ```git config --global url."git@github.com:ASUIFT401ProjectGroup19".insteadOf https://github.com/ASUIFT401ProjectGroup19```

## Runtime Config
This application is configured via the environment. The following environment
variables can be used:
```
KEY                    TYPE      DEFAULT    REQUIRED    DESCRIPTION
SERVICE_DB_DRIVER      String    mysql
SERVICE_DB_HOST        String
SERVICE_DB_DATABASE    String
SERVICE_DB_USERNAME    String
SERVICE_DB_PASSWORD    String
SERVICE_PORT           String    10000
```

