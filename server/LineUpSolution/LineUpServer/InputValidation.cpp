#include "InputValidation.hpp"

namespace InputValidation {

    bool isValidEmail(const std::string& email) {
        const std::regex pattern("(\\w+)(\\.|_)?(\\w*)@(\\w+)(\\.(\\w+))+");
        return std::regex_match(email, pattern);
    }

    bool isValidUsername(const std::string& username) {
        const std::regex pattern("^[a-zA-Z0-9_]{3,20}$");
        return std::regex_match(username, pattern);
    }

    bool isValidPassword(const std::string& password) {
        // At least 8 characters, 1 uppercase, 1 lowercase, 1 number
        const std::regex pattern("^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d).{8,}$");
        return std::regex_match(password, pattern);
    }

    std::string sanitizeInput(const std::string& input) {
        std::string sanitized;
        for (char c : input) {
            if (c == '&') sanitized += "&amp;";
            else if (c == '<') sanitized += "&lt;";
            else if (c == '>') sanitized += "&gt;";
            else if (c == '"') sanitized += "&quot;";
            else if (c == '\'') sanitized += "&#x27;";
            else if (c == '/') sanitized += "&#x2F;";
            else sanitized += c;
        }
        return sanitized;
    }

}