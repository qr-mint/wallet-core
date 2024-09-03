BEGIN TRANSACTION;

AlTER TABLE global_notifications
    RENAME to notifications;
AlTER TABLE global_notification_translations
    RENAME to notification_translations;
AlTER TABLE global_processed_notifications
    RENAME to processed_notifications;

DROP TABLE personal_notification_translations;
DROP TABLE personal_notifications;

COMMIT;