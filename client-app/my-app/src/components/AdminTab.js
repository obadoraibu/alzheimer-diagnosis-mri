/*
   AdminTab – таблица пользователей
   props:
     api       – fetch-wrapper
     onEdit    – (user) => void
     onDelete  – (user) => void
*/
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

function AdminTab({ api, onEdit, onDelete }) {
  const [users, setUsers] = useState([]);
  const [page,  setPage]  = useState(0);
  const [more,  setMore]  = useState(false);

  /* загрузка */
  useEffect(() => {
    (async () => {
      const qs = new URLSearchParams();
      qs.append('limit', PAGE_LIMIT);
      qs.append('offset', page * PAGE_LIMIT);

      const r = await api(`/admin/users?${qs.toString()}`);
      if (!r.ok) return;
      const data  = await r.json();
      const list  = Array.isArray(data) ? data : data.data || data.users || [];
      const start = page * PAGE_LIMIT;

      setUsers(list.slice(start, start + PAGE_LIMIT));
      setMore(start + PAGE_LIMIT < list.length);
    })();
  }, [api, page]);

  return (
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
              <td style={{ ...styles.td, display:'flex', gap:6 }}>
                <SmallBtn onClick={()=>onEdit(u)}>✏️</SmallBtn>
                <SmallBtn danger onClick={()=>onDelete(u)}>🗑️</SmallBtn>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

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
