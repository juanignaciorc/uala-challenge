-- Initialize database schema for microblogging platform

-- Create users table
CREATE TABLE "users" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" varchar NOT NULL,
    "email" varchar NOT NULL
);

-- Create tweets table
CREATE TABLE tweets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    message VARCHAR(280) NOT NULL
);

-- Create followers table
CREATE TABLE followers (
    user_id UUID NOT NULL REFERENCES users(id),
    follower_id UUID NOT NULL REFERENCES users(id),
    PRIMARY KEY(user_id, follower_id)
);