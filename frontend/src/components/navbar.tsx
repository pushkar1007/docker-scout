import Image from "next/image";
import Link from 'next/link'; ``

function Navbar() {
  return (
    <header className="w-screen h-fit bg-[#1D232F]">
      <Link href="/dashboard" className="block w-screen">
        <nav className="p-4 px-12">
          <div className="flex gap-2 items-center" >
            <Image
              className="dark:invert"
              src="/docker.svg"
              alt="Docker logomark"
              width={36}
              height={36}
            />
            <h1>Docker Scout</h1>
          </div>
          <div></div>
        </nav>
      </Link>
    </header>
  );
}

export default Navbar;
