import React from 'react';
import { homeStyles as styles } from '../styles/styles';

const NavBtn = ({ id, active, hidden, onClick, children }) => {
  if (hidden) return null;                       
  return (
    <span
      onClick={onClick}
      style={{
        ...styles.navItem,
        borderBottom: active ? '2px solid #008080' : '2px solid transparent',
        cursor: 'pointer',
      }}
    >
      {children}
    </span>
  );
};

/**
 * HeaderBar
 *
 * @param {string}  activeTab
 * @param {fn}      setActiveTab
 * @param {boolean} isAdmin
 * @param {boolean} hideStudies 
 * @param {fn}      onLogout
 */
function HeaderBar({ activeTab, setActiveTab, isAdmin, hideStudies, onLogout }) {
  return (
    <header
      style={{
        ...styles.header,
        position: 'relative',
        display: 'flex',
        alignItems: 'center',
      }}
    >
      <div style={styles.logo}>MRI App</div>

      {/* центр навигации */}
      <div
        style={{
          position: 'absolute',
          left: '50%',
          transform: 'translateX(-50%)',
          display: 'flex',
          gap: 24,
        }}
      >
        <NavBtn
          id="studies"
          hidden={hideStudies}
          active={activeTab === 'studies'}
          onClick={() => setActiveTab('studies')}
        >
          Исследования
        </NavBtn>

        <NavBtn
          id="profile"
          active={activeTab === 'profile'}
          onClick={() => setActiveTab('profile')}
        >
          Профиль
        </NavBtn>

        {isAdmin && (
          <NavBtn
            id="admin"
            active={activeTab === 'admin'}
            onClick={() => setActiveTab('admin')}
          >
            Администрирование
          </NavBtn>
        )}
      </div>

      {/* выход */}
      <span
        style={{ ...styles.navItem, marginLeft: 'auto', cursor: 'pointer' }}
        onClick={onLogout}
      >
        Выход
      </span>
    </header>
  );
}

export default HeaderBar;
