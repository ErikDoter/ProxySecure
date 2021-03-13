create table requests
(
    id      bigserial primary key,
    method  TEXT,
    url     TEXT,
    headers TEXT,
    body    TEXT
);