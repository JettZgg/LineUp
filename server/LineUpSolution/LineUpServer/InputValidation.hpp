#pragma once

#include <string>
#include <regex>

namespace InputValidation {
    bool isValidEmail(const std::string& email);
    bool isValidUsername(const std::string& username);
    bool isValidPassword(const std::string& password);
    std::string sanitizeInput(const std::string& input);
}