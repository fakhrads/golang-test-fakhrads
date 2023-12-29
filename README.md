## Getting Started
Database migration will be automatically applied when you run ```go run main.go```
### Prerequisites

- [Golang](https://golang.org/) installed on your machine.
- MySQL installed and a database set up.

### Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/your_username/your_project.git
2. **Navigate to the project directory**
3. **Create a .env file in the root directory and configure your MySQL database connection:**
    ```.env
    DB_USER=root
    DB_PASSWORD=password
    DB_NAME=your_database_name
    DB_HOST=localhost
    DB_PORT=3306
    API_KEY="HiJhvL\$T27@1u^%u86g"
    ```
3. **Install the necessary Go dependencies:**
    ```
    go mod tidy
    ```
4. **Run the application**
    ```
    go run main.go
    ```
    The application should now be running on http://localhost:3000.

## API Documentation
The API documentation is available at :

https://www.postman.com/fakhrads/workspace/golang-test/collection/8642679-be6ea13e-18b1-4c21-b0b7-1d4cdade8d7a?action=share&creator=8642679