Simple implementation of github.com/pion/turn server, auth done via JWT token.
There is 3 parameters:

- **Realm** - realm for the server
- **Port** - port for the server
- **jwt-sign** - signature for JWT token

Also it runs TCP healthcheck on the same port so can use it with AWS cloud
