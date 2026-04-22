CREATE TABLE IF NOT EXISTS scenes(
    id INT AUTO_INCREMENT PRIMARY KEY,
    story_id INT,
    title VARCHAR(255),
    description TEXT,
    question TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    hint TEXT,
    next_scene_id INT,
    failure_scene_id INT,
    FOREIGN KEY (story_id) REFERENCES stories(id)
);
