import React, { useState } from 'react';
import { homeStyles as styles } from '../styles/styles';

const Row = ({ label, children }) => (
  <>
    <label
      style={{
        ...styles.label,
        display: 'block',
        margin: '14px 0 6px 2px',
        fontWeight: 700,
      }}
    >
      {label}
    </label>
    {children}
  </>
);

export default function InviteUserModal({ api, close, onSuccess }) {
  const [name,  setName]  = useState('');
  const [email, setEmail] = useState('');
  const [role,  setRole]  = useState('doctor');
  const [busy,  setBusy]  = useState(false);

  const submit = async e => {
    e.preventDefault();
    setBusy(true);
    const r = await api('/admin/users', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: name, email, role }),
    });
    setBusy(false);
    if (r.ok) {
      onSuccess && onSuccess();
      close();
      alert('Приглашение отправлено');
    } else {
      alert('Ошибка');
    }
  };

  return (
    <div style={styles.modalOverlay} onClick={close}>
      <div
        onClick={e => e.stopPropagation()}
        style={{ ...styles.modalContent, padding: 0, maxWidth: 540 }}
      >
        {/* ── Header ── */}
        <div
          style={{
            background: '#2f6c71',
            color: '#fff',
            fontSize: 30,
            fontWeight: 700,
            textAlign: 'center',
            padding: '18px 16px',
            position: 'relative',
          }}
        >
          Пригласить&nbsp;пользователя
          <span
            onClick={close}
            style={{ position: 'absolute', right: 22, top: 6, fontSize: 36, cursor: 'pointer' }}
          >
            ×
          </span>
        </div>

        {/* ── Form ── */}
        <form onSubmit={submit} style={{ padding: 38 }}>
          <Row label="Имя">
            <input
              style={{ ...styles.input, width: '100%' }}
              value={name}
              onChange={e => setName(e.target.value)}
              required
            />
          </Row>

          <Row label="Email">
            <input
              style={{ ...styles.input, width: '100%' }}
              type="email"
              value={email}
              onChange={e => setEmail(e.target.value)}
              required
            />
          </Row>

          <Row label="Роль">
            <select
              style={{ ...styles.input, width: '100%' }}
              value={role}
              onChange={e => setRole(e.target.value)}
            >
              <option value="doctor">Врач</option>
              <option value="admin">Админ</option>
            </select>
          </Row>

          <div style={{ textAlign: 'center', marginTop: 26 }}>
            <button
              type="submit"
              disabled={busy}
              style={{
                ...styles.uploadButton,
                fontFamily: 'inherit',
                padding: '12px 70px',
                opacity: busy ? 0.6 : 1,
              }}
            >
              Отправить
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
