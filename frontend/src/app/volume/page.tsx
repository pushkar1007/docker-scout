import { Separator } from "@/components/ui/separator";
import VolumeCard from "@/components/ui/VolumeCard";

const Page = () => {
  return (
    <div className="w-full">
      <div className="w-full space-y-2 font-semibold">
        <h1 className="text-2xl">Volumes</h1>
      </div>
      <div>
        <Separator />
        <VolumeCard />
      </div>
    </div>
  );
}

export default Page;


