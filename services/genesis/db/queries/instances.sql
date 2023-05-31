-- name: CreateInstance :exec
INSERT INTO instances (
    id,
    provider,
    provider_id,
    zone,
    owner,
    instance_type,
    is_spot,
    cpu_count_requested,
    memory_mib_requested,
    memory_mib_reserved,
    gpu_count_requested,
    boot_image,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
);

-- name: SetInstanceBooting :exec
UPDATE instances
SET
    booted_at = sqlc.arg('booted_at')::TIMESTAMPTZ,
    hostname = sqlc.arg('hostname')::TEXT
WHERE id = $1;

-- name: SetInstanceRunning :exec
UPDATE instances
SET started_at = sqlc.arg('started_at')::TIMESTAMPTZ
WHERE id = $1;

-- name: SetInstanceDeleted :exec
UPDATE instances
SET deleted_at = sqlc.arg('deleted_at')::TIMESTAMPTZ
WHERE id = $1;

-- name: SetInstanceDeletionTriaged :exec
UPDATE instances
SET is_deletion_triaged = TRUE
WHERE id = $1;

-- name: GetInstance :one
SELECT *
FROM instances
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListRunningInstances :many
SELECT *
FROM instances
WHERE deleted_at IS NULL
    AND started_at IS NOT NULL
    AND owner = $1;

-- name: ListPastRequestedInstances :many
SELECT *
FROM instances
WHERE deleted_at IS NULL
    AND booted_at IS NULL
    AND created_at < NOW() - sqlc.arg('min_age')::INTERVAL;

-- name: ListPastBootedInstances :many
SELECT *
FROM instances
WHERE deleted_at IS NULL
    AND booted_at IS NOT NULL
    AND started_at IS NULL
    AND booted_at < NOW() - sqlc.arg('min_age')::INTERVAL;

-- name: ListPastDeletedUntriagedInstances :many
SELECT *
FROM instances
WHERE deleted_at IS NOT NULL
    AND deleted_at < NOW() - sqlc.arg('min_age')::INTERVAL
    AND is_deletion_triaged = FALSE;
