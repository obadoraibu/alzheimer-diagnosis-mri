import React, { useState } from 'react';
import { Link, useParams, useNavigate } from 'react-router-dom';
import { formStyles as styles } from '../styles/styles';

function SignUp() {
  const { code } = useParams(); 
  const navigate = useNavigate(); 
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const [isSuccess, setIsSuccess] = useState(false);

  console.log('Код из URL:', code);

  const handleSignUp = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(`http://localhost:8080/complete-invite/${code}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ password }),
      });

      if (response.ok) {
        setMessage('Регистрация завершена! Сейчас вы будете перенаправлены на вход...');
        setIsSuccess(true);

        setTimeout(() => {
          navigate('/sign-in');
        }, 2000);
      } else {
        const errorData = await response.json();
        setMessage('Ошибка: ' + (errorData.message || 'Что-то пошло не так'));
        setIsSuccess(false);
      }
    } catch (err) {
      setMessage('Ошибка сети: ' + err.message);
      setIsSuccess(false);
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.formBox}>
        <h2 style={styles.heading}>Регистрация по приглашению</h2>
        <form onSubmit={handleSignUp}>
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
            Зарегистрироваться
          </button>
        </form>

        {message && (
          <p
            style={{
              ...styles.message,
              color: isSuccess ? '#008080' : 'red',
            }}
          >
            {message}
          </p>
        )}

        <p style={styles.linkText}>
          Уже зарегистрированы?{' '}
          <Link to="/sign-in" style={styles.link}>
            Войти
          </Link>
        </p>
      </div>
    </div>
  );
}

export default SignUp;
