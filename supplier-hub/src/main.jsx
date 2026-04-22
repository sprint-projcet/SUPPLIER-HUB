import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App.jsx';
import './index.css'; // Pastikan Anda memiliki file CSS untuk Tailwind

// Menginisialisasi aplikasi React ke elemen dengan id 'root'
ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);