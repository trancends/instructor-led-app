CREATE DATABASE instructor_led_app_db;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE status_type AS ENUM('PROCESS','FINISH');
CREATE TYPE role_type AS ENUM('ADMIN','PARTICIPANT','TRAINER');

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL ,
    password VARCHAR(50) NOT NULL,
    role role_type NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE schedules (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    documentation text,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE questions (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    schedule_id uuid NOT NULL,
    description TEXT NOT NULL,
    status status_type DEFAULT 'PROCESS',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
 
);

CREATE TABLE attendances (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    schedule_id uuid NOT NULL,
    created_at timestamp  DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp  DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp
);

ALTER TABLE schedules ADD CONSTRAINT "schedules_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE questions ADD CONSTRAINT "questions_schedule_id_fkey" FOREIGN KEY (schedule_id) REFERENCES schedules(id);
ALTER TABLE attendances ADD CONSTRAINT "attendances_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE attendances ADD CONSTRAINT "attendances_schedule_id_fkey" FOREIGN KEY (schedule_id) REFERENCES schedules(id);
