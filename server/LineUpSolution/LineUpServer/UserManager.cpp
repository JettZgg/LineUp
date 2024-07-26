#include "UserManager.hpp"
#include <cmath>
#include <algorithm>
#include <openssl/sha.h>
#include <iomanip>
#include <sstream>

User::User(const std::string& email, const std::string& username, const std::string& password)
    : email(email), username(username), password_hash(hashPassword(password)), rank(1000) {}

std::string User::getDanRank() const {
    if (rank < 1000) return "Rookie";
    if (rank >= 2000) return "Radiant";
    int dan = std::min(10, (rank - 1000) / 100 + 1);
    return std::to_string(dan) + " Dan";
}

void User::updateRank(int opponent_rank, bool win) {
    int rank_diff = opponent_rank - rank;
    int rank_change = 20 + static_cast<int>(std::round(rank_diff * 0.08));
    rank += win ? rank_change : -rank_change;
    rank = std::max(0, rank); // Ensure rank doesn't go below 0
}

std::string User::hashPassword(const std::string& password) const {
    unsigned char hash[SHA256_DIGEST_LENGTH];
    SHA256_CTX sha256;
    SHA256_Init(&sha256);
    SHA256_Update(&sha256, password.c_str(), password.length());
    SHA256_Final(hash, &sha256);
    std::stringstream ss;
    for (int i = 0; i < SHA256_DIGEST_LENGTH; i++) {
        ss << std::hex << std::setw(2) << std::setfill('0') << static_cast<int>(hash[i]);
    }
    return ss.str();
}

bool UserManager::registerUser(const std::string& email, const std::string& username, const std::string& password) {
    if (users.find(email) != users.end()) {
        return false; // User already exists
    }
    users.emplace(email, User(email, username, password));
    return true;
}

bool UserManager::loginUser(const std::string& email, const std::string& password) {
    auto it = users.find(email);
    if (it == users.end()) {
        return false; // User not found
    }
    return it->second.password_hash == it->second.hashPassword(password);
}

User* UserManager::getUser(const std::string& email) {
    auto it = users.find(email);
    if (it == users.end()) {
        return nullptr;
    }
    return &it->second;
}

void UserManager::updateUserRank(const std::string& email, int opponent_rank, bool win) {
    auto user = getUser(email);
    if (user) {
        user->updateRank(opponent_rank, win);
    }
}