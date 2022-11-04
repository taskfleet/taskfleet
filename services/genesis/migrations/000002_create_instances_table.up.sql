-----------------------------------------------------------------------------------------------------------------------
-- TYPES
-----------------------------------------------------------------------------------------------------------------------

CREATE TYPE genesis.PROVIDER AS ENUM ( -- noqa: L014
    'amazon-web-services',
    'google-cloud-platform'
);
CREATE TYPE genesis.CPU_ARCHITECTURE AS ENUM ( -- noqa: L014
    'x86-64',
    'arm64'
);
CREATE TYPE genesis.GPU_KIND AS ENUM ( -- noqa: L014
    'nvidia-tesla-k80', -- Standard version with 12 GB memory
    'nvidia-tesla-m60', -- Standard version with 8 GB memory
    'nvidia-tesla-p100', -- Standard version with 16 GB memory
    'nvidia-tesla-p4', -- Standard version with 8 GB memory
    'nvidia-tesla-v100', -- Standard version with 16 GB memory
    'nvidia-tesla-t4', -- Standard version with 16 GB memory
    'nvidia-tesla-a100', -- Standard version with 40 GB memory
    'nvidia-tesla-a10', -- Standard version with 24 GB memory
    'nvidia-tesla-a100-80gb' -- Extended version with 80 GB memory
);
CREATE TYPE genesis.INSTANCE_OWNER AS ENUM ( -- noqa: L014
    'scheduler'
);

-----------------------------------------------------------------------------------------------------------------------
-- TABLE INSTANCE TYPES
-----------------------------------------------------------------------------------------------------------------------

CREATE TABLE genesis.instance_types (
    -- The provider declaring the instance type.
    provider genesis.PROVIDER NOT NULL,
    -- The provider-specific name of the instance type.
    name TEXT NOT NULL,
    -- The number of (logical) CPUs available on the instance.
    cpu_count INTEGER NOT NULL,
    -- The CPU architecture of the instance.
    cpu_architecture genesis.CPU_ARCHITECTURE NOT NULL,
    -- The actual amount of memory in MiB available on the instance.
    memory_mib INTEGER NOT NULL,
    -- The kind of GPU that is available on the instance (if any).
    gpu_kind genesis.GPU_KIND NULL,
    -- The number of GPUs that are connected to this instance (if any).
    gpu_count INTEGER NULL,

    PRIMARY KEY(provider, name),
    CHECK((gpu_kind IS NULL) = (gpu_count IS NULL))
);

-----------------------------------------------------------------------------------------------------------------------
-- TABLE INSTANCES
-----------------------------------------------------------------------------------------------------------------------

CREATE TABLE genesis.instances (
    /*
    *  IDENTIFYING INFORMATION
    */
    -- The globally unique ID of an instance. Unique forever and across providers.
    id UUID PRIMARY KEY,
    -- The provider the instance is registered with.
    provider genesis.PROVIDER NOT NULL,
    -- The unique ID of the instance on a specific provider. Only guaranteed to be unique
    -- while the instance is running.
    provider_id TEXT NOT NULL,
    -- The provider-specific availability zone where the instance was created.
    zone TEXT NOT NULL,
    -- The external owner whom this instance belongs to.
    owner genesis.INSTANCE_OWNER NOT NULL,
    /*
    *  INSTANCE CONFIGURATION
    */
    -- The provider-specific unique identifier for the type of instance launched.
    instance_type TEXT NOT NULL,
    -- An indicator whether the instance is a spot instance.
    is_spot BOOLEAN NOT NULL,
    -- The number of (logical) CPUs requested by the owner upon instance creation.
    cpu_count_requested INTEGER NOT NULL,
    -- The amount of memory in MiB requested by the owner upon instance creation.
    memory_mib_requested INTEGER NOT NULL,
    -- The amount of memory reserved by the component that the instance was launched for.
    memory_mib_reserved INTEGER NOT NULL,
    -- The number of GPUs requested by the owner upon instance creation.
    gpu_count_requested INTEGER NOT NULL,
    -- The unique identifier of the machine image/AMI that the instance was launched with.
    boot_image TEXT NOT NULL,
    /*
    *  INSTANCE STATUS
    */
    -- The hostname of the instance while it was running.
    hostname TEXT,
    -- The timestamp at which the instance was created.
    created_at TIMESTAMPTZ NOT NULL,
    -- The timestamp at which the provider declared the instance running.
    booted_at TIMESTAMPTZ,
    -- The timestamp at which Genesis declared the instance running, i.e. the timestamp
    -- at which the component running on the instance became responsive.
    started_at TIMESTAMPTZ,
    -- The timestamp at which the instance was deleted with the provider.
    deleted_at TIMESTAMPTZ,
    -- This flag is ONLY set by the garbage collector to signal that an instance has been
    -- succesfully deleted and that it does not need to be collected at any time in the future.
    is_deletion_triaged BOOLEAN NOT NULL DEFAULT FALSE,

    FOREIGN KEY (provider, instance_type) REFERENCES genesis.instance_types (provider, name)
);

CREATE INDEX idx_instances_deleted_at ON genesis.instances(deleted_at);
CREATE INDEX idx_instances_is_deletion_triaged ON genesis.instances(is_deletion_triaged);
