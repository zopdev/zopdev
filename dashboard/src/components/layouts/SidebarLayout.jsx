import { Outlet } from 'react-router-dom';
import Sidebar from '@/components/organisms/Sidebar';
import TopBar from '@/components/organisms/TopBar';

export default function SidebarLayout() {
  return (
    <div>
      <TopBar />
      <div className=" xs:inline-block md:flex h-[90vh] w-[100vw] overflow-hidden ">
        <div>
          <Sidebar
            menu={[
              {
                link: 'applications/123/environment',
                text: 'Environment',
                icon: <></>,
              },
              {
                link: 'applications/123/configDiff',
                text: 'Config Diff',
                icon: <></>,
              },
            ]}
          />
        </div>
        <main className="px-4 sm:px-6 lg:px-8 w-full overflow-auto text-left pt-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
