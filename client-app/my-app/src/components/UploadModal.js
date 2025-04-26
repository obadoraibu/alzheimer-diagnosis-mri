import React, { useState } from 'react';
import { homeStyles as styles } from '../styles/styles';

/* упрощённый label-row */
const Row = ({ label, children }) => (
  <div style={{ marginBottom: 20 }}>
    <label style={{ ...styles.label, display: 'block', marginBottom: 6 }}>{label}</label>
    {children}
  </div>
);

export default function UploadModal({ api, close, onSuccess }) {
  /* локальный state */
  const [name,     setName]     = useState('');
  const [gender,   setGender]   = useState('Male');
  const [age,      setAge]      = useState('');
  const [date,     setDate]     = useState('');
  const [file,     setFile]     = useState(null);
  const [loading,  setLoading]  = useState(false);

  /* отправка */
  const handleSubmit = async e => {
    e.preventDefault();
    if (!file) return;

    const f = new FormData();
    f.append('patient_name',   name);
    f.append('patient_gender', gender);
    f.append('patient_age',    age);
    f.append('scan_date',      date);
    f.append('file',           file);

    setLoading(true);
    const r = await api('/upload', { method: 'POST', body: f });
    setLoading(false);

    if (r.ok) {
      alert('Скан успешно загружен');
      onSuccess && onSuccess();
      close();
    } else {
      alert('Ошибка загрузки');
    }
  };

  /* кастомный radio-button */
  const Radio = ({ value, children }) => (
    <label
      style={{
        display: 'flex',
        alignItems: 'center',
        gap: 8,
        padding: '10px 24px',
        border: '1px solid #ccc',
        cursor: 'pointer',
        background: gender === value ? '#2f6c71' : '#fff',
        color: gender === value ? '#fff' : '#555',
        fontWeight: 500,
      }}
    >
      <input
        type="radio"
        value={value}
        checked={gender === value}
        onChange={() => setGender(value)}
        style={{ display: 'none' }}
      />
      <span
        style={{
          width: 16,
          height: 16,
          borderRadius: '50%',
          border: '2px solid #2f6c71',
          background: gender === value ? '#2f6c71' : 'transparent',
          display: 'inline-block',
        }}
      />
      {children}
    </label>
  );

  return (
    <div style={styles.modalOverlay} onClick={close}>
      <div
        style={{ ...styles.modalContent, maxWidth: 620, padding: 0 }}
        onClick={e => e.stopPropagation()}
      >
        {/* шапка */}
        <div
          style={{
            background: '#2f6c71',
            color: '#fff',
            padding: '18px 24px',
            fontSize: 28,
            fontWeight: 700,
            textAlign: 'center',
            position: 'relative',
          }}
        >
          Новый&nbsp;снимок
          <span
            onClick={close}
            style={{
              position: 'absolute',
              right: 22,
              top: 14,
              fontSize: 34,
              lineHeight: 0,
              cursor: 'pointer',
            }}
          >
            ×
          </span>
        </div>

        {/* форма */}
        <form onSubmit={handleSubmit} style={{ padding: 40 }}>
          <Row label="Имя пациента">
            <input
              style={{ ...styles.input, width: '100%' }}
              value={name}
              onChange={e => setName(e.target.value)}
              required
            />
          </Row>

          <Row label="Пол">
            <div style={{ display: 'flex', gap: 20 }}>
              <Radio value="Male">Мужской</Radio>
              <Radio value="Female">Женский</Radio>
            </div>
          </Row>

          <Row label="Возраст">
            <input
              style={{ ...styles.input, width: '100%' }}
              type="number"
              min="0"
              value={age}
              onChange={e => setAge(e.target.value)}
              required
            />
          </Row>

          <Row label="Дата снимка">
            <input
              style={{ ...styles.input, width: '100%' }}
              type="date"
              value={date}
              onChange={e => setDate(e.target.value)}
              required
            />
          </Row>

          <Row label="Файл снимка (диком / nifti)">
            <input
              type="file"
              accept=".dcm,.nii,.nii.gz"
              onChange={e => setFile(e.target.files[0])}
              required
              style={{
                background: '#666',
                color: '#fff',
                padding: '10px 18px',
                border: 'none',
                cursor: 'pointer',
              }}
            />
          </Row>

          {/* кнопка */}
          <div style={{ textAlign: 'center', marginTop: 20 }}>
            <button
              type="submit"
              style={{
                ...styles.uploadButton,
                fontFamily: 'inherit',
                padding: '12px 70px',
                opacity: loading ? 0.6 : 1,
              }}
              disabled={loading}
            >
              {loading ? 'Загрузка…' : 'Загрузить'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
