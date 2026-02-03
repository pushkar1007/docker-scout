export type ContainerStats = {
  id: string;
  name: string;
  state: string;
  last_used: string;
  image: string;
  labels: Record<string, string>;
  ports: string;
  cpu: string;
  memory: string;
  net_io: string;
  disk_io: string;
};

export type DockerImage = {
  id: string;
  name: string;
  size: number;
  tags: string[];
  createdAt: string;
};

export type SystemSummary = {
  active_containers: number;
  avg_cpu: string;
  avg_memory: string;
  avg_net_io: string;
  avg_disk_io: string;
};

export type DashboardData = {
  summary: SystemSummary;
  containers: ContainerStats[];
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
