# Remitly Task

This project is an implementation of a task by Remitly using Go and the Gin framework. It sets up a PostgreSQL database with initial data from csv file and creates a Gin server to handle requests. The project includes integration and unit tests to ensure the system works as expected.

## Project Structure

## How to Run

1. Clone the repository to your local machine:
    ```sh
    git clone https://github.com/Kortexik/RemitlyTask
    ```
2. Navigate to the project directory:
    ```sh
    cd RemitlyTask
    ```
3. Build and start the Docker containers:
    ```sh
    docker compose up --build
    ```
4. Once the containers are up and running, you can access the API at `http://localhost:8080/v1/swift-codes` using your browser or Postman.

## API Endpoints

- **Add New Swift Code**
    - **URL:** `POST /v1/swift-codes`
    - **Body:**
        ```json
        {
            "address": "TEST ADDRESS",
            "bankName": "TEST BANK",
            "countryISO2": "PL",
            "countryName": "POLAND",
            "isHeadquarter": false,
            "swiftCode": "TESTTESTTES"
        }
        ```
    - **Example response**
        ```json
            "message": "TESTTESTTES has been added to the database."
        ```

- **Get Swift Code Details**
    - **URL:** `GET /v1/swift-codes/:swift-code`
    - **Example response (`/v1/swift-codes/ALBPPLPWXXX`)**
      ```json
        {
            "address": "LOPUSZANSKA BUSINESS PARK LOPUSZANSKA 38 D WARSZAWA, MAZOWIECKIE, 02-232",
            "bankName": "ALIOR BANK SPOLKA AKCYJNA",
            "countryISO2": "PL",
            "countryName": "POLAND",
            "isHeadquarter": true,
            "swiftCode": "ALBPPLPWXXX",
            "branches": [
                {
                    "address": "LOPUSZANSKA BUSINESS PARK LOPUSZANSKA 38 D WARSZAWA, MAZOWIECKIE, 02-232",
                    "bankName": "ALIOR BANK SPOLKA AKCYJNA",
                    "countryISO2": "PL",
                    "isHeadquarter": false,
                    "swiftCode": "ALBPPLPWCUS"
                }
            ]
        }
      ```

- **Get Swift Codes by Country**
    - **URL:** `GET /v1/swift-codes/country/:ISO2`

- **Delete Swift Code**
    - **URL:** `DELETE /v1/swift-codes/:swift-code`
