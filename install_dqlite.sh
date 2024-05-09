#!/bin/bash

# Update package lists
sudo apt update

# Upgrade installed packages
sudo apt upgrade -y

# Clone dqlite repository
git clone https://github.com/LakshK98/dqlite.git
cd dqlite
git checkout metric-branch

# Install required packages
sudo apt install -y pkg-config autoconf automake libtool make libuv1-dev libsqlite3-dev golang-go liblz4-dev

# Prepare dqlite for building
autoreconf -i

# Configure dqlite build with raft support
./configure --enable-build-raft

# Build dqlite
make

# Install dqlite
sudo make install

cd ..

# Clone go-dqlite repository
git clone https://github.com/SaumyaSachdev/go-dqlite.git

# Navigate to go-dqlite directory
cd go-dqlite
git checkout add-conc

# Further actions can be added here if needed
go install -tags libsqlite3 ./cmd/dqlite-demo
go install -tags libsqlite3 ./cmd/dqlite

sudo ldconfig
