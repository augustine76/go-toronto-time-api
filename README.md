# Go Toronto Time API

This project is a simple Go application that retrieves the current time in Toronto, logs it to a MySQL database, and provides endpoints to access this data. It is fully containerized using Docker and Docker Compose.

## Features

- Retrieves the current time in the "America/Toronto" time zone.
- Logs the timestamps to a MySQL database.
- Provides RESTful API endpoints to get the current time and retrieve logs.
- Fully containerized with Docker and orchestrated using Docker Compose.

## Prerequisites

- **Docker** installed on your system.
- **Docker Compose** installed on your system.

## Project Structure

```
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
└── db
    └── init.sql
```

- `main.go`: The main Go application source code.
- `Dockerfile`: Dockerfile to build the Go application image.
- `docker-compose.yml`: Docker Compose configuration file.
- `go.mod` and `go.sum`: Go modules files.
- `index.html`: Frontend interface.
- `db/init.sql`: SQL script to initialize the database and create the necessary table.

## Setup Instructions

### Clone the Repository

```bash
git clone https://github.com/yourusername/go-toronto-time-api.git
cd go-toronto-time-api
```

### Environment Variables

The application uses environment variables defined in the `docker-compose.yml` file. Ensure they are set correctly:

```yaml
services:
  app:
    environment:
      - DB_USER=root
      - DB_PASSWORD=yourpassword
      - DB_HOST=db
      - DB_PORT=3306
      - DB_NAME=time_logger

  db:
    environment:
      - MYSQL_ROOT_PASSWORD=yourpassword
      - MYSQL_DATABASE=time_logger
```

- Replace `yourpassword` with a secure password of your choice.

### Initialize the Database

The database is initialized automatically using the `db/init.sql` script, which creates the `time_log` table.

**Contents of `db/init.sql`:**

```sql
USE time_logger;

CREATE TABLE IF NOT EXISTS time_log (
  id INT AUTO_INCREMENT PRIMARY KEY,
  timestamp DATETIME NOT NULL
);
```

## Running the Application

Build and start the application using Docker Compose:

```bash
docker-compose up --build
```

This command will:

- Build the Docker image for the Go application.
- Start the Go application container.
- Start the MySQL database container.
- Initialize the database and create the `time_log` table.

## API Endpoints

### GET `/current-time`

- **Description:** Retrieves the current time in Toronto and logs it to the database.
- **Example Request:**

  ```bash
  curl http://localhost:8080/current-time
  ```

- **Example Response:**

  ```json
  {
    "current_time": "2024-11-28T15:58:00-05:00"
  }
  ```

### GET `/logs`

- **Description:** Retrieves all logged timestamps from the database.
- **Example Request:**

  ```bash
  curl http://localhost:8080/logs
  ```

- **Example Response:**

  ```json
  [
    {
      "id": 1,
      "timestamp": "2024-11-28T15:58:00-05:00"
    },
    {
      "id": 2,
      "timestamp": "2024-11-28T16:00:00-05:00"
    }
  ]
  ```


## Stopping the Application

To stop the running containers, press `Ctrl+C` in the terminal where `docker-compose` is running, or run:

```bash
docker-compose down
```

## Troubleshooting

### Common Issues

#### Application Cannot Connect to Database

- **Symptom:** Application exits with an error related to database connection.
- **Solution:** Ensure the database container is running and accessible. The application includes retry logic to handle initial connection delays.

#### `time_log` Table Does Not Exist

- **Symptom:** Database error indicating the `time_log` table is missing.
- **Solution:** Ensure the `init.sql` script is correctly mapped in `docker-compose.yml` and that the database volume is initialized properly.

#### Port Conflicts

- **Symptom:** Error indicating that port `8080` or `3306` is already in use.
- **Solution:** Change the host port mappings in `docker-compose.yml` to use different ports.

