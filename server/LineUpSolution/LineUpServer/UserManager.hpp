#pragma once

#include <string>
#include <unordered_map>
#include <cstdint>

class User {
public:
    User(const std::string& email, const std::string& username, const std::string& password);

    std::string email;
    std::string username;
    std::string password_hash;
    int rank;
    std::string getDanRank() const;
    void updateRank(int opponent_rank, bool win);

private:
    std::string hashPassword(const std::string& password) const;
};

class UserManager {
public:
    bool registerUser(const std::string& email, const std::string& username, const std::string& password);
    bool loginUser(const std::string& email, const std::string& password);
    User* getUser(const std::string& email);
    void updateUserRank(const std::string& email, int opponent_rank, bool win);

private:
    std::unordered_map<std::string, User> users;
};