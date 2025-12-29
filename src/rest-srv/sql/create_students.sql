use classes;
CREATE TABLE if not exists students(
  id int auto_increment primary key,
  first_name varchar(255) not null,
  last_name varchar(255) not null,
  email varchar(255) not null unique,
  class varchar(255) not null,
  index(email),
  FOREIGN KEY (class) REFERENCES teachers(class)
) auto_increment=100;