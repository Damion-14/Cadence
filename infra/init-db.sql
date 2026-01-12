-- Cadence Workout Logger - Database Schema
-- This script initializes the PostgreSQL database with all required tables and indexes

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    username VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Workout sessions table
CREATE TABLE IF NOT EXISTS workout_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT chk_status CHECK (status IN ('active', 'completed'))
);

CREATE INDEX IF NOT EXISTS idx_workout_sessions_user_id ON workout_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_workout_sessions_status ON workout_sessions(status);
CREATE INDEX IF NOT EXISTS idx_workout_sessions_completed_at ON workout_sessions(completed_at) WHERE completed_at IS NOT NULL;

-- Exercises table (individual exercises within a workout)
CREATE TABLE IF NOT EXISTS exercises (
    id SERIAL PRIMARY KEY,
    workout_session_id INTEGER NOT NULL REFERENCES workout_sessions(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_exercises_workout_session_id ON exercises(workout_session_id);
CREATE INDEX IF NOT EXISTS idx_exercises_name ON exercises(name);

-- Sets table (individual sets within an exercise)
CREATE TABLE IF NOT EXISTS sets (
    id SERIAL PRIMARY KEY,
    exercise_id INTEGER NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    set_number INTEGER NOT NULL,
    reps INTEGER NOT NULL,
    weight DECIMAL(10, 2),
    is_bodyweight BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT chk_reps CHECK (reps > 0),
    CONSTRAINT chk_weight CHECK (weight IS NULL OR weight >= 0)
);

CREATE INDEX IF NOT EXISTS idx_sets_exercise_id ON sets(exercise_id);

-- Insert seed data for development (password is "password123")
-- Password hash generated with bcrypt cost 10
INSERT INTO users (email, password_hash, username) VALUES
    ('demo@cadence.app', '$2a$10$rXKaFWFkYaH.wIJ0HvZ0EeJ6Y8xQqVVGGY5cCWK9O3KGxqZN1QY9G', 'demouser')
ON CONFLICT (email) DO NOTHING;

-- Insert sample completed workout for demo user
DO $$
DECLARE
    demo_user_id INTEGER;
    workout_id INTEGER;
    bench_press_id INTEGER;
    squats_id INTEGER;
BEGIN
    -- Get demo user ID
    SELECT id INTO demo_user_id FROM users WHERE email = 'demo@cadence.app';

    IF demo_user_id IS NOT NULL THEN
        -- Create a completed workout from 2 days ago
        INSERT INTO workout_sessions (user_id, name, status, started_at, completed_at)
        VALUES (
            demo_user_id,
            'Push Day',
            'completed',
            NOW() - INTERVAL '2 days',
            NOW() - INTERVAL '2 days' + INTERVAL '45 minutes'
        )
        RETURNING id INTO workout_id;

        -- Add Bench Press exercise
        INSERT INTO exercises (workout_session_id, name, order_index)
        VALUES (workout_id, 'Bench Press', 0)
        RETURNING id INTO bench_press_id;

        -- Add sets for Bench Press
        INSERT INTO sets (exercise_id, set_number, reps, weight, is_bodyweight)
        VALUES
            (bench_press_id, 1, 10, 135.00, false),
            (bench_press_id, 2, 10, 135.00, false),
            (bench_press_id, 3, 8, 135.00, false);

        -- Add Squats exercise
        INSERT INTO exercises (workout_session_id, name, order_index)
        VALUES (workout_id, 'Squats', 1)
        RETURNING id INTO squats_id;

        -- Add sets for Squats
        INSERT INTO sets (exercise_id, set_number, reps, weight, is_bodyweight)
        VALUES
            (squats_id, 1, 12, 225.00, false),
            (squats_id, 2, 10, 225.00, false),
            (squats_id, 3, 10, 225.00, false);

        -- Create another completed workout from 1 day ago
        INSERT INTO workout_sessions (user_id, name, status, started_at, completed_at)
        VALUES (
            demo_user_id,
            'Pull Day',
            'completed',
            NOW() - INTERVAL '1 day',
            NOW() - INTERVAL '1 day' + INTERVAL '50 minutes'
        )
        RETURNING id INTO workout_id;

        -- Add Deadlifts exercise
        INSERT INTO exercises (workout_session_id, name, order_index)
        VALUES (workout_id, 'Deadlifts', 0)
        RETURNING id INTO bench_press_id;

        -- Add sets for Deadlifts
        INSERT INTO sets (exercise_id, set_number, reps, weight, is_bodyweight)
        VALUES
            (bench_press_id, 1, 5, 315.00, false),
            (bench_press_id, 2, 5, 315.00, false),
            (bench_press_id, 3, 5, 315.00, false);

        -- Add Pull-ups exercise (bodyweight)
        INSERT INTO exercises (workout_session_id, name, order_index)
        VALUES (workout_id, 'Pull-ups', 1)
        RETURNING id INTO squats_id;

        -- Add sets for Pull-ups
        INSERT INTO sets (exercise_id, set_number, reps, weight, is_bodyweight)
        VALUES
            (squats_id, 1, 10, NULL, true),
            (squats_id, 2, 10, NULL, true),
            (squats_id, 3, 8, NULL, true);
    END IF;
END $$;
