PRAGMA foreign_keys = ON;

CREATE TABLE client (
    clientID TEXT UNIQUE NOT NULL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    pubKey TEXT NOT NULL,
    address TEXT NOT NULL,
    jobsAvailable INTEGER DEFAULT 0
);

CREATE TABLE session (
    clientID TEXT NOT NULL PRIMARY KEY,
    workerID TEXT NOT NULL PRIMARY KEY,
    sessionID TEXT NOT NULL,
    FOREIGN KEY (clientID)
        REFERENCES client (clientID),
    FOREIGN KEY (workerID)
        REFERENCES client (clientID)
);

CREATE INDEX emailIndex ON client(email);
CREATE INDEX clientIndex ON client(clientID);
CREATE INDEX sessionIndex ON session(sessionID);