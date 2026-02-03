import BasicPie from "@/components/piechart";
import CompositionExample from "@/components/speed-gauge";
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { ChevronDownIcon, Delete, ListRestart, Pause, Play, StopCircle, Trash, Trash2 } from "lucide-react";

const Page = () => {
    return(
        <div className="w-full flex flex-col gap-4">
            <div className="w-full space-y-2 font-bold">
                <h1 className="text-3xl">Dashboard</h1>
                <Separator />
            </div>
            <div className="w-full flex gap-4 justify-between">
                <div className="border border-muted-foreground/20 bg-[#1D232F] w-1/4 h-[25vh] rounded-2xl flex flex-col justify-between p-4">
                    <h2 className="text-xl font-semibold">Total Containers</h2>
                    <div className="flex justify-center items-center">
                        <h1 className="text-6xl font-bold">12</h1>
                    </div>
                    <div className="flex justify-between">
                        <div className="flex gap-1">
                            <h3 className="text-md text-muted-foreground">Running: </h3>
                            <h3 className="text-md text-green-500">4</h3>
                        </div>
                        <Separator orientation="vertical" className="bg-muted-foreground/20"/>
                        <div className="flex gap-1">
                            <h3 className="text-md text-muted-foreground">Paused: </h3>
                            <h3 className="text-md text-yellow-500">2</h3>
                        </div>
                        <Separator orientation="vertical" className="bg-muted-foreground/20"/>
                        <div className="flex gap-1">
                            <h3 className="text-md text-muted-foreground">Stopped: </h3>
                            <h3 className="text-md text-red-500">6</h3>
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
                    <h2 className="text-xl font-semibold">Orphaned Conatiners</h2>
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
                        <CompositionExample />
                    </div>
                    <div className="flex justify-center">
                        <div className="flex gap-3">
                            <h3 className="text-md text-muted-foreground">CPU: 38%</h3>
                            <Separator orientation="vertical" className="bg-muted-foreground/20"/>
                            <h3 className="text-md text-muted-foreground">RAM: 6 GB</h3>
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
                                <Button variant="outline" size="sm" className="cursor-pointer"><Pause className="text-yellow-500"/></Button>
                                <Button variant="outline" size="sm" className="cursor-pointer"><StopCircle className="text-red-500" /></Button>
                            </div>
                        </AccordionTrigger>
                        <AccordionContent className="flex items-center justify-between">
                            <div className="rounded-full w-4 h-4 bg-chart-2 animate-pulse"></div>
                            <h3>name</h3>
                            <p>image</p>
                            <p>ram usage</p>
                            <div className="flex">
                                <div className="flex items-center">
                                    <ListRestart />
                                    <Button variant="ghost" className="cursor-pointer">Restart</Button>
                                </div>
                                <div className="flex items-center text-red-500">
                                    <Trash2 />
                                    <Button variant="ghost" className="cursor-pointer hover:text-red-500">Delete</Button>
                                </div>
                            </div>
                        </AccordionContent>
                    </AccordionItem>
                </Accordion>    
            </div>
        </div>
    );
}

export default Page;