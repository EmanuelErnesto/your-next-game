import { GlobalWindow } from 'happy-dom';

const window = new GlobalWindow({ url: 'http://localhost:3001' });
globalThis.window = window as any;
globalThis.document = window.document as any;
globalThis.navigator = window.navigator as any;
globalThis.HTMLElement = window.HTMLElement as any;
globalThis.customElements = window.customElements as any;
