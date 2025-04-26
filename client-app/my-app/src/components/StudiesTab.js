/*
 * StudiesTab – список исследований + фильтры + пагинация
 * props:
 *    api    – fetch-wrapper
 *    onOpen – (id) => void
 */
import React, { useEffect, useState } from 'react';
import { homeStyles as styles } from '../styles/styles';
import { PAGE_LIMIT }            from '../constants';
import { getStatusColor }        from '../utils/statusColor';

function StudiesTab({ api, onOpen }) {
  const [scans,   setScans]   = useState([]);
  const [page,    setPage]    = useState(0);
  const [hasMore, setHasMore] = useState(false);

  /* фильтры */
  const [searchId, setSearchId] = useState('');
  const [dateFrom, setDateFrom] = useState('');
  const [dateTo,   setDateTo]   = useState('');

  const hasFilters = searchId || dateFrom || dateTo;

  /* ---------------- загрузка ---------------- */
  useEffect(() => {
    (async () => {
      const qs = new URLSearchParams();
      qs.append('limit',  PAGE_LIMIT);
      qs.append('offset', page * PAGE_LIMIT);
      if (searchId) qs.append('id',            searchId);
      if (dateFrom) qs.append('uploaded_from', dateFrom);
      if (dateTo)   qs.append('uploaded_to',   dateTo);

      const r = await api(`/scans?${qs.toString()}`);
      if (!r.ok) { setScans([]); setHasMore(false); return; }

      const raw   = await r.json();                 // [], {data:[]}, null …
      const full  = Array.isArray(raw) ? raw : raw?.data ?? [];
      const safe  = Array.isArray(full) ? full : []; // гарантировано массив

      const start = page * PAGE_LIMIT;
      setScans(safe.slice(start, start + PAGE_LIMIT));
      setHasMore(start + PAGE_LIMIT < safe.length);
    })();
  }, [api, page, searchId, dateFrom, dateTo]);

  /* ---------------- UI ---------------- */
  return (
    <>
      {/* фильтр */}
      <form
        onSubmit={e => { e.preventDefault(); setPage(0); }}
        style={{ display:'flex', gap:8, marginBottom:12 }}
      >
        <input style={{ ...styles.input, width:80 }} placeholder="ID"
               value={searchId} onChange={e=>setSearchId(e.target.value.replace(/\D+/g,''))}/>
        <input style={{ ...styles.input, width:140 }} type="date"
               value={dateFrom} onChange={e=>setDateFrom(e.target.value)}/>
        <span style={{ lineHeight:'32px' }}>—</span>
        <input style={{ ...styles.input, width:140 }} type="date"
               value={dateTo} onChange={e=>setDateTo(e.target.value)}/>
        <button style={styles.button} type="submit">🔍</button>
        {hasFilters && (
          <button
            type="button"
            style={{ ...styles.button, background:'#ccc', color:'#000' }}
            onClick={() => { setSearchId(''); setDateFrom(''); setDateTo(''); setPage(0); }}
          >✕</button>
        )}
      </form>

      {/* таблица */}
      <table style={styles.table}>
        <thead><tr>
          <th style={styles.th}>ID</th>
          <th style={styles.th}>Пациент</th>
          <th style={styles.th}>Пол/Возраст</th>
          <th style={styles.th}>Дата снимка</th>
          <th style={styles.th}>Загружено</th>
          <th style={styles.th}>Статус</th>
        </tr></thead>
        <tbody>
          {scans.map(s=>(
            <tr key={s.ID} onClick={()=>onOpen(s.ID)} style={{ cursor:'pointer' }}>
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

      {/* пагинация */}
      <div style={{ display:'flex', gap:16, justifyContent:'center', marginTop:12 }}>
        <button style={{ ...styles.button, opacity:page===0?0.5:1 }} disabled={page===0}
                onClick={()=>setPage(p=>Math.max(0,p-1))}>← Назад</button>
        <span>Страница {page+1}</span>
        <button style={{ ...styles.button, opacity:hasMore?1:0.5 }} disabled={!hasMore}
                onClick={()=>setPage(p=>p+1)}>Вперёд →</button>
      </div>
    </>
  );
}

export default StudiesTab;
