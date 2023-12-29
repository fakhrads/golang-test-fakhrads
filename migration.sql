-- migration.sql

CREATE TABLE IF NOT EXISTS users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    address VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    photos JSON NOT NULL,
    creditcard_type VARCHAR(50) NOT NULL,
    creditcard_number VARCHAR(16) NOT NULL,
    creditcard_name VARCHAR(255) NOT NULL,
    creditcard_expired VARCHAR(10) NOT NULL,
    creditcard_cvv VARCHAR(4) NOT NULL
);