import { BrowserRouter, Routes, Route } from 'react-router-dom';

import  Login  from './pages/Login';
import  Dashboard  from './pages/Dashboard';
import  Register  from './pages/Register';
import  CreateTracker  from './pages/CreateTracker';

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/create" element={<CreateTracker />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="*" element={<Dashboard />} />
      </Routes>
    </BrowserRouter>
  );
}