# Joona

This repository is a mock backend service for a recipe browser application where people can find delicious recipes based on the ingredients and tools that they have, food preference, and their cooking mastery level. 

## Introduction

This backend application is written in Go (Golang) which the client can connect via RESTful API's with MongoDB as its main database, Elasticsearch as its secondary database for search, and Redis for caching. This backend service is also equipped with a CI/CD tool using Github Actions to test and deploy to an Azure cloud server. 
Although the tech stacks seem way too overkill for a mock service, they are chosen as a proof of concept of how such an application would be designed in a large production environment.


## Prerequisites

- Golang 1.21 or later
- Docker Coompose

## Setup Instructions

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/eifzed/joona.git
   ```

2. **Navigate to the Project Directory:**
   ```bash
   cd joona
   ```

2. **Setup DB:**
   ```bash
   docker-compose up
   ```

3. **Setup Config File**
    ```bash
   mv ./files/etc/joona-config/joona-config.development.yaml.example ./files/etc/joona-config/joona-config.development.yaml
   mv ./files/etc/joona-secret/joona-secret.development.json.example ./files/etc/joona-secret/joona-secret.development.json
   ```
3. **Run Application:**
   ```bash
   docker-compose up --build
   ```

4. **Access the Application:**
   Once the containers are up and running, you can access the GoLang application at [http://localhost:8080](http://localhost:8080).

## Directory Structure

- `cmd/`: Contains the GoLang application source code.
- `files/`: Contains the configuration files.
- `internal/`: Contains internal source code.
- `lib/`: Contains common library.
- `model/`: Contains data model.
- `Dockerfile`: Dockerfile for building the GoLang application container.
- `docker-compose.yml`: Docker Compose configuration file for managing containers.
- `README.md`: This README file.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

This project is licensed under the [MIT License](LICENSE).

