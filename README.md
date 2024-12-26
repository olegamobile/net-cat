
# TCPChat

This project is a recreation of the `NetCat` (`nc`) command-line utility in a Server-Client architecture. It enables TCP-based group chat functionality, supporting multiple clients connected to a single server.

## Features

- **TCP Connection**: Establishes a 1-to-many relationship between a server and multiple clients.
- **Client Identification**: Clients must provide a unique name upon connection.
- **Connection Control**: Manages the number of active client connections.
- **Message Transmission**: Clients can send messages to the group chat.
- **Message Format**: Each message includes:
  - Timestamp of when it was sent.
  - Name of the client who sent the message.
  - Message content.
  - Example: `[2024-12-26 15:48:41][Oleg]: Hello, Grit:Lab!`
- **Chat History**: New clients joining the chat receive all previously sent messages.
- **Client Notifications**:
  - When a client joins, all other clients are notified.
  - When a client leaves, all other clients are informed.
- **No Empty Messages**: Empty messages are not broadcast to the group.
- **Persistent Connections**: Clients leaving the chat do not cause disconnections for remaining clients.
  **Implemented commands**: Clients can leave the chat using command `/exit`.
- **Port**: The default port is `8989` if none is specified. If no port is provided, the program responds with the usage message:
  ```
  [USAGE]: ./TCPChat $port
  ```

## How to Run

1. **Build**:
   ```bash
   ./builder.sh
   ```

2. **Run**:
   Start the server using the executable:
   ```bash
   ./TCPChat 8989
   ```
   Replace `8989` with your desired port number if necessary.

3. **Connect a Client**:
   ```bash
   nc localhost 8989
   ```
   The client will be prompted to enter a name to join the chat.

## Dependencies

- Go version 1.16 or later.


## Notes

- Ensure the port number is valid and not in use by other applications.
- For more information on `NetCat`, refer to its manual using:
  ```bash
  man nc
  ```

Enjoy chatting!


## Authors

- [Anastasia Suhareva](https://01.gritlab.ax/git/asuharev) & [Oleg Balandin](https://01.gritlab.ax/git/obalandi) 
