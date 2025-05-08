import { Outlet } from 'react-router-dom';
import TopBar from '@/components/organisms/TopBar.jsx';

export default function MainLayout() {
  return (
    <div className="flex flex-col h-screen">
      <div className="sticky top-0 z-50 bg-white shadow">
        <TopBar />
      </div>
      <div className="flex-1 overflow-y-auto">
        <Outlet />
      </div>
    </div>
  );
}
