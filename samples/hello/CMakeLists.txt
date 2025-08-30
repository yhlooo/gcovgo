cmake_minimum_required(VERSION 3.10)

set(CMAKE_C_STANDARD 11)
set(CMAKE_DEPENDS_USE_COMPILER true)
set(EXECUTABLE_OUTPUT_PATH ${PROJECT_SOURCE_DIR}/bin)
set(LIBRARY_OUTPUT_PATH ${PROJECT_SOURCE_DIR}/bin)

# 用于调试
set(CMAKE_VERBOSE_MAKEFILE ON)
set(CMAKE_BUILD_TYPE Debug)

add_library(lib1 STATIC
    src/path/to/lib1/lib1.h
    src/path/to/lib1/lib1_1.c
    src/path/to/lib1/lib1_2.c
)
target_compile_options(lib1 PRIVATE -O0 --coverage)
target_link_options(lib1 PRIVATE --coverage)

add_library(branches SHARED
    src/branches/branches.h
    src/branches/branches.c
)
target_compile_options(branches PRIVATE -O0 --coverage)
target_link_options(branches PRIVATE --coverage)

add_executable(hello
    src/main.c
    src/path/to/lib1/lib1.h
    src/branches/branches.h
)
target_link_libraries(hello PRIVATE lib1 branches)
target_compile_options(hello PRIVATE -O0 --coverage)
target_link_options(hello PRIVATE --coverage)
