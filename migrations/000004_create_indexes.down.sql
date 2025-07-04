-- Drop indexes for rollback
DROP INDEX IF EXISTS idx_tweets_user_id;
DROP INDEX IF EXISTS idx_followers_follower_id;
DROP INDEX IF EXISTS idx_followers_user_id;