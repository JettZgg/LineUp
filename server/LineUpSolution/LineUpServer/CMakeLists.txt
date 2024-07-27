cmake_minimum_required(VERSION 3.10)

project(LineUpServer VERSION 1.0)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED True)

# Set Boost root directory and library directory
set(BOOST_ROOT "C:/boost_1_85_0")
set(BOOST_LIBRARYDIR "${BOOST_ROOT}/stage/lib")

set(Boost_USE_STATIC_LIBS ON)
set(Boost_USE_MULTITHREADED ON)
set(Boost_USE_STATIC_RUNTIME OFF)

# Find Boost
find_package(Boost 1.85.0 REQUIRED COMPONENTS system json)

if(NOT Boost_FOUND)
    message(FATAL_ERROR "Boost not found. Please check your Boost installation.")
endif()

# Set OpenSSL root directory
set(OPENSSL_ROOT_DIR "C:/OpenSSL-Win64")

# Find OpenSSL
find_package(OpenSSL REQUIRED)

if(NOT OpenSSL_FOUND)
    message(FATAL_ERROR "OpenSSL not found. Please check your OpenSSL installation.")
endif()

# Add all source files
add_executable(LineUpServer 
    LineUpServer.cpp
    Game.cpp
    GameManager.cpp
    UserManager.cpp
    Config.cpp
    InputValidation.cpp
    Session.cpp
)

# Include directories
target_include_directories(LineUpServer PRIVATE 
    ${CMAKE_CURRENT_SOURCE_DIR} 
    ${Boost_INCLUDE_DIRS}
    ${OPENSSL_INCLUDE_DIR}
)

# Link libraries
target_link_libraries(LineUpServer PRIVATE 
    ${Boost_LIBRARIES}
    OpenSSL::SSL
    OpenSSL::Crypto
)

# Set the output directory
set_target_properties(LineUpServer
    PROPERTIES
    RUNTIME_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}/bin"
)

# Print Boost information for debugging
message(STATUS "Boost_INCLUDE_DIRS: ${Boost_INCLUDE_DIRS}")
message(STATUS "Boost_LIBRARIES: ${Boost_LIBRARIES}")

# Print OpenSSL information for debugging
message(STATUS "OPENSSL_INCLUDE_DIR: ${OPENSSL_INCLUDE_DIR}")
message(STATUS "OPENSSL_LIBRARIES: ${OPENSSL_LIBRARIES}")