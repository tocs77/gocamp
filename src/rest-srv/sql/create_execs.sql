use classes;
CREATE TABLE IF NOT EXISTS execs(
  id int auto_increment primary key,
  first_name varchar(255) NOT NULL,
  last_name varchar(255) NOT NULL,
  email varchar(255) NOT NULL UNIQUE,
  username varchar(255) NOT NULL UNIQUE,
  password varchar(255) NOT NULL,
  password_changed_at varchar(255),
  user_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  password_reset_token varchar(255),
  password_token_expires varchar(255),
  inactive_status boolean NOT NULL DEFAULT false,
  role varchar(50) NOT NULL,
  INDEX idx_email(email),
  INDEX idx_username(username)
) auto_increment=100;