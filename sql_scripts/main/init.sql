CREATE TABLE client_services
(
    id SERIAL PRIMARY KEY,
    client_id INT NOT NULL,
    service_id INT NOT NULL,
    bill_id INT NOT NULL,
    price INT NOT NULL,
    description_uuid UUID UNIQUE NOT NULL
);