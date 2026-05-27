import { ConnectSteamButton } from '@features/auth/ui/ConnectSteamButton';
import { getSession } from '@shared/lib/auth';
import { Button } from '@shared/ui/atoms/Button';
import { HomeFooter } from '@widgets/HomeFooter';
import { HomeHeader } from '@widgets/HomeHeader';
import { BrainCircuit, Gamepad2, Library, Sparkles, Zap } from 'lucide-react';
import { redirect } from 'next/navigation';

export default async function Home() {
  const session = await getSession();

  if (session) {
    redirect('/games');
  }

  return (
    <div className="relative min-h-screen flex flex-col font-(family-name:--font-inter) bg-background-base text-gray-200">
      <HomeHeader />

      <section className="relative flex flex-col items-center justify-center min-h-[90vh] p-8 overflow-hidden text-center">
        <div className="absolute inset-0 z-0 bg-[url('https://cdn.akamai.steamstatic.com/steam/apps/1091500/library_600x900_2x.jpg')] bg-cover bg-center opacity-10 filter blur-xl mix-blend-overlay"></div>
        <div className="absolute inset-0 bg-linear-to-b from-background-base/60 via-background-base/90 to-background-base z-0"></div>

        <div className="relative z-10 max-w-3xl flex flex-col items-center gap-8 animate-in fade-in slide-in-from-bottom-4 duration-1000 mt-16">
          <div className="flex gap-2 items-center bg-background-surface px-4 py-2 rounded-full border border-primary/20 shadow-[0_0_15px_rgba(102,192,244,0.1)]">
            <Sparkles size={16} className="text-primary" />
            <span className="text-sm font-semibold tracking-wide text-primary">
              Seu Gaming Backlog Inteligente
            </span>
          </div>

          <h1 className="text-5xl sm:text-7xl font-extrabold tracking-tight text-white leading-tight">
            Pare de navegar. <br />
            <span className="bg-clip-text text-transparent bg-linear-to-r from-primary to-[#2a475e]">
              Comece a jogar.
            </span>
          </h1>

          <p className="text-lg sm:text-xl text-gray-400 max-w-xl leading-relaxed">
            Acabe com a paralisia da escolha. Conecte sua biblioteca Steam e deixe a IA escolher o
            jogo perfeito para o seu tempo e humor.
          </p>

          <div className="flex flex-col sm:flex-row gap-4 mt-4 w-full justify-center">
            <ConnectSteamButton />
            <Button
              variant="outline"
              size="lg"
              className="w-full sm:w-auto text-lg h-14 px-8 border-gray-600 hover:bg-gray-800"
            >
              Ver Demonstração
            </Button>
          </div>
        </div>
      </section>

      <section id="como-funciona" className="relative z-10 py-24 px-8 bg-background-base">
        <div className="max-w-6xl mx-auto flex flex-col items-center gap-16">
          <div className="text-center gap-4 flex flex-col items-center">
            <h2 className="text-3xl sm:text-4xl font-bold text-white">Como Funciona o YNG?</h2>
            <p className="text-gray-400 max-w-2xl text-lg">
              Em três passos simples você descobre exatamente qual jogo da poeira virtual merece sua
              atenção hoje.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 w-full">
            <FeatureCard
              icon={<Library size={32} className="text-primary" />}
              title="1. Sincronize"
              description="Conecte sua conta Steam e importamos toda a sua biblioteca de jogos de forma instantânea e segura."
            />
            <FeatureCard
              icon={<BrainCircuit size={32} className="text-primary" />}
              title="2. Diga seu tempo disponível"
              description="Diga-nos quanto tempo você tem disponível e quais os gêneros ou humor você está buscando hoje."
            />
            <FeatureCard
              icon={<Gamepad2 size={32} className="text-primary" />}
              title="3. Obtenha o melhor jogo"
              description="Nossa IA analisa sua biblioteca e retorna a melhor recomendação para sua sessão. É só clicar e jogar."
            />
          </div>
        </div>
      </section>
      <section
        id="funcionalidades"
        className="relative z-10 py-24 px-8 bg-linear-to-b from-background-surface to-background-base"
      >
        <div className="max-w-5xl mx-auto grid md:grid-cols-2 gap-16 items-center">
          <div className="flex flex-col gap-6">
            <h2 className="text-3xl font-bold text-white leading-tight">
              Chega de perder 30 minutos escolhendo o que jogar.
            </h2>
            <p className="text-gray-400 text-lg">
              A Síndrome do PC Gamer afeta milhares de jogadores. Com tantas opções, nós acabamos
              não escolhendo nada.
            </p>

            <ul className="flex flex-col gap-4 mt-4">
              <CheckListItem text="Combate do cansaço mental (Paradox of Choice)." />
              <CheckListItem text="Aproveitamento real dos jogos que você já comprou." />
              <CheckListItem text="Filtros dinâmicos que se adaptam à sua rotina corrida." />
            </ul>
          </div>

          <div className="relative w-full max-w-[280px] aspect-2/3 mx-auto rounded-xl border border-[#2a475e] shadow-[0_0_30px_rgba(102,192,244,0.15)] overflow-hidden flex items-center justify-center group">
            <div className="absolute inset-0 bg-[url('https://cdn.akamai.steamstatic.com/steam/apps/1091500/library_600x900_2x.jpg')] bg-cover bg-center transition-transform duration-700 group-hover:scale-105"></div>
            <div className="absolute inset-0 bg-linear-to-t from-background-base via-transparent to-transparent opacity-80"></div>
            <div className="z-10 flex flex-col items-center gap-2 text-center p-6 mt-auto">
              <Zap size={32} className="text-primary" />
              <p className="text-sm font-bold text-white uppercase tracking-widest">
                Ação Instantânea
              </p>
            </div>
          </div>
        </div>
      </section>

      <HomeFooter />
    </div>
  );
}

function FeatureCard({
  icon,
  title,
  description,
}: {
  icon: React.ReactNode;
  title: string;
  description: string;
}) {
  return (
    <div className="flex flex-col gap-4 p-8 rounded-2xl bg-background-surface border border-primary/50 hover:border-primary/50 transition-colors shadow-lg">
      <div className="w-14 h-14 rounded-xl bg-background-base flex items-center justify-center border border-primary/20">
        {icon}
      </div>
      <h3 className="text-xl font-bold text-white mt-2">{title}</h3>
      <p className="text-gray-400 leading-relaxed">{description}</p>
    </div>
  );
}

function CheckListItem({ text }: { text: string }) {
  return (
    <li className="flex items-start gap-3">
      <div className="mt-1 w-5 h-5 rounded-full bg-primary/20 flex items-center justify-center">
        <div className="w-2 h-2 rounded-full bg-primary"></div>
      </div>
      <span className="text-gray-300 text-lg">{text}</span>
    </li>
  );
}
