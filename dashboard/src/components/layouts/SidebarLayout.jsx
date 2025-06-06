import { Outlet, useParams } from 'react-router-dom';
import Sidebar from '@/components/organisms/Sidebar';
import TopBar from '@/components/organisms/TopBar';
import ResourceIcon from '@/assets/svg/sidebar/ResourceIcon';

export default function SidebarLayout() {
  const { cloudId } = useParams();
  return (
    <div className="flex flex-col h-screen">
      <div className="fixed top-0 w-full z-10">
        <TopBar />
      </div>
      <div className="flex flex-1 ">
        <div className="fixed left-0 h-full pt-16 hidden md:block">
          <Sidebar
            menu={[
              {
                link: `cloud-accounts/${cloudId}/resources`,
                text: 'Resources',
                icon: ResourceIcon,
              },
            ]}
          />
        </div>
        <div className="w-full md:ml-60 pt-16">
          <main className="px-4 sm:px-6 lg:px-8 w-full overflow-auto text-left pt-8 h-[calc(100vh-64px)]">
            <Outlet />
          </main>
        </div>
      </div>
    </div>
  );
}
