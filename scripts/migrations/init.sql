CREATE TABLE
  users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    balance VARCHAR(255) NOT NULL,
    salt VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
  );

CREATE TABLE
  api_credentials (
    id SERIAL PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL UNIQUE,
    client_secret VARCHAR(255) NOT NULL,
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
  );

CREATE TABLE
  recipient_accounts (
    id SERIAL PRIMARY KEY,
    account_number VARCHAR(255) NOT NULL,
    account_name VARCHAR(255) NOT NULL,
    verification_status VARCHAR(255) NOT NULL,
    last_verified_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
  );

CREATE TABLE
  transfers (
    id SERIAL PRIMARY KEY,
    recipient_account_id INTEGER REFERENCES recipient_accounts(id),
    sender_account_id INTEGER REFERENCES users(id),
    amount DECIMAL(12, 2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITHOUT TIME ZONE
  );

CREATE TABLE
  transfer_callbacks (
    id SERIAL PRIMARY KEY,
    transfer_id INTEGER REFERENCES transfers(id),
    status VARCHAR(50) NOT NULL,
    received_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
  );