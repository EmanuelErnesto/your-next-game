import 'server-only';

import { getMe } from '@shared/backend';

export async function getSession() {
  return getMe();
}
