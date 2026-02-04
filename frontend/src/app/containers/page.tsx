"use client"

import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import axios from "axios";
import { Pause, Play, StopCircle } from "lucide-react";
import { JSXElementConstructor, Key, ReactElement, ReactNode, ReactPortal, useEffect, useState } from "react";

const Page = () => {
    const [data, setData] = useState<any>(null);
    const baseUrl = process.env.NEXT_PUBLIC_BASE_URL;
    
    const fetchData = async () => {
        const response = await axios.get(`${baseUrl}/containers`);
        setData(response.data);
    };
    
    useEffect(() => {
        fetchData();
    }, []);

    const handleStart = async (containerId: string) => {
        try {
            await axios.post(`${baseUrl}/containers/start?id=${containerId}`);
            await fetchData(); // Refresh the container list
        } catch (error) {
            console.error('Error starting container:', error);
        }
    };

    const handlePause = async (containerId: string) => {
        try {
            await axios.post(`${baseUrl}/containers/pause?id=${containerId}`);
            await fetchData(); // Refresh the container list
        } catch (error) {
            console.error('Error pausing container:', error);
        }
    };

    const handleStop = async (containerId: string) => {
        try {
            await axios.post(`${baseUrl}/containers/stop?id=${containerId}`);
            await fetchData(); // Refresh the container list
        } catch (error) {
            console.error('Error stopping container:', error);
        }
    };

    console.log(data);
    return(
        <div className="w-full">
            <div className="w-full space-y-2 font-bold">
                <h1 className="text-3xl">Containers</h1>
                <Separator />
            </div>
            <div className="mt-4 flex flex-col gap-4">
                {data?.Items.map((container: { Id: Key | null | undefined; Names: (string | number | bigint | boolean | ReactElement<unknown, string | JSXElementConstructor<any>> | Iterable<ReactNode> | ReactPortal | Promise<string | number | bigint | boolean | ReactPortal | ReactElement<unknown, string | JSXElementConstructor<any>> | Iterable<ReactNode> | null | undefined> | null | undefined)[]; }) => (
                    <div key={container.Id} className="rounded-xl min-h-[20vh] h-fit bg-[#1D232F] p-6">
                        <div className="flex justify-between">
                            <div className="flex items-center gap-2">
                                <div className={`w-4 h-4 ${container.State === 'running' ? (container.State === 'exited' || container.State === 'paused' ? 'bg-chart-3' : 'bg-chart-2') : 'bg-chart-5'} animate-pulse rounded-full`}></div>
                                <h2 className="text-xl font-semibold mt-2 leading-none mt-[-5px]">{container.Names[0]}</h2>
                            </div>
                            <div className="flex gap-1">
                                <Button variant="outline" size="sm" className="cursor-pointer" onClick={() => handleStart(container.Id)}><Play className="text-green-500" /></Button>
                                <Button variant="outline" size="sm" className="cursor-pointer" onClick={() => handlePause(container.Id)}><Pause className="text-yellow-500"/></Button>
                                <Button variant="outline" size="sm" className="cursor-pointer" onClick={() => handleStop(container.Id)}><StopCircle className="text-red-500" /></Button>
                            </div>
                        </div>
                        <div>
                            <p className="text-md text-muted-foreground mt-4">Image Name: {container.Image}</p>
                            <p className="text-md text-muted-foreground mt-4">Status: {container.Status}</p>
                            <p className="text-md text-muted-foreground mt-4">Mounts</p>
                            <div className="flex flex-col gap-2 mt-2">
                                {container.Mounts.map((mount) => (
                                    <div className="bg-black w-fit">
                                        {mount.Source}
                                    </div>
                                ))}
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}

export default Page;