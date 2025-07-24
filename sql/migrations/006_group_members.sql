-- +goose Up
CREATE TABLE group_members (
  user_id UUID,
  group_id UUID,
  UNIQUE(user_id, group_id)
);

-- +goose Down
DROP TABLE group_members;
