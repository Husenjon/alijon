CREATE TABLE IF NOT EXISTS branch(
    id serial PRIMARY KEY,
    inn VARCHAR (255),
    code VARCHAR (255) UNIQUE NOT null,
    name VARCHAR (255)
);

CREATE TABLE IF NOT EXISTS additional_bank(
    id serial PRIMARY KEY,
    code VARCHAR (255) UNIQUE NOT null,
    name VARCHAR (255) NOT null
);

CREATE TABLE IF NOT EXISTS bank(
    id serial PRIMARY KEY,
    code VARCHAR (255) UNIQUE NOT null,
    name VARCHAR (255) NOT null
);

CREATE TABLE IF NOT EXISTS contract(
    id serial PRIMARY KEY,
    number VARCHAR (255) UNIQUE NOT null,
    registration_card VARCHAR (255),
    date VARCHAR (255)
);

CREATE TABLE IF NOT EXISTS client(
    id serial PRIMARY KEY,
    inn VARCHAR (255) UNIQUE NOT null,
    name VARCHAR (255)
);

CREATE TABLE IF NOT EXISTS bags(
    id serial PRIMARY KEY,
    contract_id integer,
    name VARCHAR (255),
    status VARCHAR (255),
    FOREIGN KEY (contract_id) REFERENCES contract (id)
);

CREATE TABLE IF NOT EXISTS contracts(
    id serial PRIMARY KEY,
    contract_number Varchar(255) unique not null,
    additional_bank_id integer,
    debt float8,
    route VARCHAR(255),
    rate VARCHAR(255),
    bag_numbers VARCHAR(255),
    bank_id integer,
    contract_id integer,
    state VARCHAR(255),
    client_id integer,
    branch_id integer,
    FOREIGN KEY (additional_bank_id) REFERENCES additional_bank (id),
    FOREIGN KEY (bank_id) REFERENCES bank (id),
    FOREIGN KEY (contract_id) REFERENCES contract (id),
    FOREIGN KEY (client_id) REFERENCES client (id),
    FOREIGN KEY (branch_id) REFERENCES branch (id)
);

CREATE TABLE IF NOT EXISTS revenue(
    id serial PRIMARY KEY,
    client_id integer,
    sum float8,
    contract_id integer,
    date VARCHAR (255),
    FOREIGN KEY (client_id) REFERENCES client (id),
    FOREIGN KEY (contract_id) REFERENCES contract (id)
);

CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    ism varchar(255) NOT NULL,
    familya varchar(255) NOT NULL,
    otasini_ismi varchar(255) NOT NULL,
    phone varchar(127),
    username varchar(255) UNIQUE NOT NULL,
    password varchar(512) NOT NULL,
    branch_id INTEGER NOT NULL,
    token varchar(512),
    image VARCHAR(1023),
    is_active BOOLEAN NOT NULL,
    created_time timestamp default CURRENT_TIMESTAMP,
    updated_time timestamp
);