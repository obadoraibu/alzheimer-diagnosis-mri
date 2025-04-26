import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { signInStyles as styles } from '../styles/styles';

function SignIn() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const navigate = useNavigate();

 
  const generateFingerprint = () => {
    return navigator.userAgent + Math.random().toString(36).substring(2);
  };

  const handleSignIn = async (e) => {
    e.preventDefault();

    const fingerprint = generateFingerprint();

    try {
      const response = await fetch('http://localhost:8080/sign-in', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password, fingerprint }),
      });

      if (response.ok) {
        const data = await response.json();
        if (data.access && data.refresh) {
          localStorage.setItem('accessToken', data.access);
          localStorage.setItem('refreshToken', data.refresh);
          navigate('/home');
        } else {
          setMessage('Invalid response from server.');
        }
      } else {
        const errorData = await response.json();
        setMessage('Error: ' + (errorData.message || 'Something went wrong'));
      }
    } catch (err) {
      setMessage('Network error: ' + err.message);
    }
  };


  


  return (
    <div style={styles.container}>
      <div style={styles.loginBox}>
        <h2 style={styles.heading}>Авторизация</h2>
        <form onSubmit={handleSignIn}>
          <div style={styles.formGroup}>
            <label style={styles.label}>E-mail</label>
            <input
              style={styles.input}
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>

          <div style={styles.formGroup}>
            <label style={styles.label}>Пароль</label>
            <input
              style={styles.input}
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>

          <button style={styles.button} type="submit">
            Войти
          </button>
        </form>
        {message && <p style={styles.message}>{message}</p>}
        <p style={styles.registerText}>
          Забыли пароль?{' '}
          <Link to="/sign-up" style={styles.registerLink}>
            Восстановить
          </Link>
        </p>
      </div>
    </div>
  );
}

export default SignIn;
