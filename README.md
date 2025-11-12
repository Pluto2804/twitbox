TwitBox
Welcome to Twitbox, a modern micro-social platform built for creators, thinkers, and storytellers. 
We are focused on providing a simple and meaningful space that values authenticity and expression. 
Our platform is engineered for lightning-fast performance using Go and is secure by design, featuring robust session authentication, CSRF protection, and a privacy-first architecture. 
It demonstrates clean architecture, secure authentication, and efficient session management — all implemented from scratch.



User Authentication & Authorization: Secure user registration and login system
Session Management: HTTP-only cookie-based sessions with CSRF protection
Twit Management: Create, view, and manage short messages
Template Caching: Efficient HTML template rendering with in-memory caching
Middleware Chain: Composable middleware for logging, recovery, authentication, and CSRF protection
Database Migrations: Version-controlled database schema management

https://github.com/user-attachments/assets/1a000400-972a-4e76-98ad-47cf9b53843c


Production Deployment: Deployed on DigitalOcean with systemd service management

Architecture
MVC Pattern
The application follows a clean MVC (Model-View-Controller) architecture:

Models (internal/model): Data access layer with abstraction and encapsulation
Views (ui/html): Template-based HTML rendering
Controllers (cmd/web): HTTP handlers and request processing

Key Design Principles

Polymorphism: Interface-based database access for testability
Abstraction: Clean separation between data layer and business logic
Encapsulation: Internal packages hide implementation details

Project Structure
```bash
twitbox/
│
├── cmd/web/                 # Application entry point and HTTP layer
│   ├── main.go              # Application initialization
│   ├── handlers.go          # HTTP request handlers
│   ├── middleware.go        # HTTP middleware chain
│   ├── routes.go            # URL routing configuration
│   ├── helpers.go           # Helper functions
│   ├── templates.go         # Template management
│   └── context.go           # Request context handling
│
├── internal/                # Private application code
│   ├── model/               # Data models and database logic
│   │   ├── twits.go         # Twit data access
│   │   └── errors.go        # Custom error types
│   └── validator/           # Input validation logic
│       └── validator.go
│
├── ui/                      # User interface assets
│   ├── html/                # HTML templates
│   │   ├── base.tmpl.html   # Base template layout
│   │   ├── pages/           # Page templates
│   │   │   ├── home.tmpl.html
│   │   │   ├── create.tmpl.html
│   │   │   ├── view.tmpl.html
│   │   │   ├── login.tmpl.html
│   │   │   ├── signup.tmpl.html
│   │   │   ├── account.tmpl.html
│   │   │   <img width="1512" height="982" alt="Screenshot 2025-11-12 at 10 03 02 AM" src="https://github.com/user-attachments/assets/9bc44c7e-b05c-4c14-8f61-791496e0153a" />
<img width="1512" height="982" alt="Screenshot 2025-11-12 at 9 52 44 AM" src="https://github.com/user-attachments/assets/d68c0db4-edb7-4c44-8107-a15a622c279d" />
<img width="1512" height="982" alt="Screenshot 2025-11-12 at 9 42 53 AM" src="https://github.com/user-attachments/assets/befba7e2-c6c6-4bbb-b810-a203151bbfe6" />
└── about.tmpl.html
│   │   └── partials/        # Reusable components
│   │       └── nav.tmpl.html
│   └── static/              # Static assets
│       ├── css/
│       │   └── main.css
│       ├── js/
│       │   └── main.js
│       └── img/
│           ├── favicon.ico
│           └── logo.png
│
├── migrations/              # Database migration files
│   ├── 001_initial.up.sql
│   └── 001_initial.down.sql
│
├── remote/setup/            # Deployment shell scripts
│   ├── 01.sh                # Initial system setup
│   └── 02.sh                # Application setup
│
├── bin/                     # Compiled binaries
│   └── linux_amd64/
│       └── web
│
├── go.mod                   # Go module dependencies
├── go.sum                   # Dependency checksums
├── Makefile                 # Build and deployment automation
└── LICENSE                  # Project license


Technology Stack

Language: Go 1.21+
Database: MySQL 8.0+
Web Server: Caddy (reverse proxy with automatic HTTPS)
Session Store: Encrypted cookie-based sessions
Deployment: DigitalOcean Droplet (Ubuntu)
Process Management: systemd service

Prerequisites

Go 1.21 or higher
MySQL 8.0 or higher
Make (for using Makefile commands)
