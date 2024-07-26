#pragma once

#include "Session.hpp"
#include <string>
#include <unordered_map>
#include <memory>
#include <vector>

class User {
public:
    User(const std::string& email, const std::string& username, const std::string& password);

    std::string email;
    std::string username;
    std::string password_hash;
    int rank;
    std::string getDanRank() const;
    void updateRank(int opponent_rank, bool win);

    static std::string hashPassword(const std::string& password);
    static bool verifyPassword(const std::string& password, const std::string& hash);

private:
    static std::vector<unsigned char> generateSalt();
    static std::string base64Encode(const std::vector<unsigned char>& input);
    static std::vector<unsigned char> base64Decode(const std::string& input);
};

class UserManager {
public:
    bool registerUser(const std::string& email, const std::string& username, const std::string& password);
    bool loginUser(const std::string& email, const std::string& password);
    User* getUser(const std::string& email);
    void updateUserRank(const std::string& email, int opponent_rank, bool win);
    std::string createSession(const std::string& email);
    bool validateSession(const std::string& sessionId);
    void refreshSession(const std::string& sessionId);
    void endSession(const std::string& sessionId);
    void cleanExpiredSessions();
    std::string getUserEmailFromSession(const std::string& sessionId);

private:
    std::unordered_map<std::string, User> users;
    std::unordered_map<std::string, std::unique_ptr<Session>> sessions;
    std::string generateSessionId() const;
};