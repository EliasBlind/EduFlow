-- name: CreatePerson :one
-- Создает нового пользователя и возвращает всю строку
INSERT INTO person (
    email,
    username,
    password_hash
) VALUES (
    $1, $2, $3
)
RETURNING id, email, username, password_hash, user_role;

-- name: CreateUsers :copyfrom
INSERT INTO person (
    id,
    email,
    username,
    password_hash,
    user_role
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetPersonByLogin :one
-- Поиск для авторизации
SELECT
    id,
    email,
    username,
    password_hash,
    user_role
FROM person
WHERE username = $1
ORDER BY id
LIMIT 1;

-- name: GetPersonById :one
-- Поиск профиля
SELECT
    id,
    email,
    username,
    password_hash,
    user_role
FROM person
WHERE id = $1
ORDER BY id
LIMIT 1;

-- name: UserExist :one
SELECT
    EXISTS(
        SELECT 1
        FROM person
        WHERE username = $1
    );

-- name: UpdatePassword :exec
-- Смена пароля
UPDATE person
SET password_hash = $2
WHERE id = $1;

-- name: SaveRefreshToken :one
-- Сохраняет новую сессию рефреш-токена
INSERT INTO refresh_sessions (
    user_id,
    app_id,
    token_id,
    expires_at
) VALUES (
    $1, $2, $3, $4
)
RETURNING id;

-- name: GetSessionByTokenID :one
-- Проверка токена при обновлении
SELECT
    id,
    user_id,
    app_id,
    token_id,
    expires_at
FROM refresh_sessions
WHERE
    token_id = $1
    AND expires_at > NOW()
ORDER BY expires_at
LIMIT 1;

-- name: DeleteSessionByTokenID :exec
-- Удаление токена (Logout или Refresh Token Rotation)
DELETE FROM refresh_sessions
WHERE token_id = $1;

-- name: DeleteAllUserSessions :exec
DELETE FROM refresh_sessions
WHERE user_id = $1;


-- name: ListUsers :many
SELECT
    id,
    email,
    username,
    user_role
FROM person
ORDER BY username;

-- name: UpdateRole :one
UPDATE person
SET user_role = $2
WHERE id = $1
RETURNING id, email, username, password_hash, user_role;
