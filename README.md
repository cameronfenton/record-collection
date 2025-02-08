# record-collection backend
An application for keeping track of my record and media collection.

## Building and Starting the Application

### Requirements
- Go (version 1.19 or higher)
- MySQL version 5.7 or higher
- A running MySQL server instance
- Git

### Installation
1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/record-collection-backend.git
    ```
2. Navigate to the project directory:
    ```sh
    cd record-collection-backend
    ```

### Building the Application

#### On Linux, Windows, and macOS
To build the application, run:
```sh
go build
```

### Starting the Application
To start the application, run:
```sh
./record-collection-backend
```

The application should now be running on `http://localhost:8080`(the specified port in the config.json).
