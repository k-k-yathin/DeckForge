import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from './context/AuthContext';
import { AppLayout } from './components/layout/AppLayout';
import { Landing } from './pages/Landing';
import { Login } from './pages/Login';
import { Register } from './pages/Register';
import { Dashboard } from './pages/Dashboard';
import { Upload } from './pages/Upload';
import { PresentationPage } from './pages/Presentation';
import { LoadingSpinner } from './components/ui/LoadingSpinner';

/**
 * ProtectedRoute — only renders children if user is logged in.
 */
function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { user, isLoading } = useAuth();
  if (isLoading) return <LoadingSpinner message="Loading..." />;
  if (!user) return <Navigate to="/login" replace />;
  return <>{children}</>;
}

/**
 * App.tsx defines all URL routes (React Router).
 * Public: /, /login, /register
 * Private: /dashboard, /upload, /presentation/:id
 */
export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Landing />} />
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />

      <Route
        element={
          <ProtectedRoute>
            <AppLayout />
          </ProtectedRoute>
        }
      >
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/upload" element={<Upload />} />
        <Route path="/presentation/:id" element={<PresentationPage />} />
      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
