# LineUp
A custom Gomoku

## Key Features

1.  **Customizable Game Settings**
2.  **Game Server**
3.  **Client-Side Application**
4.  **Data Collection**
5.  **Quantitative Analysis Module**
6.  **Dashboard**

## Game Rules

-   **Winning Condition**: Align a specified number of pieces in a row (3-19).
-   **Board Size**: Width and height range from 3 to 99.
-   **Minimum Board Size**: At least as large as the number of pieces needed to win.
-   **First to Win**: The first player to align the required number of pieces wins.
-   **Tie Condition**: Board is filled and no player has won.

## Tech Stack

### Server-Side (Go)

-   **Language**: Go
-   **Framework**: Gin for routing, Go WebSocket for real-time updates.
-   **Configuration**: Viper for configuration management.
-   **Testing**: Go testing tools.
-   **IDE**: Visual Studio Code

### Client-Side (React)

-   **Language**: JavaScript
-   **Framework**: React
-   **State Management**: Redux/Context API
-   **Styling**: Styled-Components/CSS Modules
-   **Real-Time Communication**: Socket.IO-client
-   **Bundler**: Vite

### Data Analysis (Python)

-   **Language**: Python
-   **Libraries**: Pandas, NumPy, Matplotlib, Seaborn, SciPy

### Database

-   **Database Management System**: PostgreSQL
-   **ORM**: GORM (Go) or SQLAlchemy (Python)

### Deployment

-   **Hosting**: AWS
-   **Web Server**: Nginx
-   **CI/CD**: GitHub Actions
-   **Domain Registration**: GoDaddy
-   **SSL Certificates**: Let's Encrypt

### Additional Tools and Services

-   **Version Control**: Git/GitHub
-   **Containerization**: Docker
-   **Monitoring and Logging**: Prometheus, Grafana, ELK Stack
-   **Testing**: Postman

## How to run
### Server
```bash
cd server
go mod tidy # only run this when it's first time
go run cmd/server/main.go
```

### client
```bash
cd client
npm run dev
```