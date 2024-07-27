-- ----------------------------
-- Table structure for article
-- ----------------------------
CREATE TABLE IF NOT EXISTS `article` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT ,
  `title` varchar(255)  NOT NULL,
  `summary` varchar(255)  NOT NULL DEFAULT '',
  `content` text  NOT NULL,
  `author` varchar(255)  NOT NULL DEFAULT '',
  `type_id` int NOT NULL DEFAULT '0',
  `type_name` varchar(255)  NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS  title_idx on article(title);
CREATE INDEX IF NOT EXISTS  created_at_idx on article(created_at);

-- ----------------------------
-- Table structure for site_config
-- ----------------------------
CREATE TABLE IF NOT EXISTS `site_config` (
  `id`  INTEGER PRIMARY KEY  AUTOINCREMENT,
  `domain` varchar(100)  NOT NULL,
  `index_title` varchar(100)  NOT NULL DEFAULT '',
  `index_keywords` varchar(255)  NOT NULL DEFAULT '',
  `index_description` varchar(255)  NOT NULL DEFAULT '',
  `template_name` varchar(100)  NOT NULL DEFAULT '',
  `routes` text  NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS domain_uni on site_config(domain);

