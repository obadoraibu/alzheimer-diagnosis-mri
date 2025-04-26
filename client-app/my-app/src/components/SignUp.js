import React, { useState } from 'react';
import { Link, useParams, useNavigate } from 'react-router-dom';

function SignUp() {
  const { code } = useParams(); // получаем код из URL
  const navigate = useNavigate(); // навигатор для перехода
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

  const styles = {
    container: {
      minHeight: '100vh',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      backgroundColor: '#f2f2f2',
      padding: '1em',
    },
    loginBox: {
      backgroundColor: '#fff',
      borderRadius: '0.5em',
      boxShadow: '0 2px 6px rgba(0,0,0,0.1)',
      padding: '2em',
      maxWidth: '400px',
      width: '100%',
      margin: '1em',
    },
    heading: {
      marginBottom: '1em',
      textAlign: 'center',
      color: '#004D4D',
    },
    formGroup: {
      marginBottom: '1em',
    },
    label: {
      display: 'block',
      marginBottom: '0.5em',
      fontWeight: 'bold',
      color: '#333',
    },
    input: {
      width: '100%',
      padding: '0.75em',
      border: '1px solid #ccc',
      borderRadius: '0.25em',
      fontSize: '1em',
      boxSizing: 'border-box',
    },
    button: {
      width: '100%',
      padding: '0.75em',
      backgroundColor: '#008080',
      color: '#fff',
      border: 'none',
      borderRadius: '0.25em',
      fontSize: '1em',
      cursor: 'pointer',
      marginTop: '1em',
    },
    message: {
      marginTop: '1em',
      textAlign: 'center',
    },
    loginText: {
      marginTop: '1em',
      textAlign: 'center',
      fontSize: '0.9em',
    },
    loginLink: {
      color: '#008080',
      textDecoration: 'none',
      fontWeight: 'bold',
    },
  };

  return (
    <div style={styles.container}>
      <div style={styles.loginBox}>
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

        <p style={styles.loginText}>
          Уже зарегистрированы?{' '}
          <Link to="/sign-in" style={styles.loginLink}>
            Войти
          </Link>
        </p>
      </div>
    </div>
  );
}

export default SignUp;
