# go-technical-test-jojonomic
backend developer technical test at jojonomic

# Technical Test Description:
Test Overview:

RestAPI using golang with kafka and docker

# Usage
Inside docker this project will running on default port \
    input-harga = localhost:8080 \
    topup = localhost:8082 \
    cek-harga = localhost:8084 \
    cek-mutasi = localhost:8085 \
    cek-saldo = localhost:8086 \
    buyback = localhost:8087

# Endpoint
API Documentation are attached in 
**postman folder** and you could import to other device to see the detailed body and response message.

# Installation

Follow these steps to install and set up your project locally. Make sure you meet the prerequisites before proceeding.

### Prerequisites

Before you begin, ensure you have the following prerequisites installed on your system:

- [Go](https://golang.org/dl/): This project is built with Go, so you need to have Go installed on your system.
- [docker]: docker is required for running inside a container with mysql.

### Step-by-Step Installation

1. **Setup the Repository:**

   For running with docker i've a command with Makefile command as a :

   ```bash
   cd misc
   docker-compose up --build


