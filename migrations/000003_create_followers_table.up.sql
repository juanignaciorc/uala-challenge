CREATE TABLE followers (
                           user_id UUID NOT NULL REFERENCES users(id),
                           follower_id UUID NOT NULL REFERENCES users(id),
                           PRIMARY KEY(user_id, follower_id)
);