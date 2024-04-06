# Money Transfer Service

This project is a simulation of a money transfer service, demonstrating the basic flow of Transferring Money from a system to another system.

## Architecture Overview

The application is structured by following the requested guidelines, which ensures that the code is organized and hopefully... maintainable & scalable. The flow of control starts from the controller, moves through use cases, interacts with the database via repositories, and finally presents the response through the presenter. This flow is bidirectional:

`Controller => UseCase => Repositories => UseCase => Presenter`

By adhering to this architecture, the application maintains a clear separation of concerns, making it easier to manage and evolve over time.

## Getting Started

To simplify the project setup and ensure a consistent environment across different machines, the application is containerized using Docker. You can start the project with the following command:

`docker-compose up`

This command builds the necessary images and starts the containers as defined in the `docker-compose.yml` file. Ensure Docker and Docker Compose are installed on your system before running this command.

OR

You could clone this repo, then

run `go mod tidy`, then

run `go run main.go`

## Implementation Details

### Task 1: Account Validation

The first endpoint implemented in this project is for account validation. It demonstrates the application's concurrency capabilities by running I/O tasks in parallel:

- One goroutine handles the database interaction, inserting the account information into the `recipient_accounts` table.
- Another goroutine makes an API call to a mock banking service to validate the account details.

The application waits for both operations to complete using an `errChan` and receives the `recipientAccountId` through a channel. Upon successful completion of both tasks, it updates the `recipient_accounts` table, setting the `status` to "Completed" and updating `last_verified_at` to reflect the time of verification.

### Task 2: Transfer/Disbursement

Similar to Task 1, the money transfer functionality is implemented with concurrency in mind:

- The application sends the money transfer request while concurrently listening for a callback from the provided `callbackUrl`.
- Once the transfer is initiated, a separate process listens for the callback to update the transaction status accordingly.

This approach allows the service to handle requests efficiently, ensuring that the system remains responsive even while waiting for external processes to complete.

## Task 3: Callback URL

- Implementation of the Transfer/Disbursement Callback endpoint is done by setting up a listener for status updates from the bank after initiating a money transfer. The process will update the transaction status in the database based on the callback received.
- This process is also used in accordance with task 2, where a go routine dispatched to call this Callback URL.

---
