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

-- Create indexes for read optimization
-- Index on tweets.user_id for efficient queries when getting user's tweets
CREATE INDEX idx_tweets_user_id ON tweets(user_id);

-- Index on followers.follower_id for efficient queries when getting who a user follows
CREATE INDEX idx_followers_follower_id ON followers(follower_id);

-- Index on followers.user_id for efficient queries when getting a user's followers
CREATE INDEX idx_followers_user_id ON followers(user_id);
