CREATE TABLE IF NOT EXISTS user_progress(
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    scene_id INT,
    attempts INT DEFAULT 0,
    completed TINYINT(1) DEFAULT 0,
    UNIQUE KEY unique_user_scene (user_id, scene_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (scene_id) REFERENCES scenes(id)
);
