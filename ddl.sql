CREATE TABLE IF NOT EXISTS car_models (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    deleted BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    car_model_id INT REFERENCES car_models(id),
    department VARCHAR(100),
    usage_date VARCHAR(20),
    destination VARCHAR(100),
    deleted BOOLEAN DEFAULT false
);

-- INSERT INTO car_models (name) VALUES ('KIA CEED - CB 2376 BE'), ('KIA CEED - CB 2378 BE');

-- INSERT INTO bookings (car_model_id, department, usage_date, destination, deleted)
-- VALUES (1, 'IT Department', '2024-05-06', 'Blagoevgrad', false),
--        (2, 'Finance', '2024-05-10', 'Sofia', false);
