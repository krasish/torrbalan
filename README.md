# Torrbalan
*Course project for 'Introduction to Go 2020/2021' @ FMI*

### Overview

*Torrbalan* is a simple peer-to-peer file exchange application. *Torrbalan* uses TCP connections for communication between its components. *Torrbalan* consists of 2 parts:
- Torrbalan server - responsible for keeping information about files and providing it to clients
- Torrbalan client - a simple CLI applicaiton which communicates to *Torrbalan* servers and downloads/uploads applications to other *Torrbalan* clients

### Usage
You can start using *Torrbalan* by simply cloning this repo. 

In order to use *Torrbalan* you will first need to start a *Torrbalan* server. To do that, go to the directory in which you cloned *Torrbalan*, replace **\<port\>** in the following command with the port you wish your server to listen and run:

`go run ./server/cmd/main.go <port>`

Afterwards, to start a *Torrbalan* client run the following command replacing **\<server-address\>** with the address of the server you previously run:

`go run ./client/cmd/main.go <server-address>`

Thus, an example start of a server and a client would be:

`go run ./server/cmd/main.go 8080`

`go run ./client/cmd/main.go localhost:8080`

