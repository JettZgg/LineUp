cmake_minimum_required(VERSION 3.10)

project(LineUpServer VERSION 1.0)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED True)

# Find Boost
find_package(Boost REQUIRED COMPONENTS system json)

# Add the executable
add_executable(LineUpServer 
    LineUpServer.cpp
    Game.cpp
    GameManager.cpp
    UserManager.cpp
)

# Include directories
target_include_directories(LineUpServer PRIVATE ${CMAKE_CURRENT_SOURCE_DIR})

# Link libraries
target_link_libraries(LineUpServer PRIVATE Boost::system Boost::json)

# Set the output directory
set_target_properties(LineUpServer
    PROPERTIES
    RUNTIME_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}/bin"
)