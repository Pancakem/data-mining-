CREATE TABLE posts(
    id VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    post_text VARCHAR(8000),
    towhom VARCHAR(8000),
    time_stamp VARCHAR(255),
    likes INTEGER,
    PRIMARY KEY (id)
);
