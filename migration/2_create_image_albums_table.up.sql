CREATE TABLE IF NOT EXISTS image_albums(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    album_id INT NOT NULL,
    image_id INT NOT NULL,
    UNIQUE(album_id, image_id),
    CONSTRAINT fk_album
    FOREIGN KEY (album_id)
    REFERENCES albums(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_image
    FOREIGN KEY (image_id)
    REFERENCES images(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);