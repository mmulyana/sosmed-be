CREATE TABLE IF NOT EXISTS comments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    postId INT NOT NULL,
    userId INT NOT NULL,
    content TEXT NOT NULL,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_post FOREIGN KEY (postId)
        REFERENCES posts(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_user_comment FOREIGN KEY (userId)
        REFERENCES users(id)
        ON DELETE CASCADE
);
