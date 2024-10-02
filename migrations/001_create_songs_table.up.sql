CREATE TABLE songs
(
    id           SERIAL PRIMARY KEY,
    group_name   VARCHAR(100) NOT NULL,
    song_name    VARCHAR(100) NOT NULL,
    lyrics       TEXT         NOT NULL,
    release_date DATE,
    link         VARCHAR(255)
);
