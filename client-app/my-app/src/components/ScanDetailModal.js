import React, { useState } from 'react';
import { homeStyles as styles } from '../styles/styles';
import jsPDF from 'jspdf';
import html2canvas from 'html2canvas';

const Row = ({ label, children, dark }) => (
  <tr>
    <td
      style={{
        padding: '10px 14px',
        fontWeight: 600,
        color: '#fff',
        background: dark ? '#2f6c71' : '#2f6c71',
        width: 170,
        whiteSpace: 'nowrap',
      }}
    >
      {label}
    </td>
    <td style={{ padding: '10px 14px', background: '#fff' }}>{children}</td>
  </tr>
);

export default function ScanDetailModal({ scan, close }) {
  const [downloading, setDownloading] = useState(false);

  if (!scan) return null;

  const gradcam = scan.gradcam_url || scan.GradCAMURL || '';

  const downloadPdf = async () => {
    setDownloading(true);
    try {
      const node = document.getElementById('scan-info-table');
      const canvas = await html2canvas(node, { scale: 2 });
      const imgData = canvas.toDataURL('image/png');

      const doc = new jsPDF('p', 'pt', 'a4');
      doc.setFont('Helvetica', 'bold');
      doc.setFontSize(18);
      doc.text(`#${scan.ID || scan.id}`, 40, 50);


      doc.addImage(imgData, 'PNG', 40, 70, 515, 0);

      if (gradcam) {
        const img = await toDataUrl(gradcam);
        doc.addImage(img, 'PNG', 40, 400, 400, 0); 
      }

      doc.save(`scan_${scan.ID || scan.id}.pdf`);
    } finally {
      setDownloading(false);
    }
  };


  const toDataUrl = url =>
    new Promise(resolve => {
      const img = new Image();
      img.crossOrigin = 'Anonymous';
      img.onload = function () {
        const c = document.createElement('canvas');
        c.width = this.naturalWidth;
        c.height = this.naturalHeight;
        c.getContext('2d').drawImage(this, 0, 0);
        resolve(c.toDataURL('image/png'));
      };
      img.src = url;
    });


  return (
    <div style={styles.modalOverlay} onClick={close}>
      <div
        style={{ ...styles.modalContent, maxWidth: 600, padding: 0 }}
        onClick={e => e.stopPropagation()}
      >
        {/* шапка */}
        <div
          style={{
            background: '#2f6c71',
            color: '#fff',
            padding: '16px 24px',
            fontSize: 24,
            fontWeight: 700,
            textAlign: 'center',
            position: 'relative',
          }}
        >
          Снимок №{scan.ID || scan.id}
          <span
            onClick={close}
            style={{
              position: 'absolute',
              right: 16,
              top: 14,
              fontSize: 28,
              cursor: 'pointer',
              lineHeight: 0,
            }}
          >
            ×
          </span>
        </div>

        {/* таблица */}
        <table
          id="scan-info-table"
          style={{
            width: '100%',
            borderCollapse: 'collapse',
            border: '1px solid #ccc',
            tableLayout: 'fixed',
          }}
        >
          <tbody>
            <Row label="Пациент">{scan.patient_name || scan.PatientName}</Row>
            <Row label="Пол">
              {(scan.patient_gender || scan.PatientGender) === 'Male'
                ? 'Мужской'
                : (scan.patient_gender || scan.PatientGender) === 'Female'
                ? 'Женский'
                : '—'}
            </Row>
            <Row label="Возраст">{scan.patient_age ?? scan.PatientAge}</Row>
            <Row label="Дата снимка">
              {new Date(scan.scan_date || scan.ScanDate).toLocaleDateString()}
            </Row>
            <Row label="Статус">
              {{
                queued: 'Ожидает анализа',
                processing: 'Обрабатывается',
                done: 'Готов',
                error: 'Ошибка',
              }[scan.status || scan.Status] || (scan.status || scan.Status)}
            </Row>
            {scan.diagnosis != null && <Row label="Диагноз">{scan.diagnosis}</Row>}
            {scan.confidence != null && (
              <Row label="Достоверность">{(scan.confidence * 100).toFixed(1)}%</Row>
            )}
          </tbody>
        </table>

        {/* кнопка */}
        <div style={{ textAlign: 'center', padding: '32px 0 40px' }}>
          <button
            onClick={downloadPdf}
            disabled={downloading}
            style={{
              ...styles.uploadButton,
              fontFamily: 'inherit',
              padding: '12px 48px',
              opacity: downloading ? 0.5 : 1,
            }}
          >
            {downloading ? 'Формирование…' : 'Скачать PDF'}
          </button>
        </div>
      </div>
    </div>
  );
}
