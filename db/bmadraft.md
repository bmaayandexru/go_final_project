.open scheduler.db

CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "20000101", 
    title VARCHAR(32) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT ""
    repeat VARCHAR(128) NOT NULL DEFAULT "",
);
// индекс создавать не нужно лн создаётся автоматически
// CREATE INDEX idx_id ON id; 
CREATE INDEX idx_date ON scheduler (date); // нужен индекс по date
CREATE UNIQUE INDEX idx_title ON scheduler (title); // (bma) будет полезным и явно не лишим
