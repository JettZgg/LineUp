# LineUp
A custom Gomoku

## Key Features
1.  Customizable Game Settings:
    -   Users specify board size (width and height).
    -   Users specify the number of pieces in a row to win.

2.  Game Server:
    -   Handle game creation and player connections.
    -   Manage game state and enforce rules.
    -   Provide real-time updates to players using WebSocket.

3.  Client-Side Application:
    -   User interface for entering game settings and playing the game.
    -   Real-time updates and interactions with the server.

4.  Data Collection:
    -   Collect game data including player moves, game outcomes, and game durations.
    -   Store this data in a database for analysis.

5.  Quantitative Analysis Module:
    -   Analyze collected data to provide insights such as win rates, average game duration, and common winning strategies.

6.  Dashboard:
    -   Web-based dashboard to display analysis results.
    -   Allow users to view performance metrics and game statistics.

## Game Rules
-   **Winning Condition**: Align a specified number of pieces in a row (between 3 and 19) to win.
-   **Board Size**: The board's width and height can range from 3 to 99.
-   **Minimum Board Size**: The board dimensions must be at least as large as the number of pieces needed to win.
-   **First to Win**: The first player to align the required number of consecutive pieces wins.
-   **Tie Condition**: If the board is filled and no player has won, the game ends in a tie.

## Tech Stack

### Server-Side (C++)
-   **Language**: C++
-   **Framework**: Boost.Beast for handling HTTP requests and WebSocket connections.
-   **Build System**: CMake for project configuration and building.
-   **Testing**: Google Test for unit testing.
-   **IDE**: Visual Studio

### Client-Side (React)
-   **Language**: JavaScript
-   **Framework**: React for building the user interface.
-   **State Management**: Redux or Context API for state management.
-   **Styling**: Styled-Components or CSS Modules for styling.
-   **Real-Time Communication**: Socket.IO-client for WebSocket communication with the server.

### Data Analysis (Python)
-   **Language**: Python
-   **Libraries**:
    -   Pandas for data manipulation.
    -   NumPy for numerical operations.
    -   Matplotlib and Seaborn for data visualization.
    -   SciPy for advanced statistical analysis.

### Database
-   **Database Management System**: PostgreSQL for storing game data.
-   **ORM**: SQLAlchemy or Psycopg2 in Python for database interactions.

### Deployment
-   **Hosting**: AWS (Amazon Web Services) for hosting the application.
-   **Web Server**: Nginx for serving the application.
-   **CI/CD**: GitHub Actions for continuous integration and deployment.
-   **Domain Registration**: GoDaddy for domain registration.
-   **SSL Certificates**: Let's Encrypt for securing the application with HTTPS.

### Additional Tools and Services
-   **Version Control**: Git for version control, with GitHub for repository hosting.
-   **Containerization**: Docker for containerizing applications and ensuring consistency across environments.
-   **Monitoring and Logging**: Prometheus and Grafana for monitoring, ELK Stack (Elasticsearch, Logstash, Kibana) for logging and log analysis.
-   **Testing**: Postman for API testing.
