-- +goose Up
CREATE TABLE group_members (
  user_id UUID REFERENCES users (id) ON DELETE CASCADE,
  group_id UUID REFERENCES groups (id) ON DELETE CASCADE,
  UNIQUE(user_id, group_id)
);

-- +goose Down
DROP TABLE group_members;
