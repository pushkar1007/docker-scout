import { Separator } from "@/components/ui/separator";
import XTerminal from "@/components/xterminal";


const Page = () => {
  return (
    <div className="overflow-none w-full h-full">
      <div className="overflow-hidden w-full h-full space-y-2 font-semibold">
        <h1 className="text-2xl">Terminal</h1>
        <Separator />
        <XTerminal />
      </div>
      <div>

      </div>
      <div></div>
    </div>
  );
}

export default Page;
