import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { formStyles as styles } from '../styles/styles';

function SignIn() {
  const [email,    setEmail]    = useState('');
  const [password, setPassword] = useState('');
  const [message,  setMessage]  = useState('');
  const [isError,  setIsError]  = useState(false);
  const navigate               = useNavigate();

  const fp = () => navigator.userAgent + Math.random().toString(36).substring(2);

  const handleSignIn = async e => {
    e.preventDefault();
    setMessage('');

    try {
      const res  = await fetch('http://localhost:8080/sign-in', {
        method :'POST',
        headers:{ 'Content-Type':'application/json' },
        body   : JSON.stringify({ email, password, fingerprint: fp() }),
      });
      console.log(res.code)

      if (res.ok) {
        const j   = await res.json();
        const tok = j?.data?.accessToken;
        if (j.success && tok) {
          localStorage.setItem('accessToken', tok);
          navigate('/home');
          return;
        }
        setIsError(true);
        setMessage('Некорректный ответ сервера.');
        return;
      }

      if (res.status === 401 || res.status === 403) {
        const err = await res.json().catch(() => ({}));
        const code = err.error?.code;

        if (code === 'WRONG_CREDENTIALS') {
          setIsError(true);
          setMessage('Неверный email или пароль');
          return;
        }

        if (code === 'USER_SUSPENDED') {
          setIsError(true);
          setMessage('Ваш аккаунт заблокирован. Обратитесь в администрацию.');
          return;
        }

        setIsError(true);
        setMessage('Ошибка авторизации');
        return;
      }

      setIsError(true);
      setMessage('Внутренняя ошибка. Попробуйте позже.');

    } catch (err) {
      setIsError(true);
      setMessage('Сетевая ошибка');
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.formBox}>
        <h2 style={styles.heading}>Авторизация</h2>

        <form onSubmit={handleSignIn}>
          <div style={styles.formGroup}>
            <label style={styles.label}>E-mail</label>
            <input
              style={styles.input}
              type="email"
              value={email}
              onChange={e => setEmail(e.target.value)}
              required
            />
          </div>

          <div style={styles.formGroup}>
            <label style={styles.label}>Пароль</label>
            <input
              style={styles.input}
              type="password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              required
            />
          </div>

          <button style={styles.button} type="submit">
            Войти
          </button>
        </form>

        {message && (
          <p style={{ ...styles.message, color: isError ? 'crimson' : styles.message.color }}>
            {message}
          </p>
        )}

        <p style={styles.linkText}>
          Забыли пароль?{' '}
          <Link to="/reset" style={styles.link}>
            Восстановить
          </Link>
        </p>
      </div>
    </div>
  );
}

export default SignIn;
