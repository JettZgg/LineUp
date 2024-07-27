#include "Config.hpp"
#include <limits>
#include <iostream>
#include <fstream>
#include <stdexcept>
#include <filesystem> // Ensure this is included correctly
#include <boost/json.hpp>

boost::json::object Config::config;

void Config::load(const std::string& filename) {
    std::ifstream file(filename);
    if (!file.is_open()) {
        std::cerr << "Failed to open config file: " << filename << std::endl;
        throw std::runtime_error("Config file not found");
    }

    std::string content((std::istreambuf_iterator<char>(file)),
        (std::istreambuf_iterator<char>()));

    try {
        config = boost::json::parse(content).as_object();
    }
    catch (const boost::system::system_error& e) {
        std::cerr << "Failed to parse config file: " << e.what() << std::endl;
    }
}

uint32_t Config::getUInt32(const std::string& key, uint32_t defaultValue) {
    auto it = config.find(key);
    if (it != config.end()) {
        if (it->value().is_int64()) {
            int64_t value = it->value().as_int64();
            if (value >= 0 && value <= std::numeric_limits<uint32_t>::max()) {
                return static_cast<uint32_t>(value);
            }
        }
        else if (it->value().is_uint64()) {
            uint64_t value = it->value().as_uint64();
            if (value <= std::numeric_limits<uint32_t>::max()) {
                return static_cast<uint32_t>(value);
            }
        }
        std::cerr << "Warning: Config value for '" << key << "' is out of range for uint32_t. Using default." << std::endl;
    }
    return defaultValue;
}

std::string Config::getString(const std::string& key, const std::string& defaultValue) {
    auto it = config.find(key);
    if (it != config.end() && it->value().is_string()) {
        return it->value().as_string().c_str();
    }
    return defaultValue;
}

std::string findConfigFile() {
    std::filesystem::path currentPath = std::filesystem::current_path();
    std::filesystem::path configPath = currentPath / "config.json";

    if (!std::filesystem::exists(configPath)) {
        configPath = currentPath.parent_path() / "config.json";
    }

    if (std::filesystem::exists(configPath)) {
        return configPath.string();
    }
    else {
        throw std::runtime_error("Config file not found");
    }
}
