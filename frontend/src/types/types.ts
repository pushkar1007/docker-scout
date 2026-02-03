export interface DockerContainer {
  Id: string;
  Names: string[];
  Image: string;
  ImageID: string;
  Command: string;
  Created: number;

  Ports: DockerPort[];

  Labels?: Record<string, string>;

  State: string;
  Status: string;

  HostConfig: {
    NetworkMode: string;
  };

  Health?: {
    Status: string;
    FailingStreak: number;
  };

  NetworkSettings: {
    Networks: Record<string, DockerNetwork>;
  };

  Mounts: DockerMount[];
}

export interface DockerPort {
  IP: string;
  PrivatePort: number;
  PublicPort?: number;
  Type: string;
}

export interface DockerNetwork {
  NetworkID: string;
  EndpointID: string;
  Gateway: string;
  IPAddress: string;
  MacAddress: string;
  IPPrefixLen: number;
}

export interface DockerMount {
  Type: string;
  Source: string;
  Destination: string;
  Mode: string;
  RW: boolean;
  Propagation: string;
}

export interface DockerImage {
  Containers: number;
  Created: number;
  Id: string;
  Labels?: Record<string, string>; // labels may exist or be empty
  ParentId: string;
  RepoDigests: string[];
  RepoTags: string[];
  SharedSize: number;
  Size: number;
}

export type SystemSummary = {
  active_containers: number;
  avg_cpu: string;
  avg_memory: string;
  avg_net_io: string;
  avg_disk_io: string;
};

export type VolumeSummary = {
  name: string;
  driver: string;
  mountpoint: string;
  scope: string;
  labels: Record<string, string>;
  created_at: string;
  in_use: boolean;
};

export type VolumeListResponse = {
  volumes: VolumeSummary[];
};

export type CreateVolumeRequest = {
  name: string;
  driver: string;
  labels: Record<string, string>;
  options: Record<string, string>;
};

export interface DockerVolumeLabels {
  [key: string]: string;
}

export interface DockerVolume {
  /** The unique hash/name of the volume */
  name: string;
  /** The driver used (e.g., 'local') */
  driver: string;
  /** The full path on the host machine where data is stored */
  mountpoint: string;
  /** Visibility of the volume (e.g., 'local') */
  scope: string;
  /** Metadata labels */
  labels: DockerVolumeLabels;
  /** ISO 8601 timestamp of creation */
  created_at: string;
  /** Whether the volume is currently attached to a container */
  in_use: boolean;
}
