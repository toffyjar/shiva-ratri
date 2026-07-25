import HttpService from '../../services/httpService';

function mockFetchOnce(options: {
  ok?: boolean;
  status?: number;
  json?: unknown;
  text?: string;
}) {
  const { ok = true, status = 200, json, text } = options;
  const response = {
    ok,
    status,
    json: jest.fn().mockResolvedValue(json),
    text: jest.fn().mockResolvedValue(text ?? ''),
  };
  (global.fetch as jest.Mock).mockResolvedValueOnce(response);
  return response;
}

describe('HttpService', () => {
  beforeEach(() => {
    (global as any).fetch = jest.fn();
  });

  describe('constructor base URL', () => {
    it('uses the provided base URL', async () => {
      mockFetchOnce({ json: {} });
      const svc = new HttpService('https://custom.example.com');
      await svc.get('/ping');

      const calledUrl = (global.fetch as jest.Mock).mock.calls[0][0];
      expect(calledUrl).toBe('https://custom.example.com/ping');
    });

    it('falls back to localhost:9393 when EXPO_PUBLIC_JYOTISH_URL is unset', async () => {
      const original = process.env.EXPO_PUBLIC_JYOTISH_URL;
      delete process.env.EXPO_PUBLIC_JYOTISH_URL;

      mockFetchOnce({ json: {} });
      const svc = new HttpService();
      await svc.get('/ping');

      const calledUrl = (global.fetch as jest.Mock).mock.calls[0][0];
      expect(calledUrl).toBe('http://localhost:9393/ping');

      if (original !== undefined) process.env.EXPO_PUBLIC_JYOTISH_URL = original;
    });
  });

  describe('get', () => {
    it('appends query params to the URL', async () => {
      mockFetchOnce({ json: { hello: 'world' } });
      const svc = new HttpService('https://api.example.com');

      await svc.get('/search', { q: 'delhi', limit: '5' });

      const calledUrl = (global.fetch as jest.Mock).mock.calls[0][0];
      expect(calledUrl).toBe('https://api.example.com/search?q=delhi&limit=5');
    });

    it('omits the query string when no params are given', async () => {
      mockFetchOnce({ json: {} });
      const svc = new HttpService('https://api.example.com');

      await svc.get('/search');

      const calledUrl = (global.fetch as jest.Mock).mock.calls[0][0];
      expect(calledUrl).toBe('https://api.example.com/search');
    });

    it('sends a GET request with a JSON content-type header', async () => {
      mockFetchOnce({ json: {} });
      const svc = new HttpService('https://api.example.com');

      await svc.get('/ping');

      const options = (global.fetch as jest.Mock).mock.calls[0][1];
      expect(options.method).toBe('GET');
      expect(options.headers['Content-Type']).toBe('application/json');
    });

    it('merges custom headers on top of the default content-type', async () => {
      mockFetchOnce({ json: {} });
      const svc = new HttpService('https://api.example.com');

      await svc.get('/ping', undefined, { Authorization: 'Bearer abc' });

      const options = (global.fetch as jest.Mock).mock.calls[0][1];
      expect(options.headers['Authorization']).toBe('Bearer abc');
      expect(options.headers['Content-Type']).toBe('application/json');
    });

    it('returns ok, status, and parsed JSON data', async () => {
      mockFetchOnce({ ok: true, status: 200, json: { id: 1 } });
      const svc = new HttpService('https://api.example.com');

      const result = await svc.get<{ id: number }>('/items/1');

      expect(result).toEqual({ ok: true, status: 200, data: { id: 1 } });
    });

    it('surfaces non-2xx responses without throwing', async () => {
      mockFetchOnce({ ok: false, status: 404, json: { error: 'not_found' } });
      const svc = new HttpService('https://api.example.com');

      const result = await svc.get('/missing');

      expect(result.ok).toBe(false);
      expect(result.status).toBe(404);
      expect(result.data).toEqual({ error: 'not_found' });
    });
  });

  describe('post', () => {
    it('sends a POST request with a JSON-stringified body', async () => {
      mockFetchOnce({ ok: true, status: 201, json: { id: 2 } });
      const svc = new HttpService('https://api.example.com');

      await svc.post('/items', { name: 'test' });

      const [url, options] = (global.fetch as jest.Mock).mock.calls[0];
      expect(url).toBe('https://api.example.com/items');
      expect(options.method).toBe('POST');
      expect(options.body).toBe(JSON.stringify({ name: 'test' }));
      expect(options.headers['Content-Type']).toBe('application/json');
    });

    it('sends undefined body when none is provided', async () => {
      mockFetchOnce({ ok: true, status: 200, json: {} });
      const svc = new HttpService('https://api.example.com');

      await svc.post('/ping');

      const options = (global.fetch as jest.Mock).mock.calls[0][1];
      expect(options.body).toBeUndefined();
    });

    it('returns parsed JSON data from the response', async () => {
      mockFetchOnce({ ok: true, status: 201, json: { token: 'abc123' } });
      const svc = new HttpService('https://api.example.com');

      const result = await svc.post<{ token: string }>('/login', { email: 'a@b.com' });

      expect(result).toEqual({ ok: true, status: 201, data: { token: 'abc123' } });
    });
  });

  describe('getSvg', () => {
    it('requests an SVG content-type and returns raw text data', async () => {
      const svgMarkup = '<svg></svg>';
      mockFetchOnce({ ok: true, status: 200, text: svgMarkup });
      const svc = new HttpService('https://api.example.com');

      const result = await svc.getSvg('/chart/svg', { size: '100' });

      const [url, options] = (global.fetch as jest.Mock).mock.calls[0];
      expect(url).toBe('https://api.example.com/chart/svg?size=100');
      expect(options.headers['Content-Type']).toBe('image/svg+xml');
      expect(result).toEqual({ ok: true, status: 200, data: svgMarkup });
    });
  });
});
