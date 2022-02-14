CREATE TABLE users (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
	name        TEXT NOT NULL,
	email       TEXT UNIQUE,
	api_key     TEXT NOT NULL UNIQUE,
	created_at  TEXT NOT NULL,
	updated_at  TEXT NOT NULL
);

CREATE TABLE auths (
	id              INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id         INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	source          TEXT NOT NULL,
	source_id       TEXT NOT NULL,
	access_token    TEXT NOT NULL,
	refresh_token   TEXT NOT NULL,
	expiry          TEXT,
	created_at      TEXT NOT NULL,
	updated_at      TEXT NOT NULL,

	UNIQUE(user_id, source),  -- one source per user
	UNIQUE(source, source_id)
);

CREATE TABLE blogs (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    title       TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at  TEXT NOT NULL,
    updated_at  TEXT NOT NULL,

    UNIQUE(title)
);

CREATE TABLE sub_blogs (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    blog_id     INTEGER NOT NULL REFERENCES blogs (id) ON DELETE CASCADE,
    title       TEXT NOT NULL,
    content     TEXT NOT NULL,
    created_at  TEXT NOT NULL,
    updated_at  TEXT NOT NULL
);

CREATE TABLE comments (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    sub_blog_id INTEGER NOT NULL REFERENCES sub_blogs (id) ON DELETE CASCADE,
    user_id     INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    content     TEXT NOT NULL,
    created_at  TEXT NOT NULL
);

CREATE TABLE blog_subscriptions (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    blog_id     INTEGER NOT NULL REFERENCES blogs (id) ON DELETE CASCADE
);

CREATE TABLE sub_blog_subscriptions (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    sub_blog_id INTEGER NOT NULL REFERENCES sub_blogs (id) ON DELETE CASCADE,
    user_id     INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE projects (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    description TEXT,
    html_url    TEXT NOT NULL,

    UNIQUE(name)
);

-- many 2 many relation.
CREATE TABLE projects_topics (
    project_id              INTEGER NOT NULL REFERENCES projects (id),
    topic_description_id    INTEGER NOT NULL REFERENCES topics_description (id) ON DELETE CASCADE
);

CREATE TABLE topics_description (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL
);