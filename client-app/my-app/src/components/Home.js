/*  =======================================================================
    Home.js — единый шрифт + цветовые статусы + всё, что было раньше
    ======================================================================= */
    import React, { useEffect, useState } from 'react';
    import { useNavigate } from 'react-router-dom';
    import { homeStyles as styles } from '../styles/styles';
    
    /* базовые константы */
    const BASE_URL   = 'http://localhost:8080';
    const PAGE_LIMIT = 20;
    
    /* цветовая подсветка статусов */
    const getStatusColor = (st) => {
      const s = (st || '').toString().toLowerCase();
      if (s === 'done' || s === 'active')            return '#2e7d32'; // green
      if (s === 'invited')                           return '#ff9800'; // orange
      if (s === 'processing' || s === 'in_progress') return '#1976d2'; // blue
      if (s === 'suspended' || s === 'error')        return '#d32f2f'; // red
      return '#666';                                                // gray
    };
    
    function Home() {
      const navigate = useNavigate();
    
      /* -------- auth / навигация -------- */
      const [activeTab, setActiveTab] = useState('studies');
      const [isAdmin,   setIsAdmin]   = useState(false);
    
      /* -------- исследования -------- */
      const [scans, setScans]               = useState([]);
      const [isUploadOpen, setUpload]       = useState(false);
      const [isDetailOpen, setDetail]       = useState(false);
      const [scanDetail,   setScanDetail]   = useState(null);
    
      /* upload-form */
      const [patientName,   setPN] = useState('');
      const [patientGender, setPG] = useState('Male');
      const [patientAge,    setPA] = useState('');
      const [scanDate,      setSD] = useState('');
      const [selectedFile,  setSF] = useState(null);
    
      /* фильтр + пагинация исследований */
      const [searchId, setSearchId] = useState('');
      const [dateFrom, setDateFrom] = useState('');
      const [dateTo,   setDateTo]   = useState('');
      const [page,     setPage]     = useState(0);
      const [hasMore,  setHasMore]  = useState(false);
    
      /* -------- профиль -------- */
      const [userData, setUserData] = useState(null);
    
      /* -------- admin users -------- */
      const [users,        setUsers]        = useState([]);
      const [userPage,     setUserPage]     = useState(0);
      const [userHasMore,  setUserHasMore]  = useState(false);
    
      /* invite / edit модалки */
      const [isInviteOpen, setInvite] = useState(false);
      const [inviteName,   setInvName]  = useState('');
      const [inviteEmail,  setInvEmail] = useState('');
      const [inviteRole,   setInvRole]  = useState('doctor');
    
      const [isEditOpen, setEdit] = useState(false);
      const [editUser,   setEU]   = useState(null);
      const [editName,   setEN]   = useState('');
      const [editRole,   setER]   = useState('doctor');
      const [editStatus, setES]   = useState('active');
    
      /* =====================================================================
         REFRESH-WRAPPER  (auto-refresh accessToken)
         ===================================================================== */
      const logoutAndRedirect = () => { localStorage.clear(); navigate('/sign-in'); };
    
      const refreshTokens = async () => {
        const refresh = localStorage.getItem('refreshToken');
        if (!refresh) return false;
        const res = await fetch(`${BASE_URL}/refresh`, {
          method:'POST',
          headers:{'Content-Type':'application/json'},
          body:JSON.stringify({ refresh }),
          credentials:'include',
        });
        if (!res.ok) return false;
        const t = await res.json();
        t.access  && localStorage.setItem('accessToken',  t.access);
        t.refresh && localStorage.setItem('refreshToken', t.refresh);
        return true;
      };
    
      const api = async (url, opts = {}, retry = true) => {
        const res = await fetch(url.startsWith('http') ? url : `${BASE_URL}${url}`, {
          ...opts,
          headers:{ ...(opts.headers||{}), Authorization:`Bearer ${localStorage.getItem('accessToken')}` },
        });
        if (res.status !== 401 || !retry) return res;
        if (!(await refreshTokens())) { logoutAndRedirect(); return res; }
        return api(url, opts, false);
      };
    
      /* =====================================================================
         INITIAL AUTH CHECK
         ===================================================================== */
      useEffect(() => {
        if (!localStorage.getItem('accessToken')) return navigate('/sign-in');
        api('/admin/users').then(r => r.ok && setIsAdmin(true));
        // eslint-disable-next-line
      }, []);
    
      /* =====================================================================
         LOAD PER TAB  (и при изменении фильтров/страниц)
         ===================================================================== */
      useEffect(() => {
        if (activeTab==='studies') fetchScans();
        if (activeTab==='profile') fetchProfile();
        if (activeTab==='admin' && isAdmin) fetchUsers();
        // eslint-disable-next-line
      }, [activeTab, isAdmin, page, searchId, dateFrom, dateTo, userPage]);
    
      /* =====================================================================
         SCANS — fetch + поиск + пагинация
         ===================================================================== */
      const fetchScans = async () => {
        const qs = new URLSearchParams();
        qs.append('limit', PAGE_LIMIT);
        qs.append('offset', page*PAGE_LIMIT);
        if (searchId) qs.append('id', searchId);
        if (dateFrom) qs.append('uploaded_from', dateFrom);
        if (dateTo)   qs.append('uploaded_to',   dateTo);
    
        const r = await api(`/scans?${qs.toString()}`);
        if (!r.ok) { setScans([]); setHasMore(false); return; }
    
        const raw  = await r.json();
        const full = Array.isArray(raw) ? raw : raw?.data ?? [];
        const safe = Array.isArray(full) ? full : [];
    
        const start = page*PAGE_LIMIT;
        setScans(safe.slice(start, start+PAGE_LIMIT));
        setHasMore(start + PAGE_LIMIT < safe.length);
      };
    
      const fetchScanDetail = async (id) => {
        const r = await api(`/scans/${id}`);
        if (r.ok) { setScanDetail(await r.json()); setDetail(true); }
      };
    
      /* =====================================================================
         PROFILE
         ===================================================================== */
      const fetchProfile = async () => {
        const r = await api('/profile');
        r.ok && setUserData(await r.json());
      };
    
      /* =====================================================================
         USERS  (admin list + pagination)
         ===================================================================== */
      const fetchUsers = async () => {
        const qs = new URLSearchParams();
        qs.append('limit', PAGE_LIMIT);
        qs.append('offset', userPage*PAGE_LIMIT);
    
        const r = await api(`/admin/users?${qs.toString()}`);
        if (!r.ok) { setUsers([]); setUserHasMore(false); return; }
    
        const raw  = await r.json();
        const full = Array.isArray(raw) ? raw : raw.data || raw.users || [];
        const safe = Array.isArray(full) ? full : [];
    
        const start = userPage*PAGE_LIMIT;
        setUsers(safe.slice(start, start+PAGE_LIMIT));
        setUserHasMore(start + PAGE_LIMIT < safe.length);
      };
    
      /* invite / edit / delete */
      const inviteUser = async (e) => {
        e.preventDefault();
        const r = await api('/admin/users', {
          method:'POST', headers:{'Content-Type':'application/json'},
          body:JSON.stringify({ username:inviteName, email:inviteEmail, role:inviteRole }),
        });
        if (r.ok) { setInvite(false); setInvName(''); setInvEmail(''); fetchUsers(); }
      };
      const openEdit   = (u)=>{ setEU(u); setEN(u.username); setER(u.role); setES(u.status||'active'); setEdit(true); };
      const updateUser = async (e)=>{ e.preventDefault();
        const r = await api(`/admin/users/${editUser.id||editUser.ID}`, {
          method:'PUT', headers:{'Content-Type':'application/json'},
          body:JSON.stringify({ username:editName, role:editRole, status:editStatus }),
        });
        if (r.ok) { setEdit(false); fetchUsers(); }
      };
      const deleteUser = async (u)=>{ if(!window.confirm(`Сделать пользователя ${u.username} suspended?`)) return;
        const r = await api(`/admin/users/${u.id||u.ID}`, { method:'DELETE' }); r.ok && fetchUsers(); };
    
      /* =====================================================================
         UPLOAD SCAN
         ===================================================================== */
      const uploadScan = async (e) => {
        e.preventDefault();
        if (!selectedFile) return;
        const f = new FormData();
        f.append('patient_name', patientName);
        f.append('patient_gender', patientGender);
        f.append('patient_age', patientAge);
        f.append('scan_date', scanDate);
        f.append('file', selectedFile);
        const r = await api('/upload', { method:'POST', body:f });
        if (r.ok) { setUpload(false); fetchScans(); }
      };
    
      /* =====================================================================
         LOGOUT
         ===================================================================== */
      const logout = () =>
        fetch(`${BASE_URL}/revoke`, { method:'POST', credentials:'include' }).finally(logoutAndRedirect);
    
      /* =====================================================================
         UI helpers
         ===================================================================== */
      const NavBtn = ({ id, children }) => (
        <span
          style={{ ...styles.navItem, borderBottom:activeTab===id?'2px solid #008080':'2px solid transparent' }}
          onClick={()=>{ setActiveTab(id); if(id==='studies')setPage(0); if(id==='admin')setUserPage(0);} }
        >{children}</span>
      );
      const ActionBtn = ({onClick,children,danger=false}) => (
        <button onClick={onClick} style={{ ...styles.actionButton, color: danger?'crimson':'#008080' }}>{children}</button>
      );
    
      const hasFilters = searchId || dateFrom || dateTo;
    
      /* =====================================================================
         RENDER
         ===================================================================== */
      return (
        /* ► единый шрифт: берём fontFamily из styles.button ◄ */
        <div style={{ ...styles.pageWrapper, fontFamily: 'Georgia, serif' }}>
          {/* -------- HEADER -------- */}
          <header style={{ ...styles.header, position:'relative', display:'flex', alignItems:'center' }}>
            <div style={styles.logo}>MRI App</div>
            <div style={{ position:'absolute', left:'50%', transform:'translateX(-50%)',
                          display:'flex', gap:'24px' }}>
              <NavBtn id="studies">Исследования</NavBtn>
              <NavBtn id="profile">Профиль</NavBtn>
              {isAdmin && <NavBtn id="admin">Администрирование</NavBtn>}
            </div>
            <span style={{ ...styles.navItem, marginLeft:'auto' }} onClick={logout}>Выход</span>
          </header>
    
          {/* -------- CONTENT -------- */}
          <div style={styles.container}>
    
            {/* ========== STUDIES ========== */}
            {activeTab==='studies' && (
              <>
                <div style={styles.titleRow}>
                  <h2 style={{ ...styles.title, margin:'0 auto' }}>Исследования</h2>
                  <button style={styles.uploadButton} onClick={()=>setUpload(true)}>Загрузить</button>
                </div>
    
                {/* минималистичный фильтр */}
                <form onSubmit={(e)=>{e.preventDefault(); setPage(0);} }
                      style={{ display:'flex', gap:'8px', marginBottom:'12px' }}>
                  <input style={{ ...styles.input, width:'80px' }} placeholder="ID"
                         value={searchId} onChange={e=>setSearchId(e.target.value.replace(/\D+/g,''))}/>
                  <input style={{ ...styles.input, width:'140px' }} type="date"
                         value={dateFrom} onChange={e=>setDateFrom(e.target.value)}/>
                  <span style={{ lineHeight:'32px' }}>—</span>
                  <input style={{ ...styles.input, width:'140px' }} type="date"
                         value={dateTo} onChange={e=>setDateTo(e.target.value)}/>
                  <button style={styles.button} type="submit">🔍</button>
                  {hasFilters && (
                    <button style={{ ...styles.button, background:'#ccc', color:'#000' }} type="button"
                            onClick={()=>{ setSearchId(''); setDateFrom(''); setDateTo(''); setPage(0); }}>
                      ✕
                    </button>
                  )}
                </form>
    
                {/* таблица исследований */}
                <table style={styles.table}>
                  <thead><tr>
                    <th style={styles.th}>ID</th><th style={styles.th}>Пациент</th><th style={styles.th}>Пол/Возраст</th>
                    <th style={styles.th}>Дата снимка</th><th style={styles.th}>Загружено</th><th style={styles.th}>Статус</th>
                  </tr></thead>
                  <tbody>
                    {scans.map(s=>(
                      <tr key={s.ID} style={{ cursor:'pointer' }} onClick={()=>fetchScanDetail(s.ID)}>
                        <td style={styles.td}>{s.ID}</td>
                        <td style={styles.td}>{s.PatientName}</td>
                        <td style={styles.td}>{s.PatientGender}, {s.PatientAge}</td>
                        <td style={styles.td}>{new Date(s.ScanDate).toLocaleDateString()}</td>
                        <td style={styles.td}>{new Date(s.CreatedAt).toLocaleDateString()}</td>
                        <td style={{ ...styles.td, color:getStatusColor(s.Status) }}>{s.Status}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
    
                {/* пагинация исследований */}
                <div style={{ display:'flex', justifyContent:'center', gap:'16px', marginTop:'12px' }}>
                  <button style={{ ...styles.button, opacity:page===0?0.5:1 }} disabled={page===0}
                          onClick={()=>setPage(p=>Math.max(0,p-1))}>← Назад</button>
                  <span>Страница {page+1}</span>
                  <button style={{ ...styles.button, opacity:hasMore?1:0.5 }} disabled={!hasMore}
                          onClick={()=>setPage(p=>p+1)}>Вперёд →</button>
                </div>
              </>
            )}
    
            {/* ========== PROFILE ========== */}
            {activeTab==='profile' && userData && (
              <>
                <h2 style={{ ...styles.title, margin:'0 auto' }}>Профиль</h2>
                <div style={{ textAlign:'center' }}>
                  <div style={styles.profileField}><span style={styles.profileLabel}>Имя: </span>{userData.username}</div>
                  <div style={styles.profileField}><span style={styles.profileLabel}>Email: </span>{userData.email}</div>
                </div>
              </>
            )}
    
            {/* ========== ADMIN ========== */}
            {activeTab==='admin' && isAdmin && (
              <>
                <div style={styles.titleRow}>
                  <h2 style={{ ...styles.title, margin:'0 auto' }}>Администрирование</h2>
                  <button style={styles.uploadButton} onClick={()=>setInvite(true)}>Пригласить</button>
                </div>
    
                {users.length===0 ? <p style={{ textAlign:'center' }}>Нет пользователей</p> : (
                  <>
                    <table style={styles.table}>
                      <thead><tr>
                        <th style={styles.th}>ID</th><th style={styles.th}>Имя</th><th style={styles.th}>Email</th>
                        <th style={styles.th}>Роль</th><th style={styles.th}>Статус</th><th style={styles.th}>Действия</th>
                      </tr></thead>
                      <tbody>
                        {users.map(u=>(
                          <tr key={u.id||u.ID}>
                            <td style={styles.td}>{u.id||u.ID}</td>
                            <td style={styles.td}>{u.username}</td>
                            <td style={styles.td}>{u.email}</td>
                            <td style={styles.td}>{u.role}</td>
                            <td style={{ ...styles.td, color:getStatusColor(u.status) }}>{u.status}</td>
                            <td style={{ ...styles.td, display:'flex', gap:'6px' }}>
                              <ActionBtn onClick={()=>openEdit(u)}>✏️</ActionBtn>
                              <ActionBtn danger onClick={()=>deleteUser(u)}>🗑️</ActionBtn>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
    
                    {/* пагинация пользователей */}
                    <div style={{ display:'flex', justifyContent:'center', gap:'16px', marginTop:'12px' }}>
                      <button style={{ ...styles.button, opacity:userPage===0?0.5:1 }} disabled={userPage===0}
                              onClick={()=>setUserPage(p=>Math.max(0,p-1))}>← Назад</button>
                      <span>Страница {userPage+1}</span>
                      <button style={{ ...styles.button, opacity:userHasMore?1:0.5 }} disabled={!userHasMore}
                              onClick={()=>setUserPage(p=>p+1)}>Вперёд →</button>
                    </div>
                  </>
                )}
              </>
            )}
          </div>
    
          {/* -------- MODALS (upload / details / invite / edit) -------- */}
          {isUploadOpen && (
            <div style={styles.modalOverlay} onClick={()=>setUpload(false)}>
              <div style={styles.modalContent} onClick={e=>e.stopPropagation()}>
                <h3 style={styles.heading}>Новый снимок</h3>
                <form onSubmit={uploadScan}>
                  <div style={styles.formGroup}><label style={styles.label}>Имя пациента</label>
                    <input style={styles.input} value={patientName} onChange={e=>setPN(e.target.value)} required/></div>
                  <div style={styles.formGroup}><label style={styles.label}>Пол</label>
                    <select style={styles.input} value={patientGender} onChange={e=>setPG(e.target.value)}>
                      <option value="Male">Мужской</option><option value="Female">Женский</option><option value="Other">Другой</option>
                    </select></div>
                  <div style={styles.formGroup}><label style={styles.label}>Возраст</label>
                    <input style={styles.input} type="number" min="0" value={patientAge}
                           onChange={e=>setPA(e.target.value)} required/></div>
                  <div style={styles.formGroup}><label style={styles.label}>Дата снимка</label>
                    <input style={styles.input} type="date" value={scanDate} onChange={e=>setSD(e.target.value)} required/></div>
                  <div style={styles.formGroup}><label style={styles.label}>Файл (DICOM/NIfTI)</label>
                    <input style={styles.input} type="file" accept=".dcm,.nii,.nii.gz"
                           onChange={e=>setSF(e.target.files[0])} required/></div>
                  <button style={styles.button} type="submit">Загрузить</button>
                </form>
              </div>
            </div>
          )}
    
          {isDetailOpen && scanDetail && (
            <div style={styles.modalOverlay} onClick={()=>setDetail(false)}>
              <div style={styles.modalContent} onClick={e=>e.stopPropagation()}>
                <h3 style={styles.heading}>Снимок #{scanDetail.ID||scanDetail.id}</h3>
                <div style={styles.detailField}><strong>Пациент:</strong> {scanDetail.patient_name||scanDetail.PatientName}</div>
                <div style={styles.detailField}><strong>Пол:</strong> {scanDetail.patient_gender||scanDetail.PatientGender}</div>
                <div style={styles.detailField}><strong>Возраст:</strong> {scanDetail.patient_age??scanDetail.PatientAge}</div>
                <div style={styles.detailField}><strong>Дата снимка:</strong> {new Date(scanDetail.scan_date||scanDetail.ScanDate).toLocaleDateString()}</div>
                <div style={styles.detailField}><strong>Статус:</strong> {scanDetail.status||scanDetail.Status}</div>
                {scanDetail.diagnosis!=null  && <div style={styles.detailField}><strong>Диагноз:</strong> {scanDetail.diagnosis}</div>}
                {scanDetail.confidence!=null && <div style={styles.detailField}><strong>Достоверность:</strong> {(scanDetail.confidence*100).toFixed(1)}%</div>}
                {(scanDetail.gradcam_url||scanDetail.GradCAMURL) &&
                  <img src={scanDetail.gradcam_url||scanDetail.GradCAMURL} alt="Grad-CAM" style={styles.gradCamImage}/>}
              </div>
            </div>
          )}
    
          {isInviteOpen && (
            <div style={styles.modalOverlay} onClick={()=>setInvite(false)}>
              <div style={styles.modalContent} onClick={e=>e.stopPropagation()}>
                <h3 style={styles.heading}>Пригласить пользователя</h3>
                <form onSubmit={inviteUser}>
                  <div style={styles.formGroup}><label style={styles.label}>Имя</label>
                    <input style={styles.input} value={inviteName} onChange={e=>setInvName(e.target.value)} required/></div>
                  <div style={styles.formGroup}><label style={styles.label}>Email</label>
                    <input style={styles.input} type="email" value={inviteEmail} onChange={e=>setInvEmail(e.target.value)} required/></div>
                  <div style={styles.formGroup}><label style={styles.label}>Роль</label>
                    <select style={styles.input} value={inviteRole} onChange={e=>setInvRole(e.target.value)}>
                      <option value="doctor">Врач</option><option value="admin">Админ</option><option value="viewer">Просмотр</option>
                    </select></div>
                  <button style={styles.button} type="submit">Отправить</button>
                </form>
              </div>
            </div>
          )}
    
          {isEditOpen && editUser && (
            <div style={styles.modalOverlay} onClick={()=>setEdit(false)}>
              <div style={styles.modalContent} onClick={e=>e.stopPropagation()}>
                <h3 style={styles.heading}>Редактировать #{editUser.id||editUser.ID}</h3>
                <form onSubmit={updateUser}>
                  <div style={styles.formGroup}><label style={styles.label}>Имя</label>
                    <input style={styles.input} value={editName} onChange={e=>setEN(e.target.value)}/></div>
                  <div style={styles.formGroup}><label style={styles.label}>Email</label>
                    <input style={styles.input} value={editUser.email} disabled/></div>
                  <div style={styles.formGroup}><label style={styles.label}>Роль</label>
                    <select style={styles.input} value={editRole} onChange={e=>setER(e.target.value)}>
                      <option value="doctor">Врач</option><option value="admin">Админ</option><option value="viewer">Просмотр</option>
                    </select></div>
                  <div style={styles.formGroup}><label style={styles.label}>Статус</label>
                    <select style={styles.input} value={editStatus} onChange={e=>setES(e.target.value)}>
                      <option value="active">active</option><option value="invited">invited</option><option value="suspended">suspended</option>
                    </select></div>
                  <button style={styles.button} type="submit">Сохранить</button>
                </form>
              </div>
            </div>
          )}
        </div>
      );
    }
    
    export default Home;
    