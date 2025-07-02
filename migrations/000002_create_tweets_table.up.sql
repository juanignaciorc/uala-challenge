CREATE TABLE tweets (
                        id UUID PRIMARY KEY,
                        user_id UUID NOT NULL REFERENCES users(id),
                        message VARCHAR(280) NOT NULL
);