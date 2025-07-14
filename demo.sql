CREATE DATABASE car_dealership;
CREATE USER myuser WITH PASSWORD 'mypassword';
GRANT ALL PRIVILEGES ON DATABASE car_dealership TO myuser;

\c car_dealership

CREATE TABLE brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE models (
    id SERIAL PRIMARY KEY,
    brand_id INTEGER REFERENCES brands(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    UNIQUE(brand_id, name)
);

CREATE TABLE cars (
    id SERIAL PRIMARY KEY,
    model_id INTEGER REFERENCES models(id) ON DELETE CASCADE,
    year INTEGER NOT NULL CHECK (year > 1900 AND year <= EXTRACT(YEAR FROM CURRENT_DATE) + 1),
    price INTEGER NOT NULL CHECK (price > 0)
);

INSERT INTO brands (name) VALUES
('Toyota'),
('Honda'),
('BMW'),
('Mercedes-Benz'),
('Audi'),
('Volkswagen'),
('Ford'),
('Chevrolet'),
('Nissan'),
('Hyundai'),
('Kia'),
('Lexus'),
('Mazda'),
('Subaru'),
('Volvo');

INSERT INTO models (brand_id, name) VALUES
(1, 'Corolla'),
(1, 'Camry'),
(1, 'RAV4'),
(1, 'Land Cruiser'),
(1, 'Prius'),
(2, 'Civic'),
(2, 'Accord'),
(2, 'CR-V'),
(2, 'Pilot'),
(2, 'Fit'),
(3, '3 Series'),
(3, '5 Series'),
(3, 'X5'),
(3, 'X3'),
(3, '7 Series'),
(4, 'C-Class'),
(4, 'E-Class'),
(4, 'S-Class'),
(4, 'GLC'),
(4, 'GLE'),
(5, 'A4'),
(5, 'A6'),
(5, 'Q5'),
(5, 'Q7'),
(5, 'TT');

INSERT INTO cars (model_id, year, price) VALUES
(1, 2018, 1200),  
(1, 2020, 1500),
(2, 2019, 1800),  
(2, 2021, 2200),
(3, 2017, 1600),  
(3, 2020, 2000),
(4, 2015, 3500),  
(4, 2022, 6000),
(5, 2019, 1700),
(6, 2018, 1100),  
(6, 2020, 1400),
(7, 2019, 1700),  
(7, 2021, 2100),
(8, 2017, 1500),  
(8, 2020, 1900),
(9, 2016, 2000),  
(10, 2019, 900),
(11, 2017, 2500), 
(11, 2020, 3200),
(12, 2018, 3500), 
(12, 2021, 4500),
(13, 2016, 4000), 
(13, 2022, 6500),
(14, 2019, 3000), 
(15, 2020, 5500),
(16, 2018, 2800), 
(16, 2021, 3800),
(17, 2017, 3200), 
(17, 2020, 4200),
(18, 2016, 4500), 
(18, 2022, 7500),
(19, 2019, 3500),  
(20, 2020, 4800),
(21, 2018, 2300),  
(21, 2020, 2900),
(22, 2017, 3000),  
(22, 2021, 4000),
(23, 2019, 3200),  
(24, 2018, 4500),  
(25, 2020, 3800);  
