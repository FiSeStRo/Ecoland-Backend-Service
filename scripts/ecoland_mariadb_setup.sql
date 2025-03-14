SET FOREIGN_KEY_CHECKS=0;

DROP TABLE IF EXISTS rel_building_def_production;
DROP TABLE IF EXISTS rel_building_product;
DROP TABLE IF EXISTS def_rel_building_production;
DROP TABLE IF EXISTS def_rel_production_product;
DROP TABLE IF EXISTS user_resources;
DROP TABLE IF EXISTS buildings;
DROP TABLE IF EXISTS def_buildings;
DROP TABLE IF EXISTS def_production;
DROP TABLE IF EXISTS def_product;
DROP TABLE IF EXISTS users;

SET FOREIGN_KEY_CHECKS=1;

CREATE TABLE `users` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `username` VARCHAR(255) UNIQUE NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `role` INT DEFAULT 0,
  `time_created` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `time_last_activity` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user_resources` (
  `user_id` INT PRIMARY KEY,
  `money` DECIMAL(10,2) NOT NULL DEFAULT 100000.00,
  `prestige` INT DEFAULT 0,
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `def_buildings` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `token_name` VARCHAR(255) NOT NULL,
  `base_construction_cost` DECIMAL(10,2) NOT NULL,
  `base_construction_time` INT NOT NULL 
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `def_product` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `token_name` VARCHAR(255) NOT NULL,
  `base_value` INT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `def_production` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `token_name` VARCHAR(255) NOT NULL,
  `cost` DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  `base_duration` INT NOT NULL 
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `def_rel_production_product` (
  `production_id` INT NOT NULL,
  `product_id` INT NOT NULL,
  `is_input` BOOLEAN NOT NULL COMMENT 'True if consumed, False if produced',
  `amount` INT NOT NULL,
  PRIMARY KEY (`production_id`, `product_id`, `is_input`),
  FOREIGN KEY (`production_id`) REFERENCES `def_production` (`id`),
  FOREIGN KEY (`product_id`) REFERENCES `def_product` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `def_rel_building_production` (
  `building_id` INT NOT NULL,
  `production_id` INT NOT NULL,
  PRIMARY KEY (`building_id`, `production_id`),
  FOREIGN KEY (`building_id`) REFERENCES `def_buildings` (`id`),
  FOREIGN KEY (`production_id`) REFERENCES `def_production` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `buildings` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `user_id` INT NOT NULL,
  `building_def_id` INT NOT NULL,
  `name` VARCHAR(255),
  `status` VARCHAR(50) NOT NULL DEFAULT 'under_construction',
  `construction_start_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `construction_end_time` TIMESTAMP NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  FOREIGN KEY (`building_def_id`) REFERENCES `def_buildings` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `rel_building_product` (
  `building_id` INT NOT NULL,
  `product_id` INT NOT NULL,
  `amount` INT NOT NULL DEFAULT 0,
  `capacity` INT NOT NULL DEFAULT 500,
  PRIMARY KEY (`building_id`, `product_id`),
  FOREIGN KEY (`building_id`) REFERENCES `buildings` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`product_id`) REFERENCES `def_product` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `rel_building_def_production` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `building_id` INT NOT NULL,
  `production_id` INT NOT NULL,
  `time_start` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `time_end` TIMESTAMP NOT NULL,
  `cycles` INT NOT NULL DEFAULT 1,
  `is_completed` BOOLEAN DEFAULT FALSE,
  FOREIGN KEY (`building_id`) REFERENCES `buildings` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`production_id`) REFERENCES `def_production` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE INDEX idx_buildings_user ON buildings(user_id);
CREATE INDEX idx_rel_building_def_production_status ON rel_building_def_production(is_completed);