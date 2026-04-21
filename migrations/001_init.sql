-- migrations/001_init.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    progress INTEGER DEFAULT 0,
    score INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stories (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE scenes (
    id SERIAL PRIMARY KEY,
    story_id INTEGER REFERENCES stories(id),
    title VARCHAR(255),
    description TEXT,
    question TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    hint TEXT,
    next_scene_id INTEGER,
    failure_scene_id INTEGER
);

CREATE TABLE user_progress (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    scene_id INTEGER REFERENCES scenes(id),
    attempts INTEGER DEFAULT 0,
    completed BOOLEAN DEFAULT FALSE,
    UNIQUE(user_id, scene_id)
);