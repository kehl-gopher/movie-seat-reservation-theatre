CREATE INDEX movies_fulltext_idx ON movies USING GIN (to_tsvector ('english', title || ' ' || synopsis));
