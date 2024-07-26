#pragma once

#include <vector>
#include <cstdint>
#include <string>

class Game {
public:
    enum class MoveResult {
        Success,
        OutOfBounds,
        PositionOccupied,
        NotPlayersTurn,
        GameAlreadyOver
    };

    Game(uint32_t width, uint32_t height, uint32_t win_condition);

    MoveResult makeMove(uint32_t x, uint32_t y, uint32_t player);
    bool checkWin(uint32_t x, uint32_t y, uint32_t player) const;
    bool isBoardFull() const;
    bool isGameOver() const;
    uint32_t getWinner() const;
    uint32_t getCurrentPlayer() const;
    std::string getSerializedState() const;
    std::string getLastError() const;

private:
    std::vector<std::vector<uint32_t>> board_;
    uint32_t width_;
    uint32_t height_;
    uint32_t win_condition_;
    uint32_t moves_made_;
    uint32_t current_player_;
    bool game_over_;
    std::string last_error_;
};