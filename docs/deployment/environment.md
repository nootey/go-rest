## Environment

The server uses a .env file to store and read secrets.

This is the desired structure:

```js
# Database Configuration
MONGO_URI=<url>
DATABASE_NAME=go-rest

# Application Configuration
PORT=8080
RELEASE=false
HOST=127.0.0.1

# Security configuration
WEB_CLIENT_DOMAIN=localhost
WEB_CLIENT_PORT=<port>
JWT_WEB_CLIENT_ACCESS=<string>
```
