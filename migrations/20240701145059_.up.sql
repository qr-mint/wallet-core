BEGIN TRANSACTION;

AlTER TABLE notifications
    RENAME TO global_notifications;
AlTER TABLE notification_translations
    RENAME TO global_notification_translations;
AlTER TABLE processed_notifications
    RENAME TO global_processed_notifications;

ALTER TABLE global_notifications DROP COLUMN IF EXISTS type;

CREATE TABLE personal_notifications
(
    id      BIGINT                                                          NOT NULl PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id BIGINT                                                          NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE personal_notification_translations
(
    id              BIGINT                                       NOT NULl PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    language        VARCHAR(10) CHECK (language IN ('ru', 'en')) NOT NULL,
    text            TEXT                                         NOT NULL,
    image_path      TEXT DEFAULT NULL,
    notification_id BIGINT                                       NOT NULL,
    CONSTRAINT fk_notification_id FOREIGN KEY (notification_id) REFERENCES personal_notifications (id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX unique_per_notification_translations_notification_id_language ON personal_notification_translations (notification_id, language);

COMMIT;