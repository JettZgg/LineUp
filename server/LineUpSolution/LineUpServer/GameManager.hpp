#pragma once

#include "Game.hpp"
#include "UserManager.hpp"
#include <map>
#include <string>
#include <memory>
#include <array>
#include <cstdint>

class GameManager {
public:
    enum class GameError {
        Success,
        GameNotFound,
        PlayerNotInGame,
        MoveFailed,
        PlayerNotAuthenticated
    };

    GameManager(UserManager& userManager);

    std::string createGame(const std::string& player1Email, uint32_t width, uint32_t height, uint32_t win_condition);
    bool joinGame(const std::string& gameId, const std::string& player2Email);
    GameError makeMove(const std::string& gameId, const std::string& playerEmail, uint32_t x, uint32_t y);
    std::string getGameState(const std::string& gameId);
    bool isGameOver(const std::string& gameId);
    std::string getCurrentPlayerEmail(const std::string& gameId);
    void handleDisconnection(const std::string& gameId, const std::string& playerEmail);
    std::string getLastGameError(const std::string& gameId);

private:
    std::map<std::string, std::unique_ptr<Game>> games_;
    std::map<std::string, std::array<std::string, 2>> players_;
    UserManager& userManager_;
};