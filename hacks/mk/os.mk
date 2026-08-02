ifeq ($(OS),Windows_NT)
	OS := Windows
endif

ifeq ($(DETECTED_OS),Linux)
	OS := Linux
endif

ifeq ($(DETECTED_OS),Darwin)
	OS := MacOS
endif