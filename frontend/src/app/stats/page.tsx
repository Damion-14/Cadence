import { useEffect, useState } from 'react';
import { useNavigate, Link } from 'react-router';
import { useAuth } from '../../context/AuthContext';
import { api } from '../../lib/api';
import { Card } from '../../components/ui/Card';
import { Spinner } from '../../components/ui/Spinner';
import type { PRsResponse, WeeklySummary } from '../../lib/types';

export default function StatsPage() {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const [prs, setPRs] = useState<PRsResponse | null>(null);
  const [weekly, setWeekly] = useState<WeeklySummary | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }

    loadStats();
  }, [isAuthenticated, navigate]);

  const loadStats = async () => {
    try {
      const [prsData, weeklyData] = await Promise.all([
        api.get<PRsResponse>('/stats/prs'),
        api.get<WeeklySummary>('/stats/weekly'),
      ]);
      setPRs(prsData);
      setWeekly(weeklyData);
    } catch (error) {
      console.error('Failed to load stats:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[60vh]">
        <Spinner size="lg" />
      </div>
    );
  }

  const daysOfWeek = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

  const workoutsByDay = new Map<number, NonNullable<typeof weekly>['workouts'][0][]>();
  if (weekly) {
    weekly.workouts.forEach((workout) => {
      const day = workout.day_of_week;
      if (!workoutsByDay.has(day)) {
        workoutsByDay.set(day, []);
      }
      workoutsByDay.get(day)?.push(workout);
    });
  }

  return (
    <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-3xl font-bold mb-8">Your Stats</h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
        <Card>
          <h3 className="text-sm font-medium text-gray-600 mb-1">This Week</h3>
          <p className="text-3xl font-bold">{weekly?.total_workouts || 0}</p>
          <p className="text-sm text-gray-600">Workouts</p>
        </Card>
        <Card>
          <h3 className="text-sm font-medium text-gray-600 mb-1">This Week</h3>
          <p className="text-3xl font-bold">{weekly?.total_exercises || 0}</p>
          <p className="text-sm text-gray-600">Exercises</p>
        </Card>
        <Card>
          <h3 className="text-sm font-medium text-gray-600 mb-1">This Week</h3>
          <p className="text-3xl font-bold">
            {weekly?.total_volume ? weekly.total_volume.toFixed(0) : 0}
          </p>
          <p className="text-sm text-gray-600">lbs Volume</p>
        </Card>
      </div>

      <div className="mb-8">
        <h2 className="text-2xl font-bold mb-4">Weekly Calendar</h2>
        <Card>
          <div className="grid grid-cols-7 gap-2">
            {daysOfWeek.map((day, index) => {
              const dayWorkouts = workoutsByDay.get(index) || [];
              const hasWorkout = dayWorkouts.length > 0;

              return (
                <div key={day} className="text-center">
                  <div className="text-sm font-medium text-gray-600 mb-2">{day}</div>
                  <div
                    className={`h-16 rounded flex items-center justify-center ${
                      hasWorkout
                        ? 'bg-blue-100 border-2 border-blue-500'
                        : 'bg-gray-100'
                    }`}
                  >
                    {hasWorkout && (
                      <div className="text-xs">
                        <div className="font-bold text-blue-700">{dayWorkouts.length}</div>
                        <div className="text-blue-600">workout{dayWorkouts.length > 1 ? 's' : ''}</div>
                      </div>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </Card>
      </div>

      <div>
        <h2 className="text-2xl font-bold mb-4">Personal Records</h2>
        {prs && prs.prs.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {prs.prs.map((pr) => (
              <Card key={pr.exercise_name}>
                <h3 className="font-semibold text-lg mb-2">{pr.exercise_name}</h3>
                <div className="space-y-1 text-sm">
                  {pr.max_weight && (
                    <div className="flex justify-between">
                      <span className="text-gray-600">Max Weight:</span>
                      <span className="font-medium">{pr.max_weight} lbs</span>
                    </div>
                  )}
                  <div className="flex justify-between">
                    <span className="text-gray-600">Max Reps:</span>
                    <span className="font-medium">{pr.max_reps}</span>
                  </div>
                  {pr.max_volume && (
                    <div className="flex justify-between">
                      <span className="text-gray-600">Max Volume:</span>
                      <span className="font-medium">{pr.max_volume.toFixed(0)} lbs</span>
                    </div>
                  )}
                  <div className="text-xs text-gray-500 pt-2 border-t mt-2">
                    {new Date(pr.achieved_at).toLocaleDateString()}
                  </div>
                </div>
                <Link
                  to={`/stats/progress/${encodeURIComponent(pr.exercise_name)}`}
                  className="block mt-3 text-sm text-blue-600 hover:underline"
                >
                  View Progress →
                </Link>
              </Card>
            ))}
          </div>
        ) : (
          <Card>
            <p className="text-center text-gray-600 py-8">
              No personal records yet. Complete workouts to see your PRs!
            </p>
          </Card>
        )}
      </div>
    </div>
  );
}
