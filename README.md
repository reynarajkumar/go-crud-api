# go-crud-api

This is a simple CRUD (Create, Read, Update, Delete) API written in Go for managing user profiles. It provides endpoints for adding, retrieving, updating, and deleting user profiles.

## Features

- **Create Profile**: Add a new user profile with department, designation, and employee details.
- **Retrieve Profiles**: Retrieve all user profiles or a specific profile by ID.
- **Update Profile**: Update an existing user profile by ID.
- **Delete Profile**: Delete a user profile by ID.

## Dependencies

- [Gorilla Mux](https://github.com/gorilla/mux): A powerful HTTP router and URL matcher for building Go web servers.

## Usage

1. Install Go and set up your Go workspace.
2. Clone this repository into your Go workspace.
3. Install dependencies by running:
   ```bash
   go mod tidy
   ```
4. Run the server:
   ```bash
   go run main.go
   ```
5. The server will start running on `http://localhost:8080`.

## Endpoints

- `POST /profiles`: Add a new profile.
- `GET /profiles`: Retrieve all profiles.
- `GET /profiles/{id}`: Retrieve a specific profile by ID.
- `PUT /profiles/{id}`: Update a specific profile by ID.
- `DELETE /profiles/{id}`: Delete a specific profile by ID.

## Concurrency

The program demonstrates parallelized and unparallelized API calls for adding profiles. It measures and compares the time taken for both methods.
