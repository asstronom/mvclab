CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(30) UNIQUE NOT NULL,
    passwrd TEXT NOT NULL, 
    firstname TEXT,
    lastname TEXT,
    email TEXT,
    status TEXT,
    description TEXT
);