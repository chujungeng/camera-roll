CREATE TABLE IF NOT EXISTS album_tags(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    album_id INT NOT NULL,
    tag_id INT NOT NULL,
    UNIQUE(album_id, tag_id),
    CONSTRAINT fk_album_tag
    FOREIGN KEY (album_id)
    REFERENCES albums(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_tag_album
    FOREIGN KEY (tag_id)
    REFERENCES tags(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);