import DockerImages from "@/components/imageHandler";
import { Separator } from "@/components/ui/separator";

const Page = () => {
  return (
    <div className="w-full h-full">
      <div className="w-full space-y-2 font-bold">
        <h1 className="text-3xl">Images</h1>
        <Separator />
      </div>
      <DockerImages />
    </div>
  );
}

export default Page;

