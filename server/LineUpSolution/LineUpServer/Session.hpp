#pragma once

#include <string>
#include <chrono>

class Session {
public:
    Session(const std::string& userEmail);
    bool isExpired() const;
    void refresh();
    const std::string& getUserEmail() const;

private:
    std::string userEmail;
    std::chrono::system_clock::time_point expirationTime;
    static const std::chrono::minutes SESSION_DURATION;
};