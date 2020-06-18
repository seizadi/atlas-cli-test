CREATE TABLE users_groups
(
    user_id int REFERENCES users(id) ON DELETE CASCADE,
    group_id int REFERENCES groups(id) ON DELETE CASCADE,
    primary key (user_id, group_id)
);
