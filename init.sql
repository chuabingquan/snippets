CREATE DATABASE snippetsDB;

CREATE TABLE account (
    id uuid PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(25) UNIQUE NOT NULL,
    password_hash text NOT NULL,
    password_salt text NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL
);

CREATE TABLE snippet (
    id uuid PRIMARY KEY,
    account_id uuid NOT NULL REFERENCES account(id),
    filename VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    is_public BOOLEAN NOT NULL
);

INSERT INTO account VALUES
('6ab591ee-519a-487d-a2b5-27e308f81242', 'admin@snippets.com', 'admin', '$2a$08$ZhgZk6noN5QfLHzAq7IyaO9/scR6O52pdRtz0THbB3ncaMAxrnF7e', 'f1nd1ngn3m0', 'Admin', 'Test'), -- P@ssw0rd --
('1c99fc26-1a69-41d7-bd31-ef8156166917', 'charlotte.l@gmail.com', 'charlottelaw', '$2a$08$/C3lPRjqPukiTVtAeJX2mu.6mS.BSQBbbie7o7nglKHo9McgrL00a', 'farm1990M0O', 'Charlotte', 'Lawerence'), -- cherrykitty --
('9b7c9167-f139-44ea-910e-60211cb389f2', 'johnmendes88@gmail.com', 'johnmendes88', '$2a$08$5tb41iw4IqOGT2lxCovyD.ZfmiQ8jDz836D.BY95zKB3OybiDm0yG', 'scr3wum@n', 'John', 'Mendes'); -- too_much88 --