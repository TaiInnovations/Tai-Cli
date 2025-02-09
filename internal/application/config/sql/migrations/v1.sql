CREATE TABLE `session` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` TEXT NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `conversation` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `session_id` INTEGER NOT NULL,
    `role`  TEXT NOT NULL,
    `message` TEXT NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX `idx_session_id` ON `conversation` (`session_id`);

CREATE TABLE `service_provider` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` TEXT NOT NULL,
    `url` TEXT NOT NULL,
    `api_key` TEXT NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ai_model (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` TEXT NOT NULL,
    `provider_id` INTEGER NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX `idx_name` ON `ai_model` (`name`,`provider_id`);

CREATE TABLE `setting` (
    `name` TEXT NOT NULL PRIMARY KEY,
    `value` TEXT NOT NULL,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO `setting` (`name`,`value`) VALUES ('data_version', '1');
INSERT INTO `service_provider` (`id`,`name`,`url`,`api_key`) VALUES (1, 'OpenRouter','https://openrouter.ai/api/v1/chat/completions','');
INSERT INTO `ai_model` (`id`,`provider_id`,`name`) VALUES
    (1, 1, 'google/gemini-2.0-flash-lite-preview-02-05:free'),
    (2, 1, 'google/gemini-2.0-pro-exp-02-05:free'),
    (3, 1, 'google/gemini-2.0-flash-thinking-exp:free'),
    (4, 1, 'google/gemini-2.0-flash-thinking-exp-1219:free'),
    (5, 1, 'google/gemini-2.0-flash-exp:free'),
    (6, 1, 'google/gemini-exp-1206:free'),
    (7, 1, 'google/learnlm-1.5-pro-experimental:free'),
    (8, 1, 'deepseek/deepseek-r1:free'),
    (9, 1, 'deepseek/deepseek-r1-distill-llama-70b:free');
INSERT INTO `setting` (`name`,`value`) VALUES ('active_ai_model_id', '1');
INSERT INTO `session` (`name`) VALUES ('New Chat');
