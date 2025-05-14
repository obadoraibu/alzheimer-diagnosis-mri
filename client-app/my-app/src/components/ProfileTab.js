import React from 'react';
import { homeStyles as styles } from '../styles/styles';

export default function ProfileTab({ user }) {
  if (!user) return null;
  return (
    <>
      <h2 style={{ ...styles.title, margin: '0 auto' }}>Профиль</h2>
      <div style={{ textAlign: 'center' }}>
        <div style={styles.profileField}>
          <span style={styles.profileLabel}>Имя:&nbsp;</span>
          {user.username}
        </div>
        <div style={styles.profileField}>
          <span style={styles.profileLabel}>Email:&nbsp;</span>
          {user.email}
        </div>
        <div style={styles.profileField}>
          <span style={styles.profileLabel}>Роль:&nbsp;</span>
          {user.role}
        </div>
        <div style={styles.profileField}>
          <span style={styles.profileLabel}>Статус:&nbsp;</span>
          {user.status}
        </div>
      </div>
    </>
  );
}
