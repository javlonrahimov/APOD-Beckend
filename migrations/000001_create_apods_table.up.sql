CREATE TABLE IF NOT EXISTS apods (
	id bigserial PRIMARY KEY,
	title text NOT NULL,
	explanation text NOT NULL,
	media_type text NOT NULL,
	date date NOT NULL,
	url text,
	hd_url text,
	version integer NOT NULL DEFAULT 1
);
