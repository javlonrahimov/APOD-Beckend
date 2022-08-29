CREATE TABLE IF NOT EXISTS apods (
	id bigserial PRIMARY KEY,
	date date NOT NULL,
	title text NOT NULL,
	explanation text NOT NULL,
	url text NOT NULL,
	hd_url text NOT NULL,
	media_type text NOT NULL,
	created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
);