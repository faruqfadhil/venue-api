create database if not exists venue_db;

use venue_db;

CREATE TABLE IF NOT EXISTS `venue` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'venue identifier',
  `name` TEXT NOT NULL,
  `min_price` DECIMAL(15, 2) NOT NULL DEFAULT 0,
  `max_price` DECIMAL(15, 2) NOT NULL DEFAULT 0,
  `capacity` int(11) NOT NULL DEFAULT 0,
  `star` DECIMAL(15, 2) NOT NULL DEFAULT 0,
  `review_count` int(11) NOT NULL DEFAULT 0,
  `thumbnail_url` TEXT NOT NULL,
  `city_id` int(11) NOT NULL,
  `description` TEXT NOT NULL,
  `website` TEXT NOT NULL,
  `phone` varchar(13) NOT NULL DEFAULT '',
  `email` varchar(255) NOT NULL DEFAULT '',
  `instagram` varchar(255) NOT NULL DEFAULT '',
  `address` TEXT NOT NULL,
  `logo` TEXT NOT NULL,
  `is_favourite` TINYINT(1) NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created date',
  `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'user who create this entity',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'update date',
  `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'user who update this entity',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `city` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'city identifier',
  `name` TEXT NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created date',
  `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'user who create this entity',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'update date',
  `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'user who update this entity',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `venue_gallery` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `venue_id` int(11) NOT NULL,
  `file_url` TEXT NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created date',
  `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'user who create this entity',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'update date',
  `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'user who update this entity',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE IF NOT EXISTS `venue_category_package` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `venue_id` int(11) NOT NULL,
  `description` TEXT NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created date',
  `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'user who create this entity',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'update date',
  `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'user who update this entity',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `category_package` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `category_id` int(11) NOT NULL,
  `name` TEXT NOT NULL,
  `thumbnail_url` TEXT NOT NULL,
  `price` DECIMAL(15, 2) NOT NULL DEFAULT 0,
  `capacity` int(11) NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created date',
  `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'user who create this entity',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'update date',
  `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'user who update this entity',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `auth` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `email` TEXT NOT NULL,
  `fullname` TEXT NOT NULL,
  `password` TEXT NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created date',
  `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'user who create this entity',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'update date',
  `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'user who update this entity',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `order` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `package_id`int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `date` timestamp NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created date',
  `created_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'user who create this entity',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'update date',
  `updated_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'user who update this entity',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- SEEDER
INSERT INTO city(id,name,created_at,created_by,updated_at,updated_by) VALUES
(1,'Surabaya',NOW(),'user',NOW(),'user'),
(2,'Sidoarjo',NOW(),'user',NOW(),'user'),
(3,'Malang',NOW(),'user',NOW(),'user'),
(4,'Kediri',NOW(),'user',NOW(),'user'),
(5,'Mojokerto',NOW(),'user',NOW(),'user'),
(6,'Blitar',NOW(),'user',NOW(),'user'),
(7,'Bangkalan',NOW(),'user',NOW(),'user'),
(8,'Sumenep',NOW(),'user',NOW(),'user'),
(9,'Sampang',NOW(),'user',NOW(),'user'),
(10,'Pamekasan',NOW(),'user',NOW(),'user'),
(11,'Gersik',NOW(),'user',NOW(),'user'),
(12,'Jember',NOW(),'user',NOW(),'user'),
(13,'Probolinggo',NOW(),'user',NOW(),'user'),
(14,'Jombang',NOW(),'user',NOW(),'user'),
(15,'Banyuwangi',NOW(),'user',NOW(),'user'),
(16,'Situbondo',NOW(),'user',NOW(),'user');