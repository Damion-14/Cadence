import { Outlet } from 'react-router';
import { AuthProvider } from '../context/AuthContext';
import { WorkoutProvider } from '../context/WorkoutContext';
import { Header } from '../components/layout/Header';

export default function Layout() {
  return (
    <AuthProvider>
      <WorkoutProvider>
        <div className="min-h-screen bg-gray-50">
          <Header />
          <main>
            <Outlet />
          </main>
        </div>
      </WorkoutProvider>
    </AuthProvider>
  );
}