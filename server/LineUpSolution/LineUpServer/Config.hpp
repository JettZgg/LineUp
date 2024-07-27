#pragma once

#include <boost/json.hpp>
#include <string>
#include <fstream>
#include <iostream>
#include <cstdint>

class Config {
public:
    static void load(const std::string& filename);
    static uint32_t getUInt32(const std::string& key, uint32_t defaultValue);
    static std::string getString(const std::string& key, const std::string& defaultValue);

private:
    static boost::json::object config;
};