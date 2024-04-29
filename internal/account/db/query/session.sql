-- name: CreateSession :one
INSERT INTO sessions (
  id ,
  email,
  refresh_token,
  user_agent,
  client_ip,
  expires_at
) VALUES (
    $1 , $2, $3, $4, $5 , $6
) RETURNING * ;

-- name: GetSessionById :one
SELECT * FROM sessions WHERE id = $1;