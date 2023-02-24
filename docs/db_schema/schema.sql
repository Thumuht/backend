CREATE TABLE user (
  user_id     INTEGER PRIMARY KEY AUTOINCREMENT,
  login_name  TEXT UNIQUE NOT NULL,
  nickname    TEXT,
  pwhash      TEXT,
  email       TEXT,
  about       TEXT
);

CREATE TABLE post (
  post_id    INTEGER PRIMARY KEY AUTOINCREMENT,
  title      VARCHAR(100),
  content    TEXT,
  user_id    INTEGER REFERENCES user ON DELETE CASCADE
);

CREATE TABLE comment (
  comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
  content    TEXT,
  user_id    INTEGER REFERENCES user ON DELETE CASCADE,
  post_id    INTEGER REFERENCES post ON DELETE CASCADE
);
