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
1.  To win, a player must align a specified number of pieces in a row, which can be between 3 and 19.
2.  The board's width and height can range from 3 to 99.
3.  The dimensions of the board (width and height) must be at least as large as the number of pieces needed in a row to win.
4.  The player who first achieves the required number of consecutive pieces wins the game.
5.  If the board is completely filled and no player has won, the game ends in a tie.
