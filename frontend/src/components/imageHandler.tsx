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
    <div className="grid m-3 mt-6 gap-6 sm:grid-cols-2 lg:grid-cols-3">
      {images.map((img) => (
        <div
          key={img.Id}
          className="
        rounded-xl border border-slate-500/70
        p-4 transition-all duration-200
        bg-gray-800/60 backdrop-blur
        hover:-translate-y-0.5
        hover:shadow-lg hover:shadow-black/40
      "
        >
          {/* Header */}
          <div className="flex items-center justify-between">
            <h3 className="text-sm font-mono text-slate-200 truncate">
              {img.RepoTags?.[0] ?? "untagged"}
            </h3>

            <span className="text-xs px-2 py-0.5 rounded-md bg-slate-800 text-slate-300">
              {formatBytes(img.Size)}
            </span>
          </div>

          {/* Body */}
          <div className="mt-4 space-y-2 text-sm text-slate-400">
            <p className="flex justify-between">
              <span className="text-slate-500">ID</span>
              <code className="text-sky-400 text-xs">
                {img.Id.replace("sha256:", "").slice(0, 12)}
              </code>
            </p>

            <p className="flex justify-between">
              <span className="text-slate-500">Created</span>
              <span>
                {new Date(img.Created * 1000).toLocaleDateString()}
              </span>
            </p>
          </div>

          {/* Labels */}
          {img.Labels && Object.keys(img.Labels).length > 0 && (
            <div className="mt-3 flex flex-wrap gap-1">
              {Object.entries(img.Labels).map(([key, value]) => (
                <span
                  key={key}
                  className="rounded-md bg-slate-800/70 px-2 py-0.5 text-xs text-slate-400"
                >
                  {key}: {value}
                </span>
              ))}
            </div>
          )}
        </div>
      ))}
    </div>
  )
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

