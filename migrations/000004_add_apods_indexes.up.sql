CREATE INDEX IF NOT EXISTS apods_title_idx ON apods USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS apods_copyright_idx ON apods USING GIN (to_tsvector('simple', copyright));