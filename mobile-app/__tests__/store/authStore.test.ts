jest.mock('expo-secure-store', () => ({
  getItemAsync: jest.fn(),
  setItemAsync: jest.fn().mockResolvedValue(undefined),
  deleteItemAsync: jest.fn().mockResolvedValue(undefined),
}));

jest.mock('../../services/httpService', () => ({
  authHttpService: {
    post: jest.fn(),
  },
}));

import * as SecureStore from 'expo-secure-store';
import { useAuthStore } from '../../store/authStore';
import { authHttpService } from '../../services/httpService';

const mockedPost = authHttpService.post as jest.Mock;
const mockedGetItemAsync = SecureStore.getItemAsync as jest.Mock;
const mockedSetItemAsync = SecureStore.setItemAsync as jest.Mock;
const mockedDeleteItemAsync = SecureStore.deleteItemAsync as jest.Mock;

const resetStore = () => {
  useAuthStore.setState({
    token: null,
    email: null,
    isAuthenticated: false,
    isHydrating: true,
    isSubmitting: false,
    error: null,
  });
};

describe('useAuthStore', () => {
  beforeEach(() => {
    resetStore();
  });

  describe('hydrate', () => {
    it('marks the user authenticated when a token is stored', async () => {
      mockedGetItemAsync.mockImplementation((key: string) =>
        Promise.resolve(key.includes('email') ? 'user@example.com' : 'stored-token')
      );

      await useAuthStore.getState().hydrate();

      const state = useAuthStore.getState();
      expect(state.token).toBe('stored-token');
      expect(state.email).toBe('user@example.com');
      expect(state.isAuthenticated).toBe(true);
      expect(state.isHydrating).toBe(false);
    });

    it('leaves the user unauthenticated when nothing is stored', async () => {
      mockedGetItemAsync.mockResolvedValue(null);

      await useAuthStore.getState().hydrate();

      const state = useAuthStore.getState();
      expect(state.token).toBeNull();
      expect(state.isAuthenticated).toBe(false);
      expect(state.isHydrating).toBe(false);
    });
  });

  describe('login', () => {
    it('sets authenticated state and persists the session on success', async () => {
      mockedPost.mockResolvedValue({
        ok: true,
        status: 200,
        data: { id: 1, email: 'user@example.com', token: 'jwt-token' },
      });

      const result = await useAuthStore.getState().login('user@example.com', 'password123');

      expect(result).toBe(true);
      expect(mockedPost).toHaveBeenCalledWith('/login', {
        email: 'user@example.com',
        password: 'password123',
      });

      const state = useAuthStore.getState();
      expect(state.token).toBe('jwt-token');
      expect(state.email).toBe('user@example.com');
      expect(state.isAuthenticated).toBe(true);
      expect(state.isSubmitting).toBe(false);
      expect(state.error).toBeNull();

      const persistedValues = mockedSetItemAsync.mock.calls.map((call) => call[1]);
      expect(persistedValues).toEqual(expect.arrayContaining(['jwt-token', 'user@example.com']));
    });

    it('sets the server-provided error message on failure', async () => {
      mockedPost.mockResolvedValue({
        ok: false,
        status: 401,
        data: { error: 'invalid_credentials', message: 'Invalid email or password' },
      });

      const result = await useAuthStore.getState().login('user@example.com', 'wrongpass');

      expect(result).toBe(false);
      const state = useAuthStore.getState();
      expect(state.isAuthenticated).toBe(false);
      expect(state.error).toBe('Invalid email or password');
      expect(mockedSetItemAsync).not.toHaveBeenCalled();
    });

    it('falls back to a default error message when the server sends none', async () => {
      mockedPost.mockResolvedValue({ ok: false, status: 500, data: {} });

      await useAuthStore.getState().login('user@example.com', 'password123');

      expect(useAuthStore.getState().error).toBe('Invalid email or password');
    });

    it('sets a network error message when the request throws', async () => {
      mockedPost.mockRejectedValue(new Error('network down'));

      const result = await useAuthStore.getState().login('user@example.com', 'password123');

      expect(result).toBe(false);
      const state = useAuthStore.getState();
      expect(state.isAuthenticated).toBe(false);
      expect(state.isSubmitting).toBe(false);
      expect(state.error).toBe('Unable to reach server. Please try again.');
    });
  });

  describe('register', () => {
    it('sets authenticated state and persists the session on success', async () => {
      mockedPost.mockResolvedValue({
        ok: true,
        status: 201,
        data: { id: 2, email: 'new@example.com', token: 'new-jwt-token' },
      });

      const result = await useAuthStore.getState().register('new@example.com', 'password123');

      expect(result).toBe(true);
      expect(mockedPost).toHaveBeenCalledWith('/register', {
        email: 'new@example.com',
        password: 'password123',
      });

      const state = useAuthStore.getState();
      expect(state.token).toBe('new-jwt-token');
      expect(state.isAuthenticated).toBe(true);
    });

    it('sets the server-provided error message on failure', async () => {
      mockedPost.mockResolvedValue({
        ok: false,
        status: 409,
        data: { error: 'email_exists', message: 'Email already registered' },
      });

      const result = await useAuthStore.getState().register('taken@example.com', 'password123');

      expect(result).toBe(false);
      expect(useAuthStore.getState().error).toBe('Email already registered');
    });

    it('sets a network error message when the request throws', async () => {
      mockedPost.mockRejectedValue(new Error('network down'));

      const result = await useAuthStore.getState().register('new@example.com', 'password123');

      expect(result).toBe(false);
      expect(useAuthStore.getState().error).toBe('Unable to reach server. Please try again.');
    });
  });

  describe('logout', () => {
    it('clears the session and removes stored credentials', async () => {
      useAuthStore.setState({ token: 'jwt-token', email: 'user@example.com', isAuthenticated: true });

      await useAuthStore.getState().logout();

      expect(mockedDeleteItemAsync).toHaveBeenCalledTimes(2);
      const state = useAuthStore.getState();
      expect(state.token).toBeNull();
      expect(state.email).toBeNull();
      expect(state.isAuthenticated).toBe(false);
    });
  });

  describe('clearError', () => {
    it('resets the error field to null', () => {
      useAuthStore.setState({ error: 'some error' });

      useAuthStore.getState().clearError();

      expect(useAuthStore.getState().error).toBeNull();
    });
  });
});
