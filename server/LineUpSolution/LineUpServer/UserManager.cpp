#include "UserManager.hpp"
#include "Session.hpp"
#include <openssl/sha.h>
#include <openssl/rand.h>
#include <openssl/bio.h>
#include <openssl/evp.h>
#include <openssl/buffer.h>
#include <cmath>
#include <random>
#include <algorithm>
#include <iomanip>
#include <sstream>
#include <cstring>
#include <cstdint>

User::User(const std::string& email, const std::string& username, const std::string& password)
    : email(email), username(username), password_hash(hashPassword(password)), rank(1000) {}

std::vector<unsigned char> User::generateSalt() {
    std::vector<unsigned char> salt(16);
    RAND_bytes(salt.data(), salt.size());
    return salt;
}

std::string User::hashPassword(const std::string& password) {
    std::vector<unsigned char> salt = generateSalt();
    std::vector<unsigned char> hash(EVP_MAX_MD_SIZE);
    unsigned int hashLen;

    std::vector<unsigned char> passwordData(password.begin(), password.end());
    std::vector<unsigned char> saltedPassword = salt;
    saltedPassword.insert(saltedPassword.end(), passwordData.begin(), passwordData.end());

    EVP_MD_CTX* mdctx = EVP_MD_CTX_new();
    const EVP_MD* md = EVP_sha256();

    EVP_DigestInit_ex(mdctx, md, nullptr);
    EVP_DigestUpdate(mdctx, saltedPassword.data(), saltedPassword.size());
    EVP_DigestFinal_ex(mdctx, hash.data(), &hashLen);
    EVP_MD_CTX_free(mdctx);

    hash.resize(hashLen);

    std::string saltBase64 = base64Encode(salt);
    std::string hashBase64 = base64Encode(hash);

    return saltBase64 + ":" + hashBase64;
}

bool User::verifyPassword(const std::string& password, const std::string& storedHash) {
    size_t colonPos = storedHash.find(':');
    if (colonPos == std::string::npos) {
        return false;
    }

    std::string saltBase64 = storedHash.substr(0, colonPos);
    std::string hashBase64 = storedHash.substr(colonPos + 1);

    std::vector<unsigned char> salt = base64Decode(saltBase64);
    std::vector<unsigned char> storedHashBytes = base64Decode(hashBase64);

    std::vector<unsigned char> passwordData(password.begin(), password.end());
    std::vector<unsigned char> saltedPassword = salt;
    saltedPassword.insert(saltedPassword.end(), passwordData.begin(), passwordData.end());

    std::vector<unsigned char> computedHash(EVP_MAX_MD_SIZE);
    unsigned int hashLen;

    EVP_MD_CTX* mdctx = EVP_MD_CTX_new();
    const EVP_MD* md = EVP_sha256();

    EVP_DigestInit_ex(mdctx, md, nullptr);
    EVP_DigestUpdate(mdctx, saltedPassword.data(), saltedPassword.size());
    EVP_DigestFinal_ex(mdctx, computedHash.data(), &hashLen);
    EVP_MD_CTX_free(mdctx);

    computedHash.resize(hashLen);

    return computedHash == storedHashBytes;
}

std::string User::base64Encode(const std::vector<unsigned char>& input) {
    BIO* bio, * b64;
    BUF_MEM* bufferPtr;

    b64 = BIO_new(BIO_f_base64());
    bio = BIO_new(BIO_s_mem());
    bio = BIO_push(b64, bio);

    BIO_set_flags(bio, BIO_FLAGS_BASE64_NO_NL);
    BIO_write(bio, input.data(), input.size());
    BIO_flush(bio);
    BIO_get_mem_ptr(bio, &bufferPtr);

    std::string result(bufferPtr->data, bufferPtr->length);
    BIO_free_all(bio);

    return result;
}

std::vector<unsigned char> User::base64Decode(const std::string& input) {
    BIO* bio, * b64;
    std::vector<unsigned char> result(input.size());

    bio = BIO_new_mem_buf(input.c_str(), -1);
    b64 = BIO_new(BIO_f_base64());
    bio = BIO_push(b64, bio);

    BIO_set_flags(bio, BIO_FLAGS_BASE64_NO_NL);
    int decodedSize = BIO_read(bio, result.data(), input.size());
    BIO_free_all(bio);

    result.resize(decodedSize);
    return result;
}

std::string User::getDanRank() const {
    if (rank < 1000) return "Rookie";
    if (rank >= 2000) return "Radiant";
    uint32_t dan = std::min<uint32_t>(10, (rank - 1000) / 100 + 1);
    return std::to_string(dan) + " Dan";
}

void User::updateRank(uint32_t opponent_rank, bool win) {
    int32_t rank_diff = static_cast<int32_t>(opponent_rank) - static_cast<int32_t>(rank);
    int32_t rank_change = 20 + static_cast<int32_t>(std::round(rank_diff * 0.08));
    int32_t new_rank = static_cast<int32_t>(rank) + (win ? rank_change : -rank_change);
    rank = static_cast<uint32_t>(std::max(0, new_rank));
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
    return it->second.password_hash == User::hashPassword(password);
}

User* UserManager::getUser(const std::string& email) {
    auto it = users.find(email);
    if (it == users.end()) {
        return nullptr;
    }
    return &it->second;
}

void UserManager::updateUserRank(const std::string& email, uint32_t opponent_rank, bool win) {
    auto user = getUser(email);
    if (user) {
        user->updateRank(opponent_rank, win);
    }
}

std::string UserManager::createSession(const std::string& email) {
    std::string sessionId = generateSessionId();
    sessions[sessionId] = std::make_unique<Session>(email);
    return sessionId;
}

bool UserManager::validateSession(const std::string& sessionId) {
    auto it = sessions.find(sessionId);
    if (it != sessions.end() && !it->second->isExpired()) {
        it->second->refresh();
        return true;
    }
    return false;
}

void UserManager::refreshSession(const std::string& sessionId) {
    auto it = sessions.find(sessionId);
    if (it != sessions.end()) {
        it->second->refresh();
    }
}

void UserManager::endSession(const std::string& sessionId) {
    sessions.erase(sessionId);
}

void UserManager::cleanExpiredSessions() {
    for (auto it = sessions.begin(); it != sessions.end();) {
        if (it->second->isExpired()) {
            it = sessions.erase(it);
        }
        else {
            ++it;
        }
    }
}

std::string UserManager::generateSessionId() const {
    static const char alphanum[] =
        "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
    std::random_device rd;
    std::mt19937 gen(rd());
    std::uniform_int_distribution<> dis(0, sizeof(alphanum) - 2);

    std::string sessionId;
    sessionId.reserve(32);
    for (int i = 0; i < 32; ++i) {
        sessionId += alphanum[dis(gen)];
    }
    return sessionId;
}

std::string UserManager::getUserEmailFromSession(const std::string& sessionId) {
    auto it = sessions.find(sessionId);
    if (it != sessions.end() && !it->second->isExpired()) {
        return it->second->getUserEmail();
    }
    return ""; // Return empty string if session is not found or expired
}