CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    age INT NOT NULL
);

INSERT INTO users (name, age) VALUES
('Alice', 30),
('Bob', 25);
