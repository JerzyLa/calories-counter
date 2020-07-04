CREATE DATABASE IF NOT EXISTS user_data;

USE user_data;

CREATE TABLE IF NOT EXISTS users
(
    id          CHAR(36) PRIMARY KEY NOT NULL,
    account_id  CHAR(36)             NOT NULL,
    username    VARCHAR(50)          NOT NULL,
    password    VARCHAR(50)          NOT NULL,
    role_id     INT                  NOT NULL,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT unique_account UNIQUE (account_id, username)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = latin1;


CREATE TABLE IF NOT EXISTS users_meals
(
    id          CHAR(36) PRIMARY KEY NOT NULL,
    user_id     CHAR(36)             NOT NULL,
    name        VARCHAR(50)          NOT NULL,
    date        DATE                 NOT NULL,
    time        TIME                 NOT NULL,
    calories    INT                  NOT NULL,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = latin1;


CREATE TABLE IF NOT EXISTS users_settings
(
    id                      INT      NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_id                 CHAR(36) NOT NULL,
    expected_daily_calories INT      NOT NULL,
    create_time             TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_time             TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT unique_user_id UNIQUE (user_id)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = latin1;


CREATE TABLE IF NOT EXISTS users_calories
(
    id               INT      NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_id          CHAR(36) NOT NULL,
    date             DATE     NOT NULL,
    total_calories   INT      NOT NULL,
    calories_deficit TINYINT  NOT NULL,
    create_time      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_time      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT unique_user_id_date UNIQUE (user_id, date)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = latin1;
