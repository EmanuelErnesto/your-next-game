export function HomeFooter() {
  return (
    <footer className="w-full bg-[#171a21] py-8 border-t border-[#2a475e]/30 mt-auto z-10 relative">
      <div className="max-w-6xl mx-auto px-8 flex flex-col md:flex-row items-center justify-between gap-4">
        <div className="flex flex-col gap-1 text-center md:text-left">
          <span className="text-white font-bold tracking-widest text-lg">YNG</span>
          <span className="text-xs text-gray-500">
            © 2026 Your Next Game. Todos os direitos reservados.
          </span>
        </div>
        <div className="flex gap-6 text-sm text-gray-400">
          <a href="#" className="hover:text-primary transition-colors">
            Política de Privacidade
          </a>
          <a href="#" className="hover:text-primary transition-colors">
            Termos de Serviço
          </a>
          <a href="#" className="hover:text-primary transition-colors">
            Contato
          </a>
        </div>
      </div>
    </footer>
  );
}
