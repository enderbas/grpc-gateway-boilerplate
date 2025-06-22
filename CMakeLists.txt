# --- Final CMakeLists.txt (v3) ---
cmake_minimum_required(VERSION 3.16)
project(grpc_gateway_example CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# --- Find Compiler Tools ---
# Find the tools using the correct executable names for your system
find_program(PROTOC_EXECUTABLE protoc)
find_program(GRPC_CPP_PLUGIN grpc_cpp_plugin) # Corrected plugin name

if(NOT PROTOC_EXECUTABLE)
    message(FATAL_ERROR "protobuf compiler 'protoc' not found. Please run 'sudo apt install protobuf-compiler'.")
endif()
if(NOT GRPC_CPP_PLUGIN)
    message(FATAL_ERROR "gRPC C++ plugin 'grpc_cpp_plugin' not found. Please check your installation (e.g., 'sudo apt install grpc-proto').")
endif()


# --- Proto Code Generation (Handles subdirectories automatically) ---
set(PROTO_SRC_DIR "${CMAKE_CURRENT_SOURCE_DIR}/proto")
set(PROTO_GEN_DIR "${CMAKE_CURRENT_BINARY_DIR}/generated-sources")
file(MAKE_DIRECTORY ${PROTO_GEN_DIR})

file(GLOB_RECURSE ALL_PROTO_FILES "${PROTO_SRC_DIR}/*.proto")
set(PROTO_GENERATED_SRCS "")
set(PROTO_GENERATED_HDRS "")

foreach(PROTO_FILE ${ALL_PROTO_FILES})
    get_filename_component(PROTO_FILENAME ${PROTO_FILE} NAME_WE)
    file(RELATIVE_PATH PROTO_REL_DIR ${PROTO_SRC_DIR} ${PROTO_FILE})
    get_filename_component(PROTO_REL_DIR ${PROTO_REL_DIR} DIRECTORY)
    set(GENERATED_H "${PROTO_GEN_DIR}/${PROTO_REL_DIR}/${PROTO_FILENAME}.pb.h")
    set(GENERATED_SRC "${PROTO_GEN_DIR}/${PROTO_REL_DIR}/${PROTO_FILENAME}.pb.cc")
    set(GENERATED_GRPC_H "${PROTO_GEN_DIR}/${PROTO_REL_DIR}/${PROTO_FILENAME}.grpc.pb.h")
    set(GENERATED_GRPC_SRC "${PROTO_GEN_DIR}/${PROTO_REL_DIR}/${PROTO_FILENAME}.grpc.pb.cc")
    list(APPEND PROTO_GENERATED_SRCS ${GENERATED_SRC} ${GENERATED_GRPC_SRC})
    list(APPEND PROTO_GENERATED_HDRS ${GENERATED_H} ${GENERATED_GRPC_H})
endforeach()

add_custom_command(
    OUTPUT ${PROTO_GENERATED_SRCS} ${PROTO_GENERATED_HDRS}
    COMMAND ${PROTOC_EXECUTABLE}
        --proto_path=${PROTO_SRC_DIR}
        --cpp_out=${PROTO_GEN_DIR}
        --grpc_out=${PROTO_GEN_DIR}
        --plugin=protoc-gen-grpc=${GRPC_CPP_PLUGIN} # Note: The flag name is still protoc-gen-grpc
        ${ALL_PROTO_FILES}
    DEPENDS ${ALL_PROTO_FILES}
    COMMENT "Generating Protobuf and gRPC C++ sources from all .proto files"
)

add_custom_target(GenerateProtoSources DEPENDS ${PROTO_GENERATED_SRCS} ${PROTO_GENERATED_HDRS})


# --- Server Executable ---
add_executable(greeter_server server/greeter_server.cc ${PROTO_GENERATED_SRCS})
add_dependencies(greeter_server GenerateProtoSources)

# Add the directory for our generated headers to the include path
target_include_directories(greeter_server PRIVATE ${PROTO_GEN_DIR})

# --- Linking (Using the direct method from your working example) ---
target_link_libraries(greeter_server PRIVATE
    grpc++
    grpc++_reflection
    protobuf
    gpr
    absl_synchronization
)