CREATE TABLE IF NOT EXISTS teams
(
    team_name
    TEXT
    PRIMARY
    KEY
);

CREATE TABLE IF NOT EXISTS users
(
    user_id TEXT PRIMARY KEY, username TEXT NOT NULL,
    team_name TEXT NOT NULL REFERENCES team ( team_nam ) ON DELETE RESTRICT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS pull_requests
(
    pull_request_idTEXTPRIMARYKEY,
    pull_request_name TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    merged_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS reviewers
(
    id BIGSERIAL PRIMARY KEY,
    pull_request_id TEXT NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    UNIQUE(pull_request_id, user_id)
);
