# Tragac na Inflacija

## Overview
This project implements a file upload and expense processing system with the following components:
- **Go API Server**: Accepts file upload requests, stores files in a shared data volume, saves file metadata into a PostgreSQL database, and publishes file metadata to a RabbitMQ queue for processing. It also consumes processed data from the RabbitMQ queue and populates the expenses table.
- **PostgreSQL Database**: Contains two tables - `files` for storing file metadata and `expenses` for storing processed expense data.
- **RabbitMQ**: Two queues are used: `pending_files` for the Go server to publish file metadata and `processed_files` for the Python application to send processed expense data back to the Go server.
- **Python Application**: Consumes messages from the `pending_files` queue, processes the associated files (PDF receipts), extracts expense information, and publishes the result to the `processed_files` queue.

The purpose of this project is to expose an API that allows users to upload receipt files, parse expense information from them, and store the data in a structured way for further analysis and price tracking. The system currently supports PDF receipt files. The Python application will process these PDF files to extract expense data, which should include:
- Expense amount
- Date of expense
- Category (e.g., food, travel, etc.)

## Components
### 1. **Go API Server**
The Go API server provides the following functionality:
- **POST /upload**: Accepts a file upload request, saves the file to a shared volume, and stores its metadata (id, name, path, type) in the `files` table.
- **Consumes processed_files**: Retrieves processed expense data from the `processed_files` queue and stores it in the `expenses` table for analysis.

### 2. **PostgreSQL Database**
The database consists of two main tables:
- `files`: Stores file metadata (file id, name, path, type).
- `expenses`: Stores the parsed expense data (expense id, amount, date, category).

### 3. **RabbitMQ**
RabbitMQ facilitates message passing between components:
- `pending_files` queue: The Go API server publishes file metadata here for processing by the Python app.
- `processed_files` queue: The Python app publishes parsed expense data here, which the Go server consumes to populate the expenses table.

### 4. **Python Application**
The Python application is responsible for:
- Consuming file metadata from the `pending_files` queue.
- Processing files (PDF receipts) based on the metadata.
- Extracting and parsing expense information from the receipts.
- Publishing processed data to the `processed_files` queue for the Go server to store in the `expenses` table.

## Setup and Installation

### Prerequisites
- Docker and Docker Compose
- Go 1.22 or higher
- Python 3.12 or higher
- PostgreSQL instance (can be set up using Docker)
- RabbitMQ instance (can be set up using Docker)

### Running the Project
1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/yourproject.git
   cd yourproject
2. **Configure environment variables**<br> 
You need to set up environment variables for each component in the project. Create two .env files for the Go API and Python application. Create additional docker-compose.env file and make sure variables are consistent across these files.

3. **Set up Docker containers**<br>
You can use Docker Compose to spin up all the necessary containers (Go API, Python app, PostgreSQL, RabbitMQ). Run the following command in the root of the repository:
    ```bash
    docker-compose up --build
4. **Access the API**<br>
Once the Docker containers are up and running, the Go API server will be available at `http://localhost:8080`. You can interact with the API as follows:<br>
`POST /upload`: Upload a file (PDF receipt) to store its metadata and initiate processing.
5. **Stop the containers**<br>
To stop the Docker containers, run:
    ```
    docker-compose down
    ```
