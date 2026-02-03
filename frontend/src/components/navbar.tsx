import Image from "next/image";

function Navbar() {
    return(
        <header className="w-screen h-fit bg-[#1D232F]">
            <nav className="p-4 px-12">
                <div className="flex gap-2 items-center">
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
        </header>
    );
}

export default Navbar;