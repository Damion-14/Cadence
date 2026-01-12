import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router';
import { useAuth } from '../../context/AuthContext';
import { api } from '../../lib/api';
import { Card } from '../../components/ui/Card';
import { Spinner } from '../../components/ui/Spinner';
import type { HistoryResponse } from '../../lib/types';

export default function HistoryPage() {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const [history, setHistory] = useState<HistoryResponse | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }

    loadHistory();
  }, [isAuthenticated, navigate]);

  const loadHistory = async () => {
    try {
      const data = await api.get<HistoryResponse>('/history?limit=20&offset=0');
      setHistory(data);
    } catch (error) {
      console.error('Failed to load history:', error);
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

  return (
    <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-3xl font-bold mb-8">Workout History</h1>

      {history && history.workouts.length > 0 ? (
        <div className="space-y-4">
          {history.workouts.map((workout) => (
            <Link key={workout.id} to={`/history/${workout.id}`}>
              <Card className="hover:shadow-lg transition-shadow cursor-pointer">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="text-xl font-semibold">{workout.name || 'Workout'}</h3>
                    <p className="text-gray-600 mt-1">
                      {new Date(workout.completed_at).toLocaleDateString('en-US', {
                        weekday: 'long',
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric',
                        hour: 'numeric',
                        minute: '2-digit',
                      })}
                    </p>
                  </div>
                  <div className="text-right text-sm text-gray-600">
                    <div>{workout.exercise_count} exercises</div>
                    <div>{workout.total_sets} sets</div>
                    {workout.total_volume > 0 && (
                      <div className="font-medium text-blue-600">
                        {workout.total_volume.toFixed(0)} lbs volume
                      </div>
                    )}
                  </div>
                </div>
              </Card>
            </Link>
          ))}
        </div>
      ) : (
        <Card>
          <div className="text-center py-12">
            <p className="text-gray-600 mb-4">No workout history yet</p>
            <Link
              to="/workout"
              className="text-blue-600 hover:underline font-medium"
            >
              Start your first workout
            </Link>
          </div>
        </Card>
      )}
    </div>
  );
}
