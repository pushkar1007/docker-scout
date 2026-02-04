"use client"

import BasicPie from "@/components/piechart";
import CompositionExample from "@/components/speed-gauge";
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import axios from "axios";
import { ChevronDownIcon, ListRestart, Pause, Play, StopCircle, Trash2 } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import type { ContainerData, DashboardData } from "@/types/types";

const WEBSOCKET =
  process.env.NEXT_PUBLIC_WEBSOCKET || "ws://localhost:3000/stats";


const Page = () => {

  const wsRef = useRef<WebSocket | null>(null);
  useEffect(() => {

    const ws = new WebSocket(WEBSOCKET);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log('Connected to server');
      ws.send(JSON.stringify({ type: 'greet', payload: 'Hello Server!' }));
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data as string);
        console.log('Received:', data);
        setStats(data);
      } catch (err) {
        console.error('Invalid message format', err);
      }
    };

    ws.onclose = () => console.log('Connection closed');
    ws.onerror = (err) => console.error('WebSocket error', err);
  }, []);



  const [stats, setStats] = useState<DashboardData | null>(null);
  const [containers, setContainers] = useState<ContainerData | null>(null);
  const [pausedContainers, setPausedContainers] = useState<number>(0);
  const [stoppedContainers, setStoppedContainers] = useState<number>(0);
  useEffect(() => {
    const fetchData = async () => {
      const [containersResponse] = await Promise.all([
        axios.get<ContainerData>("http://10.172.201.84/docker_api_server/containers"),
      ]);
      const items = containersResponse.data?.Items ?? [];
      const paused = items.filter((container) => container.State === "paused").length;
      const stopped = items.filter(
        (container) => container.State === "exited" || container.State === "created"
      ).length;

      setContainers(containersResponse.data);
      // setStats(statsResponse.data);
      setPausedContainers(paused);
      setStoppedContainers(stopped);
    };
    fetchData();
  }, []);

  const avgCpu = Number(stats?.summary?.avg_cpu ?? 0);
  const avgCpuPercent = Number.isFinite(avgCpu) ? avgCpu * 100 : 0;

  console.log(containers);
  console.log(stats);

  return (
    <div className="w-full flex flex-col gap-4">
      <div className="w-full space-y-2 font-bold">
        <h1 className="text-3xl">Dashboard</h1>
        <Separator />
      </div>
      <div className="w-full flex gap-4 justify-between">
        <div className="border border-muted-foreground/20 bg-[#1D232F] w-1/4 h-[25vh] rounded-2xl flex flex-col justify-between p-4">
          <h2 className="text-xl font-semibold">Total Containers</h2>
          <div className="flex justify-center items-center">
            <h1 className="text-6xl font-bold">{containers?.Items?.length ?? 0}</h1>
          </div>
          <div className="flex justify-between">
            <div className="flex gap-1">
              <h3 className="text-md text-muted-foreground">Running: </h3>
              <h3 className="text-md text-green-500">{stats?.summary?.active_containers ?? 0}</h3>
            </div>
            <Separator orientation="vertical" className="bg-muted-foreground/20" />
            <div className="flex gap-1">
              <h3 className="text-md text-muted-foreground">Paused: </h3>
              <h3 className="text-md text-yellow-500">{pausedContainers}</h3>
            </div>
            <Separator orientation="vertical" className="bg-muted-foreground/20" />
            <div className="flex gap-1">
              <h3 className="text-md text-muted-foreground">Stopped: </h3>
              <h3 className="text-md text-red-500">{stoppedContainers}</h3>
            </div>
          </div>
        </div>
        <div className="border border-muted-foreground/20 bg-[#1D232F] w-1/4 h-[25vh] rounded-2xl flex flex-col justify-between p-4">
          <h2 className="text-xl font-semibold">Group Count</h2>
          <div className="flex justify-center items-center">
            <h1 className="text-6xl font-bold">4</h1>
          </div>
          <div className="flex justify-center">
            <h3 className="text-md text-muted-foreground">Active Groups: 4</h3>
          </div>
        </div>
        <div className="border border-muted-foreground/20 bg-[#1D232F] w-1/4 h-[25vh] rounded-2xl flex flex-col justify-between p-4">
          <h2 className="text-xl font-semibold">Orphaned Containers</h2>
          <div className="flex justify-center items-center">
            <h1 className="text-6xl font-bold">3</h1>
          </div>
          <div className="flex justify-center">
            <h3 className="text-md text-blue-500 underline hover:cursor-pointer">View Unlabeled Lists</h3>
          </div>
        </div>
        <div className="border border-muted-foreground/20 bg-[#1D232F] w-1/4 h-[25vh] rounded-2xl flex flex-col justify-between p-4">
          <h2 className="text-xl font-semibold">Total Resource Load</h2>
          <div className="flex justify-center items-center">
            <CompositionExample value={avgCpuPercent} />
          </div>
          <div className="flex justify-center">
            <div className="flex gap-3">
              <h3 className="text-md text-muted-foreground">CPU: {stats?.summary?.avg_cpu ?? "0%"}</h3>
              <Separator orientation="vertical" className="bg-muted-foreground/20" />
              <h3 className="text-md text-muted-foreground">RAM: {stats?.summary?.avg_memory ?? "0 GB"}</h3>
            </div>
          </div>
        </div>
      </div>
      <div className="w-full flex gap-4 justify-between">
        <div className="border border-muted-foreground/20 bg-[#1D232F] w-3/5 h-[25vh] max-h-[25vh] rounded-2xl flex flex-col p-4">
          <h2 className="text-xl font-semibold">Resource Hog Leaderboard</h2>
          <div className="mt-2">
            <div className="flex">
              <h3 className="w-3/5 font-semibold">Top Container</h3>
              <h3 className="w-1/5 font-semibold">RAM</h3>
              <h3 className="w-1/5 font-semibold">CPU</h3>
            </div>
            <div className="flex">
              <h3 className="w-3/5">1. postgres-db</h3>
              <h3 className="w-1/5">1.8G RAM</h3>
              <h3 className="w-1/5">25% CPU</h3>
            </div>
            <div className="flex">
              <h3 className="w-3/5">2. analytics-worker</h3>
              <h3 className="w-1/5">1.2G RAM</h3>
              <h3 className="w-1/5">18% CPU</h3>
            </div>
            <div className="flex">
              <h3 className="w-3/5">3. cache-server</h3>
              <h3 className="w-1/5">900M RAM</h3>
              <h3 className="w-1/5">15% CPU</h3>
            </div>
          </div>
        </div>
        <div className="border border-muted-foreground/20 bg-[#1D232F] w-2/5 h-[25vh] rounded-2xl flex flex-col justify-between p-4">
          <h2 className="text-xl font-semibold">Group Distribution</h2>
          <div className="flex justify-center"><BasicPie /></div>
        </div>
      </div>
      <div className="border border-muted-foreground/20 bg-[#1D232F] w-full h-fit min-h-[25vh] rounded-2xl p-4">
        <h2 className="text-xl font-semibold">Container Explorer &#40;Overview&#41;</h2>
        <Accordion
          type="single"
          collapsible
          defaultValue="shipping"
          className=""
        >
          <AccordionItem value="shipping">
            <AccordionTrigger className="flex">
              <div className="flex items-center gap-4">
                <ChevronDownIcon className="text-muted-foreground pointer-events-none size-4 shrink-0 translate-y-0.5 transition-transform duration-200" />
                <h3 className="text-lg">Project: Auth-Service &#40;Running: 3, 1 Stopped | 1.2G RAM&#41;</h3>
              </div>
              <div className="flex gap-1">
                <Button variant="outline" size="sm" className="cursor-pointer"><Play className="text-green-500" /></Button>
                <Button variant="outline" size="sm" className="cursor-pointer"><Pause className="text-yellow-500" /></Button>
                <Button variant="outline" size="sm" className="cursor-pointer"><StopCircle className="text-red-500" /></Button>
              </div>
            </AccordionTrigger>
            <AccordionContent className="flex flex-col">
              {stats?.containers?.map((stat) => (
                <div key={stat.Id} className="flex items-center justify-between">
                  <div className="flex gap-4 items-center w-1/5">
                    <div className="rounded-full w-4 h-4 bg-chart-2 animate-pulse"></div>
                    <h3 className="text-md truncate leading-none mt-[-4px]">{stat.name ?? "Unknown"}</h3>
                  </div>
                  <p className="text-md mt-[-4px] w-1/5 truncate">{stat.image ?? "N/A"}</p>
                  <p className="text-md mt-[-4px] w-1/5 truncate">{stat.memory}</p>
                  <div className="flex w-1/5">
                    <div className="flex items-center">
                      <ListRestart />
                      <Button variant="ghost" className="cursor-pointer">Restart</Button>
                    </div>
                    <div className="flex items-center text-red-500">
                      <Trash2 />
                      <Button variant="ghost" className="cursor-pointer hover:text-red-500">Delete</Button>
                    </div>
                  </div>
                </div>
              ))}
            </AccordionContent>
          </AccordionItem>
        </Accordion>
      </div>
    </div>
  );
}

export default Page;
