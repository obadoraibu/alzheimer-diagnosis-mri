import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import SignUp from './components/SignUp';
import SignIn from './components/SignIn';
import Home from './components/Home';
// import ResetRequest from './components/ResetRequest';
// import ResetForm from './components/ResetForm';
// import ScanList from './components/ScanList';
// import ScanDetail from './components/ScanDetail';
// import UploadScan from './components/UploadScan';

function App() {
  return (
    <Router>
      <Routes>
        {/* Аутентификация */}
        <Route path="/sign-in" element={<SignIn />} />
        <Route path="/complete-invite/:code" element={<SignUp />} />

        {/* Приватные (защищённые) маршруты */}
        <Route path="/home" element={<Home />} />
        {/* <Route path="/scans" element={<ScanList />} />
        <Route path="/scans/:id" element={<ScanDetail />} />
        <Route path="/upload" element={<UploadScan />} />
        <Route path="/profile" element={<Profile />} /> */}

        {/* Fallback */}
        <Route path="*" element={<SignIn />} />
      </Routes>
    </Router>
  );
}

export default App;
