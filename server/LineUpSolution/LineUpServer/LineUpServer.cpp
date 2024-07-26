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
            std::string email = jv.at("email").as_string().c_str();
            std::string username = jv.at("username").as_string().c_str();
            std::string password = jv.at("password").as_string().c_str();
            bool success = user_manager.registerUser(email, username, password);

            json::object response;
            response["action"] = "register_result";
            response["success"] = success;
            ws.write(net::buffer(json::serialize(response)));
        }
        else if (action == "login") {
            std::string email = jv.at("email").as_string().c_str();
            std::string password = jv.at("password").as_string().c_str();
            bool success = user_manager.loginUser(email, password);

            json::object response;
            response["action"] = "login_result";
            response["success"] = success;
            if (success) {
                User* user = user_manager.getUser(email);
                response["username"] = user->username;
                response["rank"] = user->rank;
                response["dan_rank"] = user->getDanRank();
            }
            ws.write(net::buffer(json::serialize(response)));
        }
        else if (action == "create_game") {
            std::string email = jv.at("email").as_string().c_str();
            uint32_t width = jv.at("width").as_int64();
            uint32_t height = jv.at("height").as_int64();
            uint32_t win_condition = jv.at("win_condition").as_int64();

            if (width < 3 || width > 99 || height < 3 || height > 99 || win_condition < 3 || win_condition > 19) {
                send_error(ws, "Invalid game parameters");
                return;
            }

            std::string game_id = game_manager.createGame(email, width, height, win_condition);

            json::object response;
            response["action"] = "game_created";
            response["game_id"] = game_id;
            ws.write(net::buffer(json::serialize(response)));
        }
        else if (action == "join_game") {
            std::string game_id = jv.at("game_id").as_string().c_str();
            std::string email = jv.at("email").as_string().c_str();
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
            std::string game_id = jv.at("game_id").as_string().c_str();
            std::string email = jv.at("email").as_string().c_str();
            uint32_t x = jv.at("x").as_int64();
            uint32_t y = jv.at("y").as_int64();

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
        // Check command line arguments.
        if (argc != 3)
        {
            std::cerr << "Usage: LineUpServer <address> <port>\n";
            std::cerr << "Example:\n";
            std::cerr << "    LineUpServer 0.0.0.0 8080\n";
            return EXIT_FAILURE;
        }

        auto const address = net::ip::make_address(argv[1]);
        auto const port = static_cast<unsigned short>(std::atoi(argv[2]));

        // The io_context is required for all I/O
        net::io_context ioc{ 1 };

        // The acceptor receives incoming connections
        tcp::acceptor acceptor{ ioc, {address, port} };

        std::cout << "Server listening on " << address << ":" << port << std::endl;

        for (;;)
        {
            // This will receive the new connection
            tcp::socket socket{ ioc };

            // Block until we get a connection
            acceptor.accept(socket);

            // Launch the session, transferring ownership of the socket
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