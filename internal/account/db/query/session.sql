-- name: CreateSession :one
INSERT INTO sessions (
  id ,
  user_id,
  refresh_token,
  user_agent,
  client_ip
) VALUES (
    $1 , $2, $3, $4, $5 
) RETURNING * ;

-- name: GetSessionById :one
SELECT * FROM sessions WHERE id = $1;
