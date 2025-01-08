-- create slugify function, lower case, remove accents, remove special characters, remove spaces, remove duplicate -
CREATE FUNCTION slugify(text) RETURNS text AS $$
    SELECT lower(regexp_replace($1, '[^a-zA-Z0-9]+', '-', 'g'));
$$ LANGUAGE sql;

ALTER TABLE categories ADD COLUMN handle VARCHAR(255);

-- set handle as name slug
UPDATE categories SET handle = slugify(name) WHERE handle IS NULL;
