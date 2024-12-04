ALTER TABLE users
ADD notify_application_updated BOOLEAN NOT NULL DEFAULT true,
ADD notify_application_received BOOLEAN NOT NULL DEFAULT true,
ADD notify_application_withdrawn BOOLEAN NOT NULL DEFAULT true;
