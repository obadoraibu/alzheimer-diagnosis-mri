/* Цвет для различных статусов */
export const getStatusColor = (st = '') => {
    const s = st.toLowerCase();
    if (s === 'done' || s === 'active')            return '#2e7d32'; // зелёный
    if (s === 'invited')                           return '#ff9800'; // оранжевый
    if (s === 'processing' || s === 'in_progress') return '#1976d2'; // синий
    if (s === 'suspended' || s === 'error')        return '#d32f2f'; // красный
    return '#666';                                                 // серый
  };
  