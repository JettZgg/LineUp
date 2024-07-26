#include "Game.hpp"
#include <algorithm>
#include <sstream>

Game::Game(uint32_t width, uint32_t height, uint32_t win_condition)
    : width_(width), height_(height), win_condition_(win_condition),
    moves_made_(0), current_player_(1), game_over_(false)
{
    board_ = std::vector<std::vector<uint32_t>>(height, std::vector<uint32_t>(width, 0));
}

Game::MoveResult Game::makeMove(uint32_t x, uint32_t y, uint32_t player)
{
    if (x >= width_ || y >= height_) {
        last_error_ = "Move is out of bounds";
        return MoveResult::OutOfBounds;
    }
    if (board_[y][x] != 0) {
        last_error_ = "Position is already occupied";
        return MoveResult::PositionOccupied;
    }
    if (game_over_) {
        last_error_ = "Game is already over";
        return MoveResult::GameAlreadyOver;
    }
    if (player != current_player_) {
        last_error_ = "It's not this player's turn";
        return MoveResult::NotPlayersTurn;
    }

    board_[y][x] = player;
    moves_made_++;

    if (checkWin(x, y, player))
        game_over_ = true;
    else if (isBoardFull())
        game_over_ = true;
    else
        current_player_ = (current_player_ == 1) ? 2 : 1;

    return MoveResult::Success;
}

bool Game::checkWin(uint32_t x, uint32_t y, uint32_t player) const
{
    // Horizontal check
    uint32_t count = 0;
    for (int i = std::max(0, static_cast<int>(x - win_condition_ + 1));
        i < std::min(width_, x + win_condition_); i++) {
        count = (board_[y][i] == player) ? count + 1 : 0;
        if (count == win_condition_) return true;
    }

    // Vertical check
    count = 0;
    for (int i = std::max(0, static_cast<int>(y - win_condition_ + 1));
        i < std::min(height_, y + win_condition_); i++) {
        count = (board_[i][x] == player) ? count + 1 : 0;
        if (count == win_condition_) return true;
    }

    // Diagonal checks (top-left to bottom-right and top-right to bottom-left)
    for (int d = 0; d < 2; d++) {
        count = 0;
        for (int i = -static_cast<int>(win_condition_) + 1; i < static_cast<int>(win_condition_); i++) {
            int cx = x + (d == 0 ? i : -i);
            int cy = y + i;
            if (cx >= 0 && cx < width_ && cy >= 0 && cy < height_) {
                count = (board_[cy][cx] == player) ? count + 1 : 0;
                if (count == win_condition_) return true;
            }
        }
    }

    return false;
}

bool Game::isBoardFull() const
{
    return moves_made_ == width_ * height_;
}

bool Game::isGameOver() const
{
    return game_over_;
}

uint32_t Game::getWinner() const
{
    if (!isGameOver()) {
        return 0; // No winner yet
    }

    // Check if the game ended due to a win condition
    for (uint32_t y = 0; y < height_; ++y) {
        for (uint32_t x = 0; x < width_; ++x) {
            if (board_[y][x] != 0) {
                if (checkWin(x, y, board_[y][x])) {
                    return board_[y][x];
                }
            }
        }
    }

    // If the game is over but no one won, it's a draw
    return 0;
}

uint32_t Game::getCurrentPlayer() const
{
    return current_player_;
}

std::string Game::getSerializedState() const
{
    std::stringstream ss;
    ss << width_ << "," << height_ << "," << win_condition_ << ","
        << current_player_ << "," << (game_over_ ? "1" : "0") << ",";
    for (const auto& row : board_) {
        for (uint32_t cell : row) {
            ss << cell << ",";
        }
    }
    return ss.str();
}

std::string Game::getLastError() const
{
    return last_error_;
}