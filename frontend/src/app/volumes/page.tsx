import { Separator } from "@/components/ui/separator";
import VolumeCard from "@/components/ui/VolumeCard";

const Page = () => {
  return (
    <div className="w-full h-full">
      <div className="w-full space-y-2 font-bold">
        <h1 className="text-3xl">Volumes</h1>
        <Separator />
      </div>
      <VolumeCard />
    </div>
  );
}

export default Page;


