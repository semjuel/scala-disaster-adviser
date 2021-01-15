CREATE
EXTENSION postgis;

CREATE TABLE users
(
    id    SERIAL       NOT NULL,
    uuid  VARCHAR(50)  NOT NULL,
    name  VARCHAR(255) DEFAULT NULL,
    email VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);
CREATE
UNIQUE INDEX idx_db_unique_email_users ON users (email);
CREATE
UNIQUE INDEX idx_db_unique_uuid_users ON users (uuid);


CREATE TABLE events
(
    id         SERIAL       NOT NULL,
    user_id    INT              DEFAULT NULL,
    event_id   VARCHAR(50)  NOT NULL,
    summary    VARCHAR(255) NOT NULL,
    location   VARCHAR(500) NOT NULL,
    latitude   DOUBLE PRECISION DEFAULT 0 NOT NULL,
    longitude  DOUBLE PRECISION DEFAULT 0 NOT NULL,
    start_date TIMESTAMP(0) WITHOUT TIME ZONE DEFAULT NULL,
    end_date   TIMESTAMP(0) WITHOUT TIME ZONE DEFAULT NULL,
    created    TIMESTAMP(0) WITHOUT TIME ZONE DEFAULT NULL,
    updated    TIMESTAMP(0) WITHOUT TIME ZONE DEFAULT NULL,
    PRIMARY KEY (id)
);
CREATE
UNIQUE INDEX idx_db_unique_event_id_events ON events (event_id);
SELECT AddGeometryColumn('events', 'geom', 4326, 'POINT', 2);
CREATE
INDEX idx_db_geom_events ON events USING gist(geom);


/* Add test data. */
INSERT
INTO users
VALUES (1, 'e3f0fac2-567a-11eb-ae93-0242ac130002', 'Test User #1', 'test-user@email.com',
        'ya29.a0AfH6SMBDYX29gfJYC6jJxgj8oSOSrk9fBzrWKbRkaW3uxASJMuZ3GmXqPsZ7T8T3ivvXgy1pX0W-eqKSnQRmqK7bSCwisBpUHa5AAJRd76XwgFg-0NRmviWcwyBfdaX4aAgEAf5p-X6kxQ4vScYbmhqlihaLA9SSc7Kb5N60OmI');

INSERT INTO events
VALUES (1, 1, '1v2fq4tuaoc6dc4qag0pm783b5', 'Event #1', 'Lviv, Lviv Oblast, Ukraine, 79000', 49.841952, 24.0315921,
        '2021-01-14T04:00:00+02:00', '2021-01-14T05:00:00+02:00', '2021-01-13T22:39:09.000Z',
        '2021-01-13T22:39:09.968Z');

UPDATE events
SET geom = ST_SetSRID(ST_MakePoint(longitude, latitude), 4326);
