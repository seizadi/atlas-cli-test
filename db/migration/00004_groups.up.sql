
CREATE TABLE groups (
  id serial primary key,
  created_at timestamptz DEFAULT current_timestamp,
  updated_at timestamptz DEFAULT NULL,
  name varchar(255) DEFAULT NULL,
  description varchar(255) DEFAULT NULL,
  account_id int REFERENCES accounts(id) ON DELETE SET NULL

);

CREATE TRIGGER groups_updated_at
  BEFORE UPDATE OR INSERT ON groups
  FOR EACH ROW
  EXECUTE PROCEDURE set_updated_at();

