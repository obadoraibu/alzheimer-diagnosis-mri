import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { formStyles as styles } from '../styles/styles';


function ResetRequest() {
  const [email,     setEmail]     = useState('');
  const [message,   setMessage]   = useState('');
  const [isSuccess, setIsSuccess] = useState(false);

  const handleSubmit = async e => {
    e.preventDefault();
    setMessage('');

    try {
      const res  = await fetch('http://localhost:8080/reset-password', {
        method : 'POST',
        headers: { 'Content-Type':'application/json' },
        body   : JSON.stringify({ email }),
      });

      const json = await res.json().catch(() => ({}));

      if (res.ok && json.success) {
        setIsSuccess(true);
        setMessage('Письмо с инструкциями отправлено на почту.');
        return;
      }

      if (!json.success) {
        const code = json.error?.code;

        if (code === 'USER_SUSPENDED') {
          setIsSuccess(false);
          setMessage('Ваш аккаунт заблокирован. Сброс пароля невозможен.');
          return;
        }

        setIsSuccess(false);
        setMessage(json.error?.message || json.message || 'Не удалось отправить письмо.');
        return;
      }

      setIsSuccess(false);
      setMessage(json.message || 'Ошибка запроса');

    } catch (err) {
      setIsSuccess(false);
      setMessage('Ошибка сети: ' + err.message);
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.formBox}>
        <h2 style={styles.heading}>Восстановление пароля</h2>

        <form onSubmit={handleSubmit}>
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

          <button style={styles.button} type="submit">
            Отправить
          </button>
        </form>

        {message && (
          <p
            style={{
              ...styles.message,
              color: isSuccess ? '#008080' : 'crimson',
              marginTop: 14,
            }}
          >
            {message}
          </p>
        )}

        <p style={styles.linkText}>
          Вспомнили пароль?{' '}
          <Link to="/sign-in" style={styles.link}>
            Войти
          </Link>
        </p>
      </div>
    </div>
  );
}

export default ResetRequest;
