#include "Session.hpp"

const std::chrono::minutes Session::SESSION_DURATION = std::chrono::minutes(30);

Session::Session(const std::string& userEmail) : userEmail(userEmail) {
    refresh();
}

bool Session::isExpired() const {
    return std::chrono::system_clock::now() > expirationTime;
}

void Session::refresh() {
    expirationTime = std::chrono::system_clock::now() + SESSION_DURATION;
}

const std::string& Session::getUserEmail() const {
    return userEmail;
}