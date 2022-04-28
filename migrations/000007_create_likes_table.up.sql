CREATE TABLE IF NOT EXISTS likes (
	id bigint NOT NULL,
	user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
	apod_id bigint NOT NULL REFERENCES apods ON DELETE CASCADE,
	PRIMARY KEY (id, user_id, apod_id)
);
