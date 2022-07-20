CREATE TABLE IF NOT EXISTS image_tags(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    image_id INT NOT NULL,
    tag_id INT NOT NULL,
    UNIQUE(image_id, tag_id),
    CONSTRAINT fk_image_tag
    FOREIGN KEY (image_id)
    REFERENCES images(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_tag_image
    FOREIGN KEY (tag_id)
    REFERENCES tags(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);