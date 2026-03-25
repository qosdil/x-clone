# Like-X Common HTTP Module

This module provides common HTTP utilities for the Like-X backend, including header-based authentication middleware.

## Overview

The `backend/common/http` module is part of the Like-X project (a Twitter-like social media platform). It provides reusable HTTP components built on top of the [Fiber](https://gofiber.io/) web framework.

## Features

- **Header-Based Authentication Middleware**: Extracts user ID from request headers for authentication

## Dependencies

- [Fiber v3](https://github.com/gofiber/fiber) - Web framework

## Usage

### Authentication Middleware

```go
import "github.com/qosdil/like-x/backend/common/http/auth"

app.Use(auth.AuthMiddleware)
```

The middleware expects an `Auth-User-ID` header with a numeric user ID:

```
Auth-User-ID: 123
```

On successful validation (non-zero user ID), it sets the following local in the Fiber context:
- `auth_user_id`: The user ID extracted from the header

If the header is missing or contains `0`, the middleware returns a 401 Unauthorized response.

## Contributing

This is part of the Like-X backend monorepo. Please follow the main project's contribution guidelines.