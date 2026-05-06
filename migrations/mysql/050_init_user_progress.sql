CREATE TABLE IF NOT EXISTS user_progress(
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    story_id INTEGER NOT NULL REFERENCES stories(id),
    scene_id INTEGER NOT NULL,
    attempts INT DEFAULT 0,
    completed TINYINT(1) DEFAULT 0,
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_user_scene (user_id, scene_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (scene_id) REFERENCES scenes(id)
);
