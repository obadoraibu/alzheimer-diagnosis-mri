/* =====================================================================
   Home.js — корневой экран
   ===================================================================== */
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
   import UploadModal     from './UploadModal';   // модалка загрузки скана
   import { getStatusColor } from '../utils/statusColor';
   
   /* маленький универсальный ряд формы */
   const Row = ({ label, children }) => (
     <div style={styles.formGroup}>
       <label style={styles.label}>{label}</label>
       {children}
     </div>
   );
   
   function Home() {
     const navigate = useNavigate();
     const api      = useApi();
   
     /* ------------- вкладки / роли ------------- */
     const [activeTab, setActiveTab] = useState('studies');
     const [isAdmin,   setIsAdmin]   = useState(false);
     const [profile,   setProfile]   = useState(null);      // { username, email, role }
   
     /* ------------- модалки ------------- */
     const [scanModal,  setScanModal]  = useState(null);
     const [uploadOpen, setUploadOpen] = useState(false);
     const [inviteOpen, setInviteOpen] = useState(false);
     const [editUser,   setEditUser]   = useState(null);    // объект пользователя
   
     /* поля формы Invite */
     const [invName,  setInvName]  = useState('');
     const [invEmail, setInvEmail] = useState('');
     const [invRole,  setInvRole]  = useState('doctor');
   
     /* поля формы Edit */
     const [editName,   setEditName]   = useState('');
     const [editRole,   setEditRole]   = useState('doctor');
     const [editStatus, setEditStatus] = useState('active');
   
     /* =================================================================
        Авторизация и профиль
        ================================================================= */
     useEffect(() => {
       if (!localStorage.getItem('accessToken')) {
         navigate('/sign-in');
         return;
       }
   
       api('/admin/users').then(r => r.ok && setIsAdmin(true));
   
       api('/profile').then(r => {
         if (!r.ok) return;
         r.json().then(p => {
           setProfile(p);
           if (p.role === 'admin') setActiveTab('admin');   // админ стартует с «Администрирование»
         });
       });
       // eslint-disable-next-line
     }, []);
   
     /* ------------------------------------------------------------------
        logout
        ------------------------------------------------------------------ */
     const logout = useCallback(
       () =>
         fetch(`${BASE_URL}/revoke`, { method: 'POST', credentials: 'include' }).finally(() => {
           localStorage.clear();
           navigate('/sign-in');
         }),
       [navigate]
     );
   
     /* ------------------------------------------------------------------
        открыть детали снимка
        ------------------------------------------------------------------ */
     const openScanDetail = async id => {
       const r = await api(`/scans/${id}`);
       if (r.ok) setScanModal(await r.json());
     };
   
     /* ------------------------------------------------------------------
        admin callbacks
        ------------------------------------------------------------------ */
     const handleEdit = u => {
       setEditUser(u);
       setEditName(u.username);
       setEditRole(u.role);
       setEditStatus(u.status);
     };
   
     const handleDelete = async u => {
       if (!window.confirm(`Сделать пользователя ${u.username} suspended?`)) return;
       await api(`/admin/users/${u.id || u.ID}`, { method: 'DELETE' });
     };
   
     /* ------------------------------------------------------------------
        Invite (POST /admin/users)
        ------------------------------------------------------------------ */
     const inviteUser = async e => {
       e.preventDefault();
       const r = await api('/admin/users', {
         method:'POST',
         headers:{'Content-Type':'application/json'},
         body:JSON.stringify({ username:invName, email:invEmail, role:invRole }),
       });
       if (r.ok) {
         alert('Приглашение отправлено');
         setInviteOpen(false);
         setInvName('');
         setInvEmail('');
       }
     };
   
     /* ------------------------------------------------------------------
        Edit (PUT /admin/users/:id)
        ------------------------------------------------------------------ */
     const saveEditUser = async e => {
       e.preventDefault();
       if (!editUser) return;
       const r = await api(`/admin/users/${editUser.id || editUser.ID}`, {
         method:'PUT',
         headers:{'Content-Type':'application/json'},
         body:JSON.stringify({ username:editName, role:editRole, status:editStatus }),
       });
       if (r.ok) {
         alert('Сохранено');
         setEditUser(null);
       }
     };
   
     /* ------------------------------------------------------------------
        Render
        ------------------------------------------------------------------ */
     const hideStudies = profile?.role === 'admin';
   
     return (
       <div style={{ ...styles.pageWrapper, fontFamily:'Georgia, serif' }}>
         <HeaderBar
           activeTab   ={activeTab}
           setActiveTab={setActiveTab}
           isAdmin     ={isAdmin}
           hideStudies ={hideStudies}
           onLogout    ={logout}
         />
   
         <div style={styles.container}>
           {/* -------- Исследования -------- */}
           {!hideStudies && activeTab === 'studies' && (
             <>
               <div style={styles.titleRow}>
                 <h2 style={{ ...styles.title, margin:'0 auto' }}>Исследования</h2>
                 <button style={styles.uploadButton} onClick={()=>setUploadOpen(true)}>Загрузить</button>
               </div>
               <StudiesTab api={api} onOpen={openScanDetail}/>
             </>
           )}
   
           {/* -------- Профиль -------- */}
           {activeTab === 'profile' && <ProfileTab user={profile}/>}
   
           {/* -------- Администрирование -------- */}
           {activeTab === 'admin' && isAdmin && (
             <>
               <div style={styles.titleRow}>
                 <h2 style={{ ...styles.title, margin:'0 auto' }}>Администрирование</h2>
                 <button style={styles.uploadButton} onClick={()=>setInviteOpen(true)}>Пригласить</button>
               </div>
               <AdminTab api={api} onEdit={handleEdit} onDelete={handleDelete}/>
             </>
           )}
         </div>
   
         {/* ---------------- модалки ---------------- */}
         {scanModal  && <ScanDetailModal scan={scanModal} close={()=>setScanModal(null)} />}
         {uploadOpen && <UploadModal    api={api}        close={()=>setUploadOpen(false)} />}
   
         {/* Invite inline */}
         {inviteOpen && (
           <div style={styles.modalOverlay} onClick={()=>setInviteOpen(false)}>
             <div style={styles.modalContent} onClick={e=>e.stopPropagation()}>
               <h3 style={styles.heading}>Пригласить пользователя</h3>
               <form onSubmit={inviteUser}>
                 <Row label="Имя">
                   <input style={styles.input} value={invName}
                          onChange={e=>setInvName(e.target.value)} required/>
                 </Row>
                 <Row label="Email">
                   <input style={styles.input} type="email" value={invEmail}
                          onChange={e=>setInvEmail(e.target.value)} required/>
                 </Row>
                 <Row label="Роль">
                   <select style={styles.input} value={invRole}
                           onChange={e=>setInvRole(e.target.value)}>
                     <option value="doctor">Врач</option>
                     <option value="admin">Админ</option>
                     <option value="viewer">Просмотр</option>
                   </select>
                 </Row>
                 <button style={styles.button} type="submit">Отправить</button>
               </form>
             </div>
           </div>
         )}
   
         {/* Edit inline */}
         {editUser && (
           <div style={styles.modalOverlay} onClick={()=>setEditUser(null)}>
             <div style={styles.modalContent} onClick={e=>e.stopPropagation()}>
               <h3 style={styles.heading}>Редактировать #{editUser.id || editUser.ID}</h3>
               <form onSubmit={saveEditUser}>
                 <Row label="Имя">
                   <input style={styles.input} value={editName}
                          onChange={e=>setEditName(e.target.value)}/>
                 </Row>
                 <Row label="Email">
                   <input style={styles.input} value={editUser.email} disabled/>
                 </Row>
                 <Row label="Роль">
                   <select style={styles.input} value={editRole}
                           onChange={e=>setEditRole(e.target.value)}>
                     <option value="doctor">Врач</option>
                     <option value="admin">Админ</option>
                     <option value="viewer">Просмотр</option>
                   </select>
                 </Row>
                 <Row label="Статус">
                   <select style={styles.input} value={editStatus}
                           onChange={e=>setEditStatus(e.target.value)}>
                     <option value="active">active</option>
                     <option value="invited">invited</option>
                     <option value="suspended">suspended</option>
                   </select>
                 </Row>
                 <button style={styles.button} type="submit">Сохранить</button>
               </form>
             </div>
           </div>
         )}
       </div>
     );
   }
   
   export default Home;
   