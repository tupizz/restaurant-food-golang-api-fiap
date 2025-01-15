CREATE FUNCTION slugify(text) RETURNS text AS $$
    SELECT lower(regexp_replace($1, '[^a-zA-Z0-9]+', '-', 'g'));
$$ LANGUAGE sql;

ALTER TABLE categories ADD COLUMN handle VARCHAR(255);

UPDATE categories SET handle = slugify(name) WHERE handle IS NULL;

ALTER TABLE categories ADD CONSTRAINT unique_handle UNIQUE (handle);

ALTER TABLE categories ALTER COLUMN handle SET NOT NULL;
