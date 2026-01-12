import { Link } from 'react-router';
import { useAuth } from '../context/AuthContext';
import { Button } from '../components/ui/Button';
import { Card } from '../components/ui/Card';

export default function Page() {
  const { isAuthenticated, user } = useAuth();

  if (isAuthenticated) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Welcome back, {user?.username}!
          </h1>
          <p className="text-xl text-gray-600">
            Ready to track your fitness journey?
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-4xl mx-auto">
          <Link to="/workout">
            <Card className="hover:shadow-lg transition-shadow cursor-pointer text-center">
              <h3 className="text-xl font-semibold mb-2">Start Workout</h3>
              <p className="text-gray-600">Log exercises and track your progress</p>
            </Card>
          </Link>

          <Link to="/history">
            <Card className="hover:shadow-lg transition-shadow cursor-pointer text-center">
              <h3 className="text-xl font-semibold mb-2">View History</h3>
              <p className="text-gray-600">See all your completed workouts</p>
            </Card>
          </Link>

          <Link to="/stats">
            <Card className="hover:shadow-lg transition-shadow cursor-pointer text-center">
              <h3 className="text-xl font-semibold mb-2">Your Stats</h3>
              <p className="text-gray-600">Track PRs and progress over time</p>
            </Card>
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="text-center">
        <h1 className="text-5xl font-bold text-gray-900 mb-6">
          Welcome to Cadence
        </h1>
        <p className="text-xl text-gray-600 mb-8 max-w-2xl mx-auto">
          Your personal workout logger. Track exercises, monitor progress, and
          achieve your fitness goals.
        </p>
        <div className="flex justify-center gap-4">
          <Link to="/register">
            <Button size="lg">Get Started</Button>
          </Link>
          <Link to="/login">
            <Button variant="secondary" size="lg">
              Login
            </Button>
          </Link>
        </div>
      </div>

      <div className="mt-20 grid grid-cols-1 md:grid-cols-3 gap-8">
        <Card className="text-center">
          <h3 className="text-xl font-semibold mb-3">Track Workouts</h3>
          <p className="text-gray-600">
            Log exercises, sets, reps, and weight with an intuitive interface
          </p>
        </Card>
        <Card className="text-center">
          <h3 className="text-xl font-semibold mb-3">Monitor Progress</h3>
          <p className="text-gray-600">
            View your personal records and track improvement over time
          </p>
        </Card>
        <Card className="text-center">
          <h3 className="text-xl font-semibold mb-3">Stay Consistent</h3>
          <p className="text-gray-600">
            See your weekly workout calendar and build lasting habits
          </p>
        </Card>
      </div>
    </div>
  );
}