SET search_path TO genesis;

CREATE TABLE instance_types (
    -- The provider declaring the instance type.
    provider         PROVIDER         NOT NULL,
    -- The provider-specific name of the instance type.
    name             TEXT             NOT NULL,
    -- The number of (logical) CPUs available on the instance.
    cpu_count        INTEGER          NOT NULL,
    -- The CPU architecture of the instance.
    cpu_architecture CPU_ARCHITECTURE NOT NULL,
    -- The actual amount of memory in MiB available on the instance.
    memory_mib       INTEGER          NOT NULL,
    -- The kind of GPU that is available on the instance (if any).
    gpu_kind         GPU_KIND         NULL,
    -- The number of GPUs that are connected to this instance (if any).
    gpu_count        INTEGER          NULL,

    PRIMARY KEY (provider, name),
    CHECK (name != ''),
    CHECK (cpu_count >= 1),
    CHECK (memory_mib >= 1024),
    CHECK ((gpu_kind IS NULL) = (gpu_count IS NULL)),
    CHECK (gpu_count IS NULL OR gpu_count >= 1)
);
