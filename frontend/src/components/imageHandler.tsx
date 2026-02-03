"use client";

// export const mockDockerImages: DockerImage[] = [
//   {
//     id: "sha256:9f3a2c7b4d1e8a22f0a123456789abcd",
//     name: "metrics-api",
//     size: 245_678_912, // ~234 MB
//     tags: ["latest", "prod"],
//     createdAt: "2026-02-01T10:32:00Z",
//   },
//   {
//     id: "sha256:ab12cd34ef56a7890bcdef1234567890",
//     name: "metrics-worker",
//     size: 198_432_102,
//     tags: ["latest"],
//     createdAt: "2026-01-30T18:05:00Z",
//   },
//   {
//     id: "sha256:deadbeefcafebabefeedface12345678",
//     name: "frontend-dashboard",
//     size: 156_901_888,
//     tags: ["dev", "latest"],
//     createdAt: "2026-01-29T14:11:00Z",
//   },
//   {
//     id: "sha256:1234abcd5678efgh9012ijklmnopqrst",
//     name: "postgres-backup",
//     size: 512_344_990,
//     tags: ["nightly"],
//     createdAt: "2026-01-28T02:44:00Z",
//   },
//   {
//     id: "sha256:c0ffee1234567890badf00dabcdef9876",
//     name: "redis-cache",
//     size: 63_204_321,
//     tags: ["stable"],
//     createdAt: "2026-01-27T09:20:00Z",
//   },
// ];
//

import { DockerImage } from "@/types/types";
import { useEffect, useState } from "react";

export default function DockerImages() {
  const [images, setImages] = useState<DockerImage[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchImages() {
      const res = await fetch("/api/docker/images");
      const data = await res.json();
      setImages(data);
      setLoading(false);
    }

    fetchImages();
  }, []);

  // useEffect(() => {
  //   setLoading(false)
  //   setImages(mockDockerImages)
  // }, [])


  if (loading)
    return (
      <div className="flex text-muted-foreground h-full justify-center items-center">
        Inspecting containers…
      </div>
    );

  if (!images.length)
    return (
      <div className="text-muted-foreground">
        No Docker images. Build something real.
      </div>
    );

  return (
    <div className="grid m-3 mt-6 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {images.map((img) => (
        <div
          key={img.id}
          className="
    rounded-xl
    border border-border
    backdrop-blur
    p-4
    bg-blue-900       shadow-sm
    transition
    hover:bg-blue-600/80   "
        >

          <div className="flex items-start justify-between">
            <div className="font-medium truncate">{img.name}</div>
            <span className="text-xs text-muted-foreground">
              {formatBytes(img.size)}
            </span>
          </div>

          <div className="mt-3 space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="">ID</span>
              <code className="text-xs">
                {img.id.replace("sha256:", "").slice(0, 12)}
              </code>
            </div>

            <div className="flex justify-between">
              <span className="text-muted-foreground">Created</span>
              <span>
                {new Date(img.createdAt).toLocaleDateString()}
              </span>
            </div>
          </div>

          <div className="mt-3 flex flex-wrap gap-1">
            {img.tags.map((tag) => (
              <span
                key={tag}
                className="rounded-md bg-muted px-2 py-0.5 text-xs"
              >
                {tag}
              </span>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}

function formatBytes(bytes: number) {
  const units = ["B", "KB", "MB", "GB"];
  let i = 0;

  while (bytes >= 1024 && i < units.length - 1) {
    bytes /= 1024;
    i++;
  }

  return `${bytes.toFixed(1)} ${units[i]}`;
}

