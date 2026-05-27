import { getSession } from '@shared/lib/auth';
import { Sidebar } from '@widgets/Sidebar';
import { TopNavbar } from '@widgets/TopNavbar';
import { redirect } from 'next/navigation';

export default async function DashboardLayout({ children }: { children: React.ReactNode }) {
  const session = await getSession();

  if (!session) {
    redirect('/');
  }

  return (
    <div className="flex flex-col h-screen overflow-hidden bg-background-base">
      <TopNavbar />
      <div className="flex flex-1 overflow-hidden relative">
        <Sidebar />
        <main className="flex-1 bg-background-surface relative overflow-y-auto">{children}</main>
      </div>
    </div>
  );
}
