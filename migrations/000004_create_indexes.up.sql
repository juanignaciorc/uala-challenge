-- Create indexes for read optimization
-- Index on tweets.user_id for efficient queries when getting user's tweets
CREATE INDEX idx_tweets_user_id ON tweets(user_id);

-- Index on followers.follower_id for efficient queries when getting who a user follows
CREATE INDEX idx_followers_follower_id ON followers(follower_id);

-- Index on followers.user_id for efficient queries when getting a user's followers
CREATE INDEX idx_followers_user_id ON followers(user_id);