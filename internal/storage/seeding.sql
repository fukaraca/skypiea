
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

WITH new_users AS (
    INSERT INTO users
        (user_uuid, firstname, lastname, email, password, role, status)
        VALUES
            (uuid_generate_v4(), 'John', 'Tester', 'johndoe@example.com', '$2a$04$jLaYijt7BEcQ2ZRygQdEpe347jdq38zGNQ5QabY9E6FuiI6nCpQq2', 'admin', 'active'),
            (uuid_generate_v4(), 'Jane',   'Tester', 'janedoe@example.com',   '$2a$04$jLaYijt7BEcQ2ZRygQdEpe347jdq38zGNQ5QabY9E6FuiI6nCpQq2', 'user_std', 'active')
        RETURNING id, user_uuid, firstname
),


     convs AS (
         INSERT INTO conversations
             (user_uuid, title, metadata)
             SELECT
                 u.user_uuid,
                 format('Chat #%s with %s', gs, u.firstname),
                 NULL
             FROM new_users u
                      CROSS JOIN generate_series(1, 5) AS gs
             RETURNING id, created_at
     )


INSERT INTO messages
(conv_id, model_id, by_user, message, metadata, created_at)
SELECT
    c.id,
    'gpt-4o' AS model_id,
    (m % 2 = 1) AS by_user,
    CASE WHEN m % 2 = 1
             THEN format('Lorem ipsum question %s?', m)
         ELSE format('Lorem ipsum answer %s.',   m)
        END AS message,
    NULL,
    c.created_at + ((m - 1) * INTERVAL '10 seconds')
FROM convs c
         CROSS JOIN generate_series(1, 10) AS m;
