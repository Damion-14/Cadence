import { Link } from 'react-router';
import { useAuth } from '../../context/AuthContext';
import { Button } from '../ui/Button';

export function Header() {
  const { user, logout, isAuthenticated } = useAuth();

  return (
    <header className="bg-white shadow">
      <nav className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div className="flex justify-between items-center">
          <Link to="/" className="text-2xl font-bold text-blue-600">
            Cadence
          </Link>

          {isAuthenticated ? (
            <div className="flex items-center gap-6">
              <Link
                to="/workout"
                className="text-gray-700 hover:text-blue-600 transition-colors"
              >
                Workout
              </Link>
              <Link
                to="/history"
                className="text-gray-700 hover:text-blue-600 transition-colors"
              >
                History
              </Link>
              <Link
                to="/stats"
                className="text-gray-700 hover:text-blue-600 transition-colors"
              >
                Stats
              </Link>

              <div className="flex items-center gap-4 ml-4 pl-4 border-l">
                <span className="text-sm text-gray-600">{user?.username}</span>
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={logout}
                >
                  Logout
                </Button>
              </div>
            </div>
          ) : (
            <div className="flex items-center gap-4">
              <Link to="/login">
                <Button variant="secondary" size="sm">
                  Login
                </Button>
              </Link>
              <Link to="/register">
                <Button size="sm">Register</Button>
              </Link>
            </div>
          )}
        </div>
      </nav>
    </header>
  );
}
