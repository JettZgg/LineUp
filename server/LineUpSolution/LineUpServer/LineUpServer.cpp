#define _WIN32_WINNT 0x0A00

#include <boost/beast/core.hpp>
#include <boost/beast/websocket.hpp>
#include <boost/asio/ip/tcp.hpp>
#include <boost/json.hpp>
#include <iostream>
#include <string>
#include <thread>
#include <vector>
#include "GameManager.hpp"
#include "UserManager.hpp"
#include "InputValidation.hpp"
#include "Config.hpp"
#include <cstdint>


namespace beast = boost::beast;
namespace websocket = beast::websocket;
namespace net = boost::asio;
namespace json = boost::json;
using tcp = net::ip::tcp;

UserManager user_manager;
GameManager game_manager(user_manager);

void send_error(websocket::stream<tcp::socket>& ws, const std::string& error_message)
{
    json::object response;
    response["action"] = "error";
    response["message"] = error_message;
    ws.write(net::buffer(json::serialize(response)));
}

void handle_message(websocket::stream<tcp::socket>& ws, const std::string& message)
{
    try {
        json::value jv = json::parse(message);
        std::string action = jv.at("action").as_string().c_str();

        if (action == "register") {
            std::string email = InputValidation::sanitizeInput(jv.at("email").as_string().c_str());
            std::string username = InputValidation::sanitizeInput(jv.at("username").as_string().c_str());
            std::string password = jv.at("password").as_string().c_str();

            if (!InputValidation::isValidEmail(email) || !InputValidation::isValidUsername(username) || !InputValidation::isValidPassword(password)) {
                send_error(ws, "Invalid input");
                return;
            }

            bool success = user_manager.registerUser(email, username, password);

            json::object response;
            response["action"] = "register_result";
            response["success"] = success;
            ws.write(net::buffer(json::serialize(response)));
        }
        else if (action == "login") {
            std::string email = InputValidation::sanitizeInput(jv.at("email").as_string().c_str());
            std::string password = jv.at("password").as_string().c_str();

            if (!InputValidation::isValidEmail(email) || !InputValidation::isValidPassword(password)) {
                send_error(ws, "Invalid input");
                return;
            }

            bool success = user_manager.loginUser(email, password);

            json::object response;
            response["action"] = "login_result";
            response["success"] = success;
            if (success) {
                std::string sessionId = user_manager.createSession(email);
                User* user = user_manager.getUser(email);
                response["session_id"] = sessionId;
                response["username"] = user->username;
                response["rank"] = user->rank;
                response["dan_rank"] = user->getDanRank();
            }
            ws.write(net::buffer(json::serialize(response)));
        }
        else if (action == "logout") {
            std::string sessionId = jv.at("session_id").as_string().c_str();
            user_manager.endSession(sessionId);

            json::object response;
            response["action"] = "logout_result";
            response["success"] = true;
            ws.write(net::buffer(json::serialize(response)));
        }
        else if (action == "create_game") {
            std::string sessionId = jv.at("session_id").as_string().c_str();
            if (!user_manager.validateSession(sessionId)) {
                send_error(ws, "Invalid session");
                return;
            }

            uint32_t width = static_cast<uint32_t>(jv.at("width").as_int64());
            uint32_t height = static_cast<uint32_t>(jv.at("height").as_int64());
            uint32_t win_condition = static_cast<uint32_t>(jv.at("win_condition").as_int64());

            if (width < Config::getUInt32("min_board_width", 3) ||
                width > Config::getUInt32("max_board_width", 99) ||
                height < Config::getUInt32("min_board_height", 3) ||
                height > Config::getUInt32("max_board_height", 99) ||
                win_condition < Config::getUInt32("min_win_condition", 3) ||
                win_condition > Config::getUInt32("max_win_condition", 19)) {
                send_error(ws, "Invalid game parameters");
                return;
            }

            std::string email = user_manager.getUserEmailFromSession(sessionId);
            std::string game_id = game_manager.createGame(email, width, height, win_condition);

            json::object response;
            response["action"] = "game_created";
            response["game_id"] = game_id;
            ws.write(net::buffer(json::serialize(response)));
        }
        else if (action == "join_game") {
            std::string sessionId = jv.at("session_id").as_string().c_str();
            if (!user_manager.validateSession(sessionId)) {
                send_error(ws, "Invalid session");
                return;
            }

            std::string game_id = jv.at("game_id").as_string().c_str();
            std::string email = user_manager.getUserEmailFromSession(sessionId);
            bool success = game_manager.joinGame(game_id, email);

            json::object response;
            response["action"] = "join_result";
            response["success"] = success;
            if (success) {
                response["game_state"] = game_manager.getGameState(game_id);
            }
            else {
                response["error"] = "Unable to join game";
            }
            ws.write(net::buffer(json::serialize(response)));
        }
        else if (action == "make_move") {
            std::string sessionId = jv.at("session_id").as_string().c_str();
            if (!user_manager.validateSession(sessionId)) {
                send_error(ws, "Invalid session");
                return;
            }

            std::string game_id = jv.at("game_id").as_string().c_str();
            std::string email = user_manager.getUserEmailFromSession(sessionId);
            uint32_t x = static_cast<uint32_t>(jv.at("x").as_int64());
            uint32_t y = static_cast<uint32_t>(jv.at("y").as_int64());

            GameManager::GameError result = game_manager.makeMove(game_id, email, x, y);

            json::object response;
            response["action"] = "move_result";
            response["success"] = (result == GameManager::GameError::Success);
            if (result == GameManager::GameError::Success) {
                response["game_state"] = game_manager.getGameState(game_id);
                response["game_over"] = game_manager.isGameOver(game_id);
                response["current_player"] = game_manager.getCurrentPlayerEmail(game_id);
                if (game_manager.isGameOver(game_id)) {
                    User* user = user_manager.getUser(email);
                    response["new_rank"] = user->rank;
                    response["new_dan_rank"] = user->getDanRank();
                }
            }
            else {
                response["error"] = game_manager.getLastGameError(game_id);
            }
            ws.write(net::buffer(json::serialize(response)));
        }
        else {
            send_error(ws, "Unknown action");
        }
    }
    catch (const std::exception& e) {
        send_error(ws, std::string("Error processing message: ") + e.what());
    }
}

void handle_session(websocket::stream<tcp::socket>& ws)
{
    try
    {
        ws.accept();

        for (;;)
        {
            beast::flat_buffer buffer;
            ws.read(buffer);
            std::string message = beast::buffers_to_string(buffer.data());
            handle_message(ws, message);
        }
    }
    catch (beast::system_error const& se)
    {
        if (se.code() != websocket::error::closed)
            std::cerr << "Error: " << se.code().message() << std::endl;
    }
    catch (std::exception const& e)
    {
        std::cerr << "Error: " << e.what() << std::endl;
    }
}

int main(int argc, char* argv[])
{
    try
    {
        Config::load("config.json");

        auto const address = net::ip::make_address(Config::getString("server_address", "0.0.0.0"));
        auto const port = static_cast<unsigned short>(Config::getUInt32("server_port", 8080));

        net::io_context ioc{ 1 };

        tcp::acceptor acceptor{ ioc, {address, port} };

        std::cout << "Server listening on " << address << ":" << static_cast<unsigned int>(port) << std::endl;

        for (;;)
        {
            tcp::socket socket{ ioc };
            acceptor.accept(socket);
            std::thread{ std::bind(
                &handle_session,
                websocket::stream<tcp::socket>(std::move(socket))
            ) }.detach();
        }
    }
    catch (const std::exception& e)
    {
        std::cerr << "Error: " << e.what() << std::endl;
        return EXIT_FAILURE;
    }

    return EXIT_SUCCESS;
}