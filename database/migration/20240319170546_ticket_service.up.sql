CREATE TABLE IF NOT EXISTS tickets (
    id BIGINT PRIMARY KEY,
    capacity BIGINT,
    region VARCHAR(255),
    Level VARCHAR(255),
    event_date TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ticket_details (
    id BIGINT PRIMARY KEY,
    ticket_id BIGINT,
    base_price FLOAT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (TicketID) REFERENCES ticket(ID)
);