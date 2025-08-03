-- Create the Todos Table
CREATE TABLE IF NOT EXISTS todos_db.todos (
  id UUID PRIMARY KEY,
  content TEXT,
  due TIMESTAMP,
  priority INT,
  completed BOOLEAN
);
