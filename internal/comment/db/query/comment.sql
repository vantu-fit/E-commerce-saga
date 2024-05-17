-- name: CreateComment :one
INSERT INTO comments (
  product_id, 
  user_id, 
  content, 
  left_index,
  right_index,
  parent_id
) 
VALUES ($1, $2, $3, $4, $5 , $6) RETURNING *;

-- name: ListComment :many
SELECT * FROM comments 
WHERE product_id = $1 
AND left_index < $2 
AND right_index > $3 
ORDER BY left_index ASC;

-- name: UpdateContentComment :one
UPDATE comments
SET content = $1
WHERE id = $2
RETURNING *;

-- name: UpdateLeftIndexComment :many
UPDATE comments
SET left_index = left_index + $3
WHERE (parent_id = $1 or id = $1) and left_index > $2
RETURNING *;

-- name: UpdateRightIndexComment :many
UPDATE comments
SET right_index = right_index + $3
WHERE (parent_id = $1 or id = $1) and right_index >= $2
RETURNING *;

-- name: DeleteComment :many
DELETE FROM comments
WHERE (parent_id = $1 or id = $1 )and left_index >= $2 and right_index <= $3
RETURNING *;

-- name: DeleteCommentById :one
DELETE FROM comments
WHERE id = $1
RETURNING *;

-- name: GetCommentForUpdate :one
SELECT * FROM comments
WHERE id = $1
FOR NO KEY UPDATE;


-- name: GetMaxRightIndex :one
SELECT MAX(right_index) as max_right_index 
FROM comments
WHERE parent_id = $1;

-- name: GetCommentByProductID :many
SELECT * FROM comments
WHERE product_id = $1
ORDER BY left_index ASC;

-- name: GetCommentByID :one
SELECT * FROM comments
WHERE id = $1;

-- name: GetAllComments :many
select * from comments 
where (parent_id = $1 or id = $1) and left_index >= $2 and right_index <= $3 ;


