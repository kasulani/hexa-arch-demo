CREATE TABLE url
(
    id         BIGSERIAL    NOT NULL,
    code       VARCHAR(50)  NOT NULL,
    raw_url    VARCHAR(250) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX url_short_url ON url (code);
