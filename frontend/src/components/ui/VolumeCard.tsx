"use client";

import { DockerVolume } from '@/types/types';
import { useEffect, useState } from 'react';


const VolumeCard = () => {
  const [volume, setVolume] = useState<DockerVolume[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchVolumes() {
      try {
        const res = await fetch("http://10.172.201.84/docker_api_server/volumes");

        if (!res.ok) throw new Error("API failed");

        const data = await res.json();
        setVolume(data.volumes ?? []);
      } catch (e) {
        console.error(e);
        setVolume([]);
      } finally {
        setLoading(false);
      }
    }

    fetchVolumes();
  }, []);

  if (loading)
    return (
      <div className="flex text-muted-foreground h-full justify-center items-center">
        Inspecting containers…
      </div>
    );

  if (!volume.length)
    return (
      <div className="text-muted-foreground">
        No Docker volumes. Build something real.
      </div>
    );



  return (
    <div className="grid m-3 mt-6 gap-6 sm:grid-cols-2 lg:grid-cols-3">
      {volume.map((vol) => (
        <div
          key={vol.name}
          className={`
        rounded-xl border p-4 transition-all duration-200
        bg-slate-900/70 backdrop-blur
        hover:-translate-y-1 hover:shadow-xl
        ${vol.in_use
              ? "border-emerald-500/50 hover:shadow-emerald-500/20"
              : "border-rose-500/40 hover:shadow-rose-500/20"
            }
      `}
        >
          <div className="flex items-center justify-between">
            <h3 className="text-sm font-mono text-slate-200 break-all">
              {vol.name.slice(0, 12)}…
            </h3>

            <span
              className={`
            text-xs font-semibold px-2 py-1 rounded-md
            ${vol.in_use
                  ? "bg-emerald-500/20 text-emerald-400"
                  : "bg-rose-500/20 text-rose-400"
                }
          `}
            >
              {vol.in_use ? "IN USE" : "UNUSED"}
            </span>
          </div>

          <div className="mt-4 space-y-2 text-sm text-slate-400">
            <p>
              <span className="text-slate-500">Mount:</span>{" "}
              <code className="text-sky-400 text-xs break-all">
                {vol.mountpoint}
              </code>
            </p>

            <p>
              <span className="text-slate-500">Created:</span>{" "}
              {new Date(vol.created_at).toLocaleDateString()}
            </p>
          </div>
        </div>
      ))}
    </div>

  );
};

export default VolumeCard;
