import { Separator } from "@/components/ui/separator";
import XTerminal from "@/components/xterminal";


const Page = () => {
  return (
    <div className="overflow-none w-full h-full">
      <div className="overflow-hidden w-full h-full space-y-2 font-semibold">
        <div className="w-full space-y-2 font-bold">
          <h1 className="text-3xl">Terminal</h1>
          <Separator />
        </div>
        <XTerminal />
      </div>
      <div>

      </div>
      <div></div>
    </div>
  );
}

export default Page;
