import { useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { BASE_URL } from '../constants';

/**
 * Хук-обёртка вокруг fetch:
 *  ▸ автоматически ставит Bearer-токен;
 *  ▸ при 401 пытается обновить access-token;
 *  ▸ при провале refresh — редиректит на /sign-in.
 */
export function useApi() {
  const navigate = useNavigate();

  const logout = useCallback(() => {
    localStorage.clear();
    navigate('/sign-in');
  }, [navigate]);

  /* попытка обновить пары токенов */
  const refreshTokens = useCallback(async () => {
    const refresh = localStorage.getItem('refreshToken');
    if (!refresh) return false;

    const res = await fetch(`${BASE_URL}/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body:   JSON.stringify({ refresh }),
      credentials: 'include',
    });

    if (!res.ok) return false;
    const data = await res.json();
    data.access  && localStorage.setItem('accessToken',  data.access);
    data.refresh && localStorage.setItem('refreshToken', data.refresh);
    return true;
  }, []);

  /* сам fetch-wrapper */
  const api = useCallback(
    async (url, opts = {}, retry = true) => {
      const res = await fetch(
        url.startsWith('http') ? url : `${BASE_URL}${url}`,
        {
          ...opts,
          headers: {
            ...(opts.headers || {}),
            Authorization: `Bearer ${localStorage.getItem('accessToken')}`,
          },
        }
      );

      if (res.status !== 401 || !retry) return res;
      if (!(await refreshTokens())) {
        logout();
        return res;
      }
      return api(url, opts, false); // повторяем один раз
    },
    [logout, refreshTokens]
  );

  return api;
}
