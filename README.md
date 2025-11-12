# 🐦 Twitbox – A Go Web Application
Welcome to Twitbox, a modern micro-social platform built for creators, thinkers, and storytellers. 
We are focused on providing a simple and meaningful space that values authenticity and expression. 
Our platform is engineered for lightning-fast performance using Go and is secure by design, featuring robust session authentication, CSRF protection, and a privacy-first architecture. 
It demonstrates clean architecture, secure authentication, and efficient session management — all implemented from scratch.



- 🔐  User Authentication & Authorization: Secure user registration and login system
- ⚙️  Session Management: HTTP-only cookie-based sessions with CSRF protection
- 🧠  Twit Management: Create, view, and manage short messages
- 🧩  Template Caching: Efficient HTML template rendering with in-memory caching
- 🧱  Middleware Chain: Composable middleware for logging, recovery, authentication, and CSRF protection
- 🗄️  Database Migrations: Version-controlled database schema management

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
│   │   │  
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
```


Prerequisites

Go 1.21 or higher
MySQL 8.0 or higher
Make (for using Makefile commands)

Demo Video:
https://github.com/user-attachments/assets/1a000400-972a-4e76-98ad-47cf9b53843c


<img width="1512" height="982" alt="Screenshot 2025-11-12 at 11 23 04 AM" src="https://github.com/user-attachments/assets/6a622563-60b8-4320-96ad-af651d8e183f" />
<img width="1512" height="982" alt="Screenshot 2025-11-12 at 9 41 57 AM" src="https://github.com/user-attachments/assets/cf40371d-aea3-4a12-a992-bbc88b7e29ab" />
<img width="1512" height="982" alt="Screenshot 2025-11-12 at 9 42 53 AM" src="https://github.com/user-attachments/assets/d5477c7d-afcc-403e-8f77-9647afc4327f" />
<img width="1512" height="982" alt="Screenshot 2025-11-12 at 9 52 44 AM" src="https://github.com/user-attachments/assets/d54e8f0d-a8e5-41b6-b358-437c7f20620f" />
<img width="1512" height="982" alt="Screenshot 2025-11-12 at 10 03 02 AM" src="https://github.com/user-attachments/assets/b072add6-6ec6-46f9-ad8b-42f0c1c3422b" />

🔗 **Live App:** 
[https://twitbox.app](https://twitbox.app) 
