/*CREATE INDEX books_idx ON books USING GIN ((setweight(to_tsvector(config_name, coalesce(author,'')), 'A') ||
       setweight(to_tsvector(config_name, coalesce(title,'')), 'B') ||
       setweight(to_tsvector(config_name, coalesce(body,'')), 'C')));
*/

ALTER TABLE books
    ADD COLUMN textsearchable_index_col tsvector
        GENERATED ALWAYS AS (setweight(to_tsvector(config_name, coalesce(author,'')), 'A') ||
                             setweight(to_tsvector(config_name, coalesce(title,'')), 'B') ||
                             setweight(to_tsvector(config_name, coalesce(body,'')), 'C')) STORED;

CREATE INDEX textsearch_idx ON books USING GIN (textsearchable_index_col);
