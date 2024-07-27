#ifndef CONFIG_HPP
#define CONFIG_HPP

#include <string>
#include <boost/json.hpp>

class Config {
public:
    static boost::json::object config;

    static void load(const std::string& filename);
    static uint32_t getUInt32(const std::string& key, uint32_t defaultValue);
    static std::string getString(const std::string& key, const std::string& defaultValue);
};

std::string findConfigFile();

#endif // CONFIG_HPP
