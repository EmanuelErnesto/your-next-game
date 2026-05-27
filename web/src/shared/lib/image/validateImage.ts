import { logger } from '../logger';

export async function validateImageUrl(url: string): Promise<boolean> {
  if (!url) return false;

  if (url.startsWith('/') || url.includes('placeholder')) return true;

  try {
    const res = await fetch(url, {
      method: 'HEAD',
      next: { revalidate: 3600 * 24 },
    });

    if (res.ok) {
      logger.info('Image validated [200 OK]:', { url });
    } else {
      logger.warn('Image failed validation [404+]:', { url, status: res.status, method: 'HEAD' });
    }

    return res.ok;
  } catch (e) {
    logger.error('Falha na validação de imagem (Erro de Rede):', e as Error, { url });
    return false;
  }
}
