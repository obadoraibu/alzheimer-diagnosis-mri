import { useNavigate } from 'react-router-dom';
import { BASE_URL }    from '../constants';

/**
 * useApi   — fetch-wrapper с авто-refresh-токена
 *
 * ▸ ставит Bearer accessToken из localStorage
 * ▸ при 401 вызывает POST /refresh { fingerprint }  (cookie ➜ refresh)
 * ▸ если пришёл новый токен → сохраняет и повторяет исходный запрос
 * ▸ если refresh неудачный → очищаем хранилище и редирект на /sign-in
 */
export const useApi = () => {
  const navigate = useNavigate();

  /* fingerprint, сохранённый в sign-in */
  const fp = () => navigator.userAgent + Math.random().toString(36).substring(2);

  /* выход с очисткой */
  const forceLogout = () => {
    localStorage.clear();
    navigate('/sign-in');
  };

  /* собственно fetch-обёртка */
  const api = async (url, opts = {}, retry = true) => {
    const res = await fetch(url.startsWith('http') ? url : `${BASE_URL}${url}`, {
      ...opts,
      credentials: 'include',                          // <── важен для refresh-cookie
      headers: {
        ...(opts.headers || {}),
        Authorization: `Bearer ${localStorage.getItem('accessToken') || ''}`,
      },
    });

    /* если всё нормально или уже была повторная попытка */
    if (res.status !== 401 || !retry) return res;

    /* ---------------- refresh ---------------- */
    try {
      const ref = await fetch(`${BASE_URL}/refresh`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ fingerprint: fp() }),
      });

      const body = await ref.json().catch(() => ({}));

      if (ref.ok && body.success && body.data?.accessToken) {
        localStorage.setItem('accessToken', body.data.accessToken);
        /* повторяем оригинальный запрос один раз                */
        return api(url, opts, false); // retry = false → чтобы не попасть в цикл
      }
    } catch (_) {
      /* ignore – перейдём к logout */
    }

    /* refresh не помог → выходим */
    forceLogout();
    return res;                       // отдаём исходный 401 (можно не использовать)
  };

  return api;
};
