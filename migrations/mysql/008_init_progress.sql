CREATE TABLE progress (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    story_id INTEGER NOT NULL REFERENCES stories(id),
    scene_id INTEGER NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для ускорения запросов
CREATE INDEX idx_progress_user_id ON progress(user_id);
CREATE INDEX idx_progress_story_id ON progress(story_id);
CREATE INDEX idx_progress_user_story ON progress(user_id, story_id);
