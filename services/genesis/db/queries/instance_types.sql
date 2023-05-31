-- name: CreateInstanceType :exec
INSERT INTO instance_types (
    provider,
    name,
    cpu_count,
    cpu_architecture,
    memory_mib,
    gpu_kind,
    gpu_count
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: GetInstanceType :one
SELECT *
FROM instance_types
WHERE provider = $1 AND name = $2;
