import React, { useEffect, useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { homeStyles as styles } from '../styles/styles';

import { BASE_URL } from '../constants';
import { useApi }   from '../hooks/useApi';

import HeaderBar       from './HeaderBar';
import StudiesTab      from './StudiesTab';
import ProfileTab      from './ProfileTab';
import AdminTab        from './AdminTab';
import ScanDetailModal from './ScanDetailModal';
import UploadModal     from './UploadModal';
import InviteUserModal from './InviteUserModal';
import EditUserModal   from './EditUserModal';

function Home() {
  const navigate = useNavigate();
  const api      = useApi();

  /* вкладки / профиль */
  const [activeTab, setActiveTab] = useState('studies');
  const [isAdmin,   setIsAdmin]   = useState(false);
  const [profile,   setProfile]   = useState(null);

  /* «версии» списков: при ++ заставляют вложенные табы перезагрузиться */
  const [scanVer, setScanVer]  = useState(0);
  const [userVer, setUserVer]  = useState(0);

  /* модалки */
  const [scanModal,  setScanModal]  = useState(null);
  const [uploadOpen, setUploadOpen] = useState(false);
  const [inviteOpen, setInviteOpen] = useState(false);
  const [editUser,   setEditUser]   = useState(null);

  /* ───────── initial auth / profile ───────── */
  useEffect(() => {
    if (!localStorage.getItem('accessToken')) {
      navigate('/sign-in');
      return;
    }

    api('/admin/users').then(r => r.ok && setIsAdmin(true));

    api('/profile').then(async r => {
      if (!r.ok) return;
      const me = (await r.json())?.data;
      if (me) {
        setProfile(me);
        if (me.role === 'admin') setActiveTab('admin');
      }
    });
    // eslint-disable-next-line
  }, []);

  /* logout */
  const logout = useCallback(
    () =>
      fetch(`${BASE_URL}/revoke`, { method: 'POST', credentials: 'include' }).finally(() => {
        localStorage.clear();
        navigate('/sign-in');
      }),
    [navigate]
  );

  /* открыть детали снимка */
  const openScanDetail = async id => {
    const r = await api(`/scans/${id}`);
    if (r.ok) setScanModal((await r.json()).data);
  };

  /* admin callbacks */
  const handleEdit   = u => setEditUser(u);
  const handleDelete = async u => {
    if (!window.confirm(`Сделать пользователя ${u.username} suspended?`)) return;
    const r = await api(`/admin/users/${u.id || u.ID}`, { method: 'DELETE' });
    if (r.ok) setUserVer(v => v + 1);                     // ► обновить таблицу
  };

  /* колбэки, приходящие из модалок */
  const onScanChange = () => setScanVer(v => v + 1);
  const onUserChange = () => setUserVer(v => v + 1);

  const hideStudies = profile?.role === 'admin';

  return (
    <div style={styles.pageWrapper}>
      <HeaderBar
        activeTab   ={activeTab}
        setActiveTab={setActiveTab}
        isAdmin     ={isAdmin}
        hideStudies ={hideStudies}
        onLogout    ={logout}
      />

      <div style={styles.container}>

        {/* ─── Исследования ─── */}
        {!hideStudies && activeTab === 'studies' && (
          <>
            <div style={styles.titleRow}>
              <h2 style={{ ...styles.title, margin:'0 auto' }}>Исследования</h2>
              <button style={styles.uploadButton} onClick={()=>setUploadOpen(true)}>Загрузить</button>
            </div>
            <StudiesTab api={api} onOpen={openScanDetail} refreshKey={scanVer}/>
          </>
        )}

        {/* ─── Профиль ─── */}
        {activeTab === 'profile' && <ProfileTab user={profile}/>}

        {/* ─── Администрирование ─── */}
        {activeTab === 'admin' && isAdmin && (
          <>
            <div style={styles.titleRow}>
              <h2 style={{ ...styles.title, margin:'0 auto' }}>Администрирование</h2>
              <button style={styles.uploadButton} onClick={()=>setInviteOpen(true)}>Пригласить</button>
            </div>
            <AdminTab
              api={api}
              onEdit={handleEdit}
              onDelete={handleDelete}
              refreshKey={userVer}
            />
          </>
        )}
      </div>

      {/* модалки */}
      {scanModal  && <ScanDetailModal scan={scanModal} close={()=>setScanModal(null)} />}
      {uploadOpen && (
        <UploadModal
          api={api}
          close={()=>setUploadOpen(false)}
          onSuccess={onScanChange}
        />
      )}
      {inviteOpen && (
        <InviteUserModal
          api={api}
          close={()=>setInviteOpen(false)}
          onSuccess={onUserChange}
        />
      )}
      {editUser && (
        <EditUserModal
          api={api}
          user={editUser}
          close={()=>setEditUser(null)}
          onSuccess={onUserChange}
        />
      )}
    </div>
  );
}

export default Home;
