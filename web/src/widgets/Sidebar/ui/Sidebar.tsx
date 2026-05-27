import { useAuth } from '@shared/auth';
import { Checkbox } from '@shared/ui/atoms/Checkbox';
import Image from 'next/image';

export function Sidebar() {
  const { session, isLoading } = useAuth();
  const userName = session?.name || 'Convidado';
  const userImage = session?.image || 'https://i.pravatar.cc/150?img=11';

  return (
    <aside className="w-64 border-r border-[#2C2940] bg-background-base h-full p-4 flex flex-col gap-6 overflow-y-auto">
      <div className="flex items-center justify-between hover:bg-background-surface p-2 rounded-lg cursor-pointer transition text-gray-200">
        <div className="flex items-center gap-3">
          <div className="avatar">
            <div className="w-10 h-10 relative rounded-full overflow-hidden border border-primary/20">
              <Image src={userImage} alt={userName} fill className="object-cover" />
            </div>
          </div>
          <span className="font-semibold text-sm truncate max-w-[120px]">
            {isLoading ? 'Carregando...' : userName}
          </span>
        </div>
        <span className="text-gray-400">⌄</span>
      </div>

      <div className="divider m-0 opacity-20"></div>

      <div className="flex flex-col gap-3">
        <h3 className="font-semibold text-sm text-gray-300">Filtros</h3>
        <Checkbox label="Gênero" />
        <Checkbox label="Plataforma" />
        <Checkbox label="Status" />
      </div>
    </aside>
  );
}
