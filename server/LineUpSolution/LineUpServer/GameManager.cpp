#include "GameManager.hpp"
#include <random>
#include <sstream>
#include <cstdint>

GameManager::GameManager(UserManager& userManager) : userManager_(userManager) {}

std::string GameManager::createGame(const std::string& player1Email, uint32_t width, uint32_t height, uint32_t win_condition)
{
    User* player1 = userManager_.getUser(player1Email);
    if (!player1) return "";

    static std::random_device rd;
    static std::mt19937 gen(rd());
    static std::uniform_int_distribution<> dis(0, 15);
    static const char* hex_digits = "0123456789ABCDEF";

    std::string game_id;
    do {
        game_id.clear();
        for (int i = 0; i < 8; ++i) {
            game_id += hex_digits[dis(gen)];
        }
    } while (games_.find(game_id) != games_.end());

    games_[game_id] = std::make_unique<Game>(width, height, win_condition);
    players_[game_id] = { player1Email, "" };
    return game_id;
}

bool GameManager::joinGame(const std::string& gameId, const std::string& player2Email)
{
    User* player2 = userManager_.getUser(player2Email);
    if (!player2) return false;

    auto it = players_.find(gameId);
    if (it == players_.end()) return false;

    if (it->second[1].empty() && it->second[0] != player2Email) {
        it->second[1] = player2Email;
        return true;
    }
    return false;
}

GameManager::GameError GameManager::makeMove(const std::string& gameId, const std::string& playerEmail, uint32_t x, uint32_t y)
{
    User* player = userManager_.getUser(playerEmail);
    if (!player) return GameError::PlayerNotAuthenticated;

    auto game_it = games_.find(gameId);
    auto player_it = players_.find(gameId);
    if (game_it == games_.end() || player_it == players_.end())
        return GameError::GameNotFound;

    if (player_it->second[0] != playerEmail && player_it->second[1] != playerEmail)
        return GameError::PlayerNotInGame;

    uint32_t player_number = (player_it->second[0] == playerEmail) ? 1 : 2;
    Game::MoveResult result = game_it->second->makeMove(x, y, player_number);

    if (result == Game::MoveResult::Success && game_it->second->isGameOver()) {
        // Update ranks
        std::string opponent_email = (player_it->second[0] == playerEmail) ? player_it->second[1] : player_it->second[0];
        User* opponent = userManager_.getUser(opponent_email);
        if (opponent) {
            uint32_t winner = game_it->second->getWinner();
            bool player_won = (winner == player_number);
            bool draw = (winner == 0);
            if (!draw) {
                userManager_.updateUserRank(playerEmail, opponent->rank, player_won);
                userManager_.updateUserRank(opponent_email, player->rank, !player_won);
            }
        }
    }

    return (result == Game::MoveResult::Success) ? GameError::Success : GameError::MoveFailed;
}

std::string GameManager::getGameState(const std::string& gameId)
{
    auto it = games_.find(gameId);
    if (it == games_.end()) return "";

    return it->second->getSerializedState();
}

bool GameManager::isGameOver(const std::string& gameId)
{
    auto it = games_.find(gameId);
    if (it == games_.end()) return true;

    return it->second->isGameOver();
}

std::string GameManager::getCurrentPlayerEmail(const std::string& gameId)
{
    auto game_it = games_.find(gameId);
    auto player_it = players_.find(gameId);
    if (game_it == games_.end() || player_it == players_.end()) return "";

    uint32_t current_player = game_it->second->getCurrentPlayer();
    return player_it->second[current_player - 1];
}

void GameManager::handleDisconnection(const std::string& gameId, const std::string& playerEmail)
{
    auto player_it = players_.find(gameId);
    if (player_it != players_.end()) {
        if (player_it->second[0] == playerEmail)
            player_it->second[0] = "";
        else if (player_it->second[1] == playerEmail)
            player_it->second[1] = "";

        if (player_it->second[0].empty() && player_it->second[1].empty()) {
            games_.erase(gameId);
            players_.erase(player_it);
        }
    }
}

std::string GameManager::getLastGameError(const std::string& gameId)
{
    auto it = games_.find(gameId);
    if (it == games_.end())
        return "Game not found";

    return it->second->getLastError();
}