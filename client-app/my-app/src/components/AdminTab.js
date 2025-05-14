import React, { useEffect, useState } from 'react';
import { homeStyles as styles } from '../styles/styles';
import { PAGE_LIMIT }            from '../constants';
import { getStatusColor }        from '../utils/statusColor';

const SmallBtn = ({ danger, ...rest }) => (
  <button
    {...rest}
    style={{ ...styles.actionButton, fontFamily:'inherit', color: danger ? 'crimson' : '#008080' }}
  />
);

function AdminTab({ api, onEdit = () => {}, onDelete = () => {}, refreshKey = 0 }) {
  const [users, setUsers] = useState([]);
  const [page,  setPage]  = useState(0);
  const [more,  setMore]  = useState(false);

  const formatRole = role => {
    switch (role) {
      case 'admin': return 'Администратор';
      case 'doctor': return 'Врач';
      default: return role;
    }
  };

  const formatStatus = status => {
    switch (status) {
      case 'active': return 'Активен';
      case 'suspended': return 'Отстранён';
      case 'invited': return 'Приглашён';
      default: return status;
    }
  };

  useEffect(() => {
    (async () => {
      const qs = new URLSearchParams();
      qs.append('limit',  PAGE_LIMIT);
      qs.append('offset', page * PAGE_LIMIT);

      const r = await api(`/admin/users?${qs.toString()}`);
      if (!r.ok) return;

      const raw   = await r.json();
      const full  = Array.isArray(raw) ? raw : raw?.data ?? raw?.users ?? [];
      const safe  = Array.isArray(full) ? full : [];

      const start = page * PAGE_LIMIT;
      setUsers(safe.slice(start, start + PAGE_LIMIT));
      setMore(start + PAGE_LIMIT < safe.length);
    })();
  }, [api, page, refreshKey]);

  return (
    <>
      <table style={styles.table}>
        <thead><tr>
          <th style={styles.th}>ID</th><th style={styles.th}>Имя</th><th style={styles.th}>Email</th>
          <th style={styles.th}>Роль</th><th style={styles.th}>Статус</th><th style={styles.th}>Действия</th>
        </tr></thead>
        <tbody>
          {users.map(u => (
            <tr key={u.id || u.ID}>
              <td style={styles.td}>{u.id || u.ID}</td>
              <td style={styles.td}>{u.username}</td>
              <td style={styles.td}>{u.email}</td>
              <td style={styles.td}>{formatRole(u.role)}</td>
              <td style={{ ...styles.td, color: getStatusColor(u.status) }}>
                {formatStatus(u.status)}
              </td>
              <td style={{ ...styles.td, display: 'flex', gap: 6 }}>
                <SmallBtn onClick={() => onEdit(u)}>✏️</SmallBtn>
                <SmallBtn danger onClick={() => onDelete(u)}>🗑️</SmallBtn>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* пагинация */}
      <div style={{ display:'flex', gap:16, justifyContent:'center', marginTop:12 }}>
        <button style={{ ...styles.button, opacity:page===0?0.5:1 }} disabled={page===0}
                onClick={()=>setPage(p=>Math.max(0,p-1))}>← Назад</button>
        <span>Страница {page+1}</span>
        <button style={{ ...styles.button, opacity:more?1:0.5 }} disabled={!more}
                onClick={()=>setPage(p=>p+1)}>Вперёд →</button>
      </div>
    </>
  );
}

export default AdminTab;
