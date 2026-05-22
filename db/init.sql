CREATE TABLE IF NOT EXISTS users (
    id_user SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE users
ADD CONSTRAINT users_email_unique UNIQUE (email);

ALTER TABLE users
ADD CONSTRAINT users_username_unique UNIQUE (username);

CREATE TABLE IF NOT EXISTS approved_users (
    id_approved_users SERIAL PRIMARY KEY,
    id_user INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE approved_users
ADD COLUMN negated BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE approved_users
ADD COLUMN email_user VARCHAR(255) NOT NULL DEFAULT '';

CREATE TABLE IF NOT EXISTS account_to_pay (
id_account_to_pay SERIAL PRIMARY KEY,
id_user INT,
description varchar(100),
description_details varchar(255),
date_action date,
date_previous date,
value_pag DECIMAL(15,2),
value_add DECIMAL(15,2),
value_discount DECIMAL(15,2),
name_pag VARCHAR(100),
paid BOOLEAN);

CREATE TABLE IF NOT EXISTS account_to_pay_payments (
id_account_to_pay_payments SERIAL PRIMARY KEY,
id_payments INT,
id_account_to_pay INT
);

CREATE TABLE IF NOT EXISTS account_to_pay_documents (
id_account_to_pay_documents SERIAL PRIMARY KEY,
id_account_to_pay INT,
documents bytea
);

CREATE TABLE IF NOT EXISTS payments (
id_payments INT,
name varchar(100)
);

INSERT INTO payments (id_payments, name) VALUES
(1, 'Dinheiro'),
(2, 'Cartão de Crédito'),
(3, 'Cartão de Débito'),
(4, 'Boleto'),
(5, 'Pix'),
(6, 'Transferência Bancária'),
(7, 'Vale Alimentação'),
(8, 'Vale Refeição'),
(9, 'Cheque');