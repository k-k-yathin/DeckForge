import { Outlet } from 'react-router-dom';
import { Sidebar } from './Sidebar';

/** App shell with sidebar — used for all authenticated pages */
export function AppLayout() {
  return (
    <div className="flex min-h-screen">
      <Sidebar />
      <main className="flex-1 overflow-auto bg-slate-50 p-8">
        <Outlet />
      </main>
    </div>
  );
}
