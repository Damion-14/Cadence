const TOKEN_KEY = 'cadence_token';
const USER_KEY = 'cadence_user';

export const setToken = (token: string): void => {
  localStorage.setItem(TOKEN_KEY, token);
};

export const getToken = (): string | null => {
  return localStorage.getItem(TOKEN_KEY);
};

export const clearToken = (): void => {
  localStorage.removeItem(TOKEN_KEY);
};

export const setUser = (user: unknown): void => {
  localStorage.setItem(USER_KEY, JSON.stringify(user));
};

export const getUser = (): unknown | null => {
  const userStr = localStorage.getItem(USER_KEY);
  if (!userStr) return null;
  try {
    return JSON.parse(userStr);
  } catch {
    return null;
  }
};

export const clearUser = (): void => {
  localStorage.removeItem(USER_KEY);
};

export const clearAuth = (): void => {
  clearToken();
  clearUser();
};

export const isAuthenticated = (): boolean => {
  return !!getToken();
};
