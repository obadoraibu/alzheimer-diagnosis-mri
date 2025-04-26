/*
 * Home.js — корневой экран.
 * Администратор не видит вкладку «Исследования» и сразу попадает
 * в «Администрирование».
 */

import React, { useEffect, useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { homeStyles as styles } from '../styles/styles';

import { BASE_URL }      from '../constants';
import { useApi }        from '../hooks/useApi';

import HeaderBar         from './HeaderBar';
import StudiesTab        from './StudiesTab';
import ProfileTab        from './ProfileTab';
import AdminTab          from './AdminTab';
import ScanDetailModal   from './ScanDetailModal';
import UploadModal       from './UploadModal';   // ваша модалка загрузки

function Home() {
  const navigate = useNavigate();
  const api      = useApi();

  /* вкладки / роли */
  const [activeTab, setActiveTab] = useState('studies'); // по умолчанию
  const [isAdmin,   setIsAdmin]   = useState(false);     // доступ к /admin
  const [profile,   setProfile]   = useState(null);      // { username, email, role, … }

  /* модалки */
  const [scanModal,  setScanModal]  = useState(null);
  const [uploadOpen, setUploadOpen] = useState(false);

  /* ---------------- Авторизация + профиль ---------------- */
  useEffect(() => {
    if (!localStorage.getItem('accessToken')) {
      navigate('/sign-in');
      return;
    }

    /* проверяем, есть ли доступ к /admin */
    api('/admin/users').then(r => r.ok && setIsAdmin(true));

    /* загружаем профиль (нужна роль) */
    api('/profile').then(r => {
      if (!r.ok) return;
      r.json().then(p => {
        setProfile(p);
        /* если это админ — сразу открываем «Администрирование» */
        if (p.role === 'admin') setActiveTab('admin');
      });
    });
    // eslint-disable-next-line
  }, []);

  /* ---------------- logout ---------------- */
  const logout = useCallback(
    () =>
      fetch(`${BASE_URL}/revoke`, { method: 'POST', credentials: 'include' }).finally(() => {
        localStorage.clear();
        navigate('/sign-in');
      }),
    [navigate]
  );

  /* ---------------- детали снимка ---------------- */
  const openScanDetail = async id => {
    const r = await api(`/scans/${id}`);
    if (r.ok) setScanModal(await r.json());
  };

  /* скрыть ли вкладку «Исследования»? */
  const hideStudies = profile?.role === 'admin';

  /* ---------------- render ---------------- */
  return (
    <div style={{ ...styles.pageWrapper, fontFamily: 'Georgia, serif' }}>
      <HeaderBar
        activeTab   ={activeTab}
        setActiveTab={setActiveTab}
        isAdmin     ={isAdmin}
        hideStudies ={hideStudies}
        onLogout    ={logout}
      />

      <div style={styles.container}>
        {/* === Исследования (для non-admin) === */}
        {!hideStudies && activeTab === 'studies' && (
          <>
            <div style={styles.titleRow}>
              <h2 style={{ ...styles.title, margin: '0 auto' }}>Исследования</h2>
              <button style={styles.uploadButton} onClick={() => setUploadOpen(true)}>
                Загрузить
              </button>
            </div>

            <StudiesTab api={api} onOpen={openScanDetail} />
          </>
        )}

        {/* === Профиль === */}
        {activeTab === 'profile' && <ProfileTab user={profile} />}

        {/* === Администрирование === */}
        {activeTab === 'admin' && isAdmin && <AdminTab api={api} />}
      </div>

      {/* ---------- модалки ---------- */}
      {scanModal  && <ScanDetailModal scan={scanModal} close={() => setScanModal(null)} />}
      {uploadOpen && <UploadModal    api={api}        close={() => setUploadOpen(false)} />}
    </div>
  );
}

export default Home;
