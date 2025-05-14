// src/styles/styles.js

export const formStyles = {
  container: {
    minHeight: '100vh',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#f2f2f2',
    padding: '1em',
  },
  formBox: {
    backgroundColor: '#fff',
    borderRadius: '0.5em',
    boxShadow: '0 2px 6px rgba(0,0,0,0.1)',
    padding: '2em',
    maxWidth: '400px',
    width: '100%',
    margin: '1em',
  },
  heading: {
    marginBottom: '1em',
    textAlign: 'center',
    color: '#004D4D',
  },
  formGroup: {
    marginBottom: '1em',
  },
  label: {
    display: 'block',
    marginBottom: '0.5em',
    fontWeight: 'bold',
    color: '#333',
  },
  input: {
    width: '100%',
    padding: '0.75em',
    border: '1px solid #ccc',
    borderRadius: '0.25em',
    fontSize: '1em',
    boxSizing: 'border-box',
  },
  button: {
    width: '100%',
    padding: '0.75em',
    backgroundColor: '#008080',
    color: '#fff',
    border: 'none',
    borderRadius: '0.25em',
    fontSize: '1em',
    cursor: 'pointer',
    marginTop: '1em',
  },
  message: {
    marginTop: '1em',
    textAlign: 'center',
    color: 'red',
  },
  linkText: {
    marginTop: '1em',
    textAlign: 'center',
    fontSize: '0.9em',
  },
  link: {
    color: '#008080',
    textDecoration: 'none',
    fontWeight: 'bold',
  },
};


  

export const homeStyles = {
    pageWrapper: {
      backgroundColor: '#f2f2f2',
      minHeight: '100vh',
    },
    header: {
      backgroundColor: '#008080', 
      color: '#fff',
      padding: '1em',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'space-between',
    },
    logo: {
      fontSize: '1.25em',
      fontWeight: 'bold',
    },
    navItems: {
      display: 'flex',
      gap: '1em',
    },
    navItem: {
      cursor: 'pointer',
      fontWeight: 'bold',
    },
    container: {
      maxWidth: '1100px',
      margin: '2em auto',
      backgroundColor: '#fff',
      padding: '2em',
      borderRadius: '0.5em',
    },
    titleRow: {
      display: 'flex',
      justifyContent: 'space-between',
      alignItems: 'center',
      marginBottom: '1em',
    },
    title: {
      margin: 0,
      color: '#008080',
    },
    uploadButton: {
      backgroundColor: '#008080',
      color: '#fff',
      border: 'none',
      padding: '0.75em 1.5em',
      borderRadius: '0.25em',
      fontSize: '1em',
      cursor: 'pointer',
    },
    table: {
      width: '100%',
      borderCollapse: 'collapse',
    },
    th: {
      backgroundColor: '#e6e6e6',
      padding: '0.75em',
      border: '1px solid #ccc',
      textAlign: 'left',
    },
    td: {
      padding: '0.75em',
      border: '1px solid #ccc',
    },
    // Modal overlay styles for file upload
    modalOverlay: {
      position: 'fixed',
      top: 0,
      left: 0,
      right: 0,
      bottom: 0,
      backgroundColor: 'rgba(0,0,0,0.5)',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      zIndex: 1000,
    },
    modalContent: {
      backgroundColor: '#fff',
      padding: '2em',
      borderRadius: '0.5em',
      maxWidth: '500px',
      width: '90%',
      textAlign: 'center',
    },
    dropArea: {
      border: '2px dashed #008080',
      borderRadius: '0.5em',
      padding: '2em',
      marginTop: '1em',
      cursor: 'pointer',
      backgroundColor: '#f9f9f9',
    },
    closeButton: {
      marginTop: '1em',
      backgroundColor: '#ccc',
      border: 'none',
      borderRadius: '0.25em',
      padding: '0.5em 1em',
      cursor: 'pointer',
    },
    // Styles for the profile view
    profileField: {
      marginBottom: '1em',
      fontSize: '1em',
    },
    profileLabel: {
      fontWeight: 'bold',
      color: '#333',
    },
  };