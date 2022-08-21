ALTER TABLE albums 
    ADD CONSTRAINT fk_cover
    FOREIGN KEY (cover_id)
    REFERENCES images(id)
        ON UPDATE CASCADE
        ON DELETE SET NULL;