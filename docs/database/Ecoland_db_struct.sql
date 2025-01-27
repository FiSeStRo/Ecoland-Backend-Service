CREATE TABLE `users` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `username` varchar(255) UNIQUE NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `role` int,
  `time_created` timestamp,
  `time_last_activity` timestamp
);

CREATE TABLE `user_resources` (
  `user_id` int PRIMARY KEY,
  `money` decimal NOT NULL DEFAULT 100000,
  `prestige` int
);

CREATE TABLE `buildings` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `def_id` int NOT NULL,
  `name` varchar(255),
  `time_build` timestamp NOT NULL
);

CREATE TABLE `def_products` (
  `id` int PRIMARY KEY,
  `base_value` int NOT NULL,
  `token_name` sting
);

CREATE TABLE `def_production` (
  `id` int PRIMARY KEY,
  `token_name` varchar(255),
  `cost` decimal,
  `base_duration` int
);

CREATE TABLE `def_rel_production_product` (
  `production_id` int,
  `product_id` int,
  `is_input` bool,
  `amount` int
);

CREATE TABLE `def_buildings` (
  `id` int PRIMARY KEY,
  `token_name` varchar(255),
  `base_construction_cost` decimal,
  `base_construction_time` int
);

CREATE TABLE `def_rel_building_production` (
  `building_id` int NOT NULL,
  `production_id` int NOT NULL
);

CREATE TABLE `rel_building_product` (
  `building_id` int,
  `product_id` int,
  `amount` int,
  `capacity` int NOT NULL DEFAULT 500
);

CREATE TABLE `rel_buildng_def_production` (
  `id` int PRIMARY KEY,
  `building_id` int,
  `production_id` int,
  `time_sart` timestamp NOT NULL,
  `time_end` timestamp NOT NULL,
  `cycles` int NOT NULL,
  `is_completed` bool DEFAULT false
);

ALTER TABLE `user_resources` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `buildings` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `buildings` ADD FOREIGN KEY (`def_id`) REFERENCES `def_buildings` (`id`);

ALTER TABLE `def_rel_production_product` ADD FOREIGN KEY (`production_id`) REFERENCES `def_production` (`id`);

ALTER TABLE `def_rel_production_product` ADD FOREIGN KEY (`product_id`) REFERENCES `def_products` (`id`);

ALTER TABLE `def_rel_building_production` ADD FOREIGN KEY (`building_id`) REFERENCES `def_buildings` (`id`);

ALTER TABLE `def_production` ADD FOREIGN KEY (`id`) REFERENCES `def_rel_building_production` (`production_id`);

ALTER TABLE `rel_building_product` ADD FOREIGN KEY (`building_id`) REFERENCES `buildings` (`id`);

ALTER TABLE `rel_building_product` ADD FOREIGN KEY (`product_id`) REFERENCES `def_products` (`id`);

ALTER TABLE `rel_buildng_def_production` ADD FOREIGN KEY (`building_id`) REFERENCES `buildings` (`id`);

ALTER TABLE `rel_buildng_def_production` ADD FOREIGN KEY (`production_id`) REFERENCES `def_production` (`id`);

ALTER TABLE `def_rel_building_production` ADD FOREIGN KEY (`production_id`) REFERENCES `def_rel_building_production` (`building_id`);
