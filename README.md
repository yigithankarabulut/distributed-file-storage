# Distributed File Storage in Go

Build a decentralized, fully distributed, content-addressable file storage system using Go that can handle and stream very large files. This project covers system design, low-level programming, and network protocols, all while building a highly practical and scalable application.

## Table of Contents
- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
- [Example Usage](#example-usage)
- [Reference](#reference)
- [Contributing](#contributing)
- [License](#license)

## Introduction
This project aims to create a decentralized, fully distributed, content-addressable file storage system using Go. The system is designed to handle and stream very large files efficiently. Key features and concepts covered in this project include:

- **Decentralization**: The file storage system does not rely on a central server. Instead, it uses a peer-to-peer (P2P) network where each node can store and share files.
- **File Server**: Manages file storage, retrieval, and deletion. Each file server operates independently but can communicate with other file servers in the network.
- **File Storage and Retrieval**: The project covers efficient methods for storing, retrieving, and deleting files in a distributed environment.
- **TCP Transport Layer**: Custom TCP transport is used for peer-to-peer communication, including custom decoders and handshake functions.
- **Low-Level Programming**: The project involves low-level programming techniques in Go, providing a deep understanding of system-level operations and memory management.
- **Network Protocols**: The system uses custom network protocols for efficient data transfer and communication between nodes. TCP is the primary transport layer used.
- **Peer-to-Peer Communication**: Nodes in the network communicate directly with each other to exchange files and information. This involves implementing custom P2P communication protocols.
- **Content Addressable Storage (CAS)**: Uses CAS mechanisms to ensure that each piece of data is uniquely identified and stored based on its content.
- **Streaming and Caching**: Techniques for streaming large files and caching frequently accessed data to improve performance are implemented.
- **Encryption and Security**: Files are encrypted to ensure data security and privacy during storage and transmission.
- **Bootstrap Nodes**: Initial nodes used to join the P2P network and discover other nodes.


## Installation
To install and set up this project, follow these steps:

1. **Clone the repository:**
    ```sh
    git clone https://github.com/yigithankarabulut/distributed-file-storage.git
    cd distributed-file-storage
    ```

2. **Install dependencies:**
    ```sh
    go mod download
    ```

## Usage
We have some make commands, use the following commands:

```sh
make build
make run
make test
```

## Example Usage
The main program sets up multiple file servers on different ports, demonstrating how they interact and handle file operations. Here's a brief overview of what happens in the `main.go` file:

- Three file servers are created on different ports (`:3000`, `:4000`, and `:5000`).
- The servers are interconnected, with the third server (`:5000`) knowing about the first two servers (`:3000`, `:4000`).
- The servers start running in separate goroutines.
- The main program performs file storage and retrieval operations, demonstrating the system's functionality.

This project provides a comprehensive example of building a distributed file storage system from scratch, leveraging Go's capabilities for low-level programming and efficient network communication.

## Reference
Thanks to Anthony GG.
- [Anthony GG's GitHub Repository](https://github.com/anthdm/distributedfilesystemgo)
- [YouTube Tutorial](https://www.youtube.com/watch?v=IoY6bE--A54)


## Contributing
Contributions are welcome! Please follow these steps to contribute:

Fork the repository.\
Create a new branch (git checkout -b feature/YourFeature).\
Commit your changes (git commit -m 'Add some feature').\
Push to the branch (git push origin feature/YourFeature).\
Open a pull request.

## License
This project is licensed under the MIT License - see the LICENSE file for details.
