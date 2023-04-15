CREATE TABLE historyOfQueries
(
    ID            INTEGER primary key autoincrement,
    Author        TEXT    not null references Users (name),
    DBName        TEXT    not null references userDBs (dbName),
    Query         TEXT    not null,
    Status        INTEGER not null,
    ExecutionTime TEXT,
    Date          INTEGER not null
);
CREATE TABLE Ports
(
    Port TEXT unique
);
CREATE TABLE userDBs
(
    id      INTEGER primary key,
    connStr TEXT not null,
    owner   TEXT not null references Users (name),
    driver  TEXT not null,
    dbName  TEXT,
    docker  TEXT
);
CREATE TABLE Users
(
    Name     TEXT    not null unique,
    ID       INTEGER not null primary key autoincrement,
    Role     TEXT,
    Password TEXT,
    Nickname TEXT
);