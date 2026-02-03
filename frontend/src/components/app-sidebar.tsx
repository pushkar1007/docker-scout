import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"
import { Container, FolderOpenDot, LayoutDashboard, PackagePlus, SquareTerminal, Vault } from "lucide-react"
import Link from "next/link";

export function AppSidebar() {

  const tabs = [
    {
      name: "Dashboard",
      icon: <LayoutDashboard />,
      link: "/dashboard"
    },
    {
      name: "Containers",
      icon: <Container />,
      link: "/containers"
    },
    {
      name: "Images",
      icon: <Vault />,
      link: "/images"
    },
    {
      name: "Volumes",
      icon: <PackagePlus />,
      link: "/volumes"
    },
    {
      name: "Terminal",
      icon: <SquareTerminal />,
      link: "/terminal"
    }
  ];
  return (
    <Sidebar collapsible="none" className="bg-[#1D232F] h-[calc(100vh-72px)]">
      <SidebarContent className="items-center mt-4">
        {tabs.map((tab) => (
          <Link href={tab.link} key={tab.name}>
            <SidebarGroup className="gap-4 items-center hover:bg-blue-600 hover:cursor-pointer ">
              {tab.icon}
              <SidebarGroupContent>{tab.name}</SidebarGroupContent>
            </SidebarGroup>
          </Link>
        ))}
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton className="text-md">
              <FolderOpenDot /> Project Name
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  )
}
