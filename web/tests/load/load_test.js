import { sleep } from 'k6';
import { APIClient, assertStatus } from './modules/client.js';
import { USER_ID, INVALID_ACCESS_TOKEN } from './modules/constants.js';

export const options = {
  stages: [
    { duration: '5s', target: 10 },  // Ramp up
    { duration: '15s', target: 15 }, // Stay at 15 VUs
    { duration: '5s', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_failed: ['rate<0.15'], // Account for expected Steam 502 failures (10%) + Chaos scenarios
    http_req_duration: ['p(95)<1500'],
  },
};

export default function () {
  const client = new APIClient();
  const unauthorizedClient = new APIClient(INVALID_ACCESS_TOKEN);
  const anonymousClient = new APIClient(null);

  // ==========================================
  // SCENARIO 1: Public / Health Check Endpoints
  // ==========================================
  assertStatus(client.get('/healthz'), 'Healthz public endpoint', 200);
  assertStatus(client.get('/readyz'), 'Readyz public endpoint', 200);
  assertStatus(client.get('/api/v1/version'), 'Version public endpoint', 200);
  assertStatus(client.get('/api/v1/auth/steam/login'), 'Steam login redirect', [200, 302]);

  // ==========================================
  // SCENARIO 2: Public Catalog (Outbound Calls)
  // ==========================================
  assertStatus(client.get('/api/v1/catalog/steam/apps/400'), 'Steam Catalog details', [200, 404, 502, 504]);

  // ==========================================
  // SCENARIO 3: Library Operations (Authorized)
  // ==========================================
  assertStatus(client.get(`/api/v1/users/${USER_ID}/games`), 'List user games', 200);
  assertStatus(client.get(`/api/v1/users/${USER_ID}/games/00000001-f75a-4b62-bb34-000000000001`), 'Get single game', [200, 404]);

  const newGamePayload = {
    id: `00000060-0060-0060-0060-${Math.floor(Math.random() * 1000000000000).toString().padStart(12, '0')}`,
    title: 'K6 Modular Chaos Game',
    coverUrl: 'https://example.com/cover.jpg',
    status: 'played',
    genre: 'RPG',
    platform: 'steam',
    description: 'Chaos Resilience Test Game',
    developer: 'Chaos Studio',
    releaseDate: '2026-05-26',
    hoursPlayed: 44.5,
  };
  assertStatus(client.post(`/api/v1/users/${USER_ID}/games`, newGamePayload), 'Create game', 201);

  // ==================================================================
  // SCENARIO 4: CHAOS ENGINEERING & FAULT INJECTION
  // ==================================================================

  // A. Fault Injection: Anonymous Access Attempt
  assertStatus(
    anonymousClient.get(`/api/v1/users/${USER_ID}/games`),
    'Fault Injection: Anonymous access to user library should be rejected',
    401
  );

  // B. Fault Injection: Invalid JWT Signature/Expired Token
  assertStatus(
    unauthorizedClient.get(`/api/v1/users/${USER_ID}/games`),
    'Fault Injection: Expired/Malformed JWT should be rejected',
    401
  );

  // C. Fault Injection: Corrupted Body Payload (Missing required columns)
  const badPayload = { id: '', title: '' }; // missing ID/Title should trigger validation block
  assertStatus(
    client.post(`/api/v1/users/${USER_ID}/games`, badPayload),
    'Fault Injection: Empty ID/Title game creation payload should be rejected',
    400
  );

  // D. Fault Injection: Malicious Input / SQL Injection attempt in Path
  // Parameterized queries resolve this safely as a normal string, yielding a clean 200 OK (empty array)
  assertStatus(
    client.get(`/api/v1/users/\' OR \'1\'=\'1/games`),
    'Fault Injection: SQL Injection attempt in path variable should be handled safely',
    [200, 400, 401, 404]
  );

  sleep(0.5); // Pace iterations
}
