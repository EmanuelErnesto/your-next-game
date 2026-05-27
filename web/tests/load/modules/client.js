import http from 'k6/http';
import { check } from 'k6';
import { BASE_URL, ACCESS_TOKEN } from './constants.js';

export class APIClient {
  constructor(token = ACCESS_TOKEN) {
    this.token = token;
  }

  getHeaders(customHeaders) {
    const headers = Object.assign({}, customHeaders || {});
    if (this.token) {
      headers['Cookie'] = `access_token=${this.token}`;
    }
    return { headers };
  }

  get(path, customHeaders) {
    const url = `${BASE_URL}${path}`;
    const params = this.getHeaders(customHeaders);
    return http.get(url, params);
  }

  post(path, body, customHeaders) {
    const url = `${BASE_URL}${path}`;
    const headers = Object.assign({ 'Content-Type': 'application/json' }, customHeaders || {});
    const params = this.getHeaders(headers);
    const payload = typeof body === 'string' ? body : JSON.stringify(body);
    return http.post(url, payload, params);
  }
}

// Global assertion helper to keep scenarios clean
export function assertStatus(response, description, expectedStatuses) {
  const statuses = Array.isArray(expectedStatuses) ? expectedStatuses : [expectedStatuses];
  const passed = statuses.includes(response.status);
  
  check(response, {
    [`[${response.status}] ${description}`]: () => passed,
  });

  return passed;
}
