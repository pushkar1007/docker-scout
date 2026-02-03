"use client";

import { DockerImage } from "@/types/types";
import { useEffect, useState } from "react";

export default function DockerImages() {
  const [images, setImages] = useState<DockerImage[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchImages() {
      const res = await fetch("http://10.172.201.84/docker_api_server/images");
      const data = await res.json();
      const images = data.Items;
      setImages(images);
      setLoading(false);
    }

    fetchImages();
  }, []);



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
          key={img.Id}
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
            <div className="font-medium truncate">{img.RepoTags}</div>
            <span className="text-xs text-muted-foreground">
              {formatBytes(img.Size)}
            </span>
          </div>

          <div className="mt-3 space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="">ID</span>
              <code className="text-xs">
                {img.Id.replace("sha256:", "").slice(0, 12)}
              </code>
            </div>

            <div className="flex justify-between">
              <span className="text-muted-foreground">Created</span>
              <span>
                {new Date(img.Created * 1000).toLocaleDateString()}
              </span>
            </div>
          </div>

          <div className="mt-3 flex flex-wrap gap-1">
            {img.Labels && Object.entries(img.Labels).map(([key, value]) => (
              <span
                key={key}
                className="rounded-md bg-muted px-2 py-0.5 text-xs"
              >
                {key}: {value}
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

