import DockerImages from "@/components/imageHandler";
import { Separator } from "@/components/ui/separator";

const Page = () => {
  return (
    <div className="w-full h-full overflow-y-hidden">
      <div className="w-full space-y-2 font-semibold">
        <h1 className="text-2xl">Images</h1>
      </div>
      <Separator />
      <DockerImages />
    </div>
  );
}

export default Page;

