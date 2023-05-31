SET search_path TO genesis;

CREATE TYPE PROVIDER AS ENUM ( -- noqa: L014
    'amazon-web-services',
    'google-cloud-platform'
);

CREATE TYPE CPU_ARCHITECTURE AS ENUM ( -- noqa: L014
    'x86-64',
    'arm64'
);

CREATE TYPE GPU_KIND AS ENUM ( -- noqa: L014
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

CREATE TYPE INSTANCE_OWNER AS ENUM ( -- noqa: L014
    'scheduler'
);
