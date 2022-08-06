CREATE TABLE IF NOT EXISTS images(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    path VARCHAR(256) NOT NULL UNIQUE,
    width INT,
    height INT,
    thumbnail VARCHAR(256) NOT NULL UNIQUE,
    width_thumb INT,
    height_thumb INT,
    title VARCHAR(32) DEFAULT NULL,
    description VARCHAR(256) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);