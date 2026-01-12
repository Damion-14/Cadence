import { useEffect, useState } from 'react';
import { useNavigate, useParams, Link } from 'react-router';
import { useAuth } from '../../../../context/AuthContext';
import { api } from '../../../../lib/api';
import { Button } from '../../../../components/ui/Button';
import { Card } from '../../../../components/ui/Card';
import { Spinner } from '../../../../components/ui/Spinner';
import type { ProgressResponse } from '../../../../lib/types';

export default function ProgressPage() {
  const navigate = useNavigate();
  const { exerciseName } = useParams<{ exerciseName: string }>();
  const { isAuthenticated } = useAuth();
  const [progress, setProgress] = useState<ProgressResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [period, setPeriod] = useState('30d');

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }

    if (exerciseName) {
      loadProgress(exerciseName, period);
    }
  }, [exerciseName, period, isAuthenticated, navigate]);

  const loadProgress = async (name: string, periodValue: string) => {
    try {
      const data = await api.get<ProgressResponse>(
        `/stats/progress/${encodeURIComponent(name)}?period=${periodValue}`
      );
      setProgress(data);
    } catch (error) {
      console.error('Failed to load progress:', error);
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

  if (!progress || progress.data_points.length === 0) {
    return (
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Link to="/stats">
          <Button variant="secondary" size="sm">
            ← Back to Stats
          </Button>
        </Link>
        <Card className="mt-4">
          <p className="text-center text-gray-600 py-8">
            No progress data available for {exerciseName}
          </p>
        </Card>
      </div>
    );
  }

  const maxWeight = Math.max(...progress.data_points.map((p) => p.max_weight || 0));
  const maxVolume = Math.max(...progress.data_points.map((p) => p.volume));
  const maxReps = Math.max(...progress.data_points.map((p) => p.max_reps));

  return (
    <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-6">
        <Link to="/stats">
          <Button variant="secondary" size="sm">
            ← Back to Stats
          </Button>
        </Link>
      </div>

      <div className="mb-6">
        <h1 className="text-3xl font-bold">{progress.exercise_name}</h1>
        <p className="text-gray-600 mt-1">Progress over time</p>
      </div>

      <div className="mb-6 flex gap-2">
        <Button
          variant={period === '30d' ? 'primary' : 'secondary'}
          size="sm"
          onClick={() => setPeriod('30d')}
        >
          30 Days
        </Button>
        <Button
          variant={period === '90d' ? 'primary' : 'secondary'}
          size="sm"
          onClick={() => setPeriod('90d')}
        >
          90 Days
        </Button>
        <Button
          variant={period === '365d' ? 'primary' : 'secondary'}
          size="sm"
          onClick={() => setPeriod('365d')}
        >
          1 Year
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        {maxWeight > 0 && (
          <Card>
            <h3 className="text-sm font-medium text-gray-600 mb-1">Peak Weight</h3>
            <p className="text-2xl font-bold">{maxWeight} lbs</p>
          </Card>
        )}
        <Card>
          <h3 className="text-sm font-medium text-gray-600 mb-1">Peak Reps</h3>
          <p className="text-2xl font-bold">{maxReps}</p>
        </Card>
        <Card>
          <h3 className="text-sm font-medium text-gray-600 mb-1">Peak Volume</h3>
          <p className="text-2xl font-bold">{maxVolume.toFixed(0)} lbs</p>
        </Card>
      </div>

      <Card>
        <h2 className="text-xl font-semibold mb-4">Progress Chart</h2>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="border-b">
              <tr>
                <th className="text-left py-2">Date</th>
                {maxWeight > 0 && <th className="text-left py-2">Max Weight</th>}
                <th className="text-left py-2">Max Reps</th>
                <th className="text-left py-2">Volume</th>
              </tr>
            </thead>
            <tbody>
              {progress.data_points.map((point, index) => (
                <tr key={index} className="border-b">
                  <td className="py-2">
                    {new Date(point.date).toLocaleDateString('en-US', {
                      month: 'short',
                      day: 'numeric',
                    })}
                  </td>
                  {maxWeight > 0 && (
                    <td className="py-2">{point.max_weight || '-'} lbs</td>
                  )}
                  <td className="py-2">{point.max_reps}</td>
                  <td className="py-2 font-medium">{point.volume.toFixed(0)} lbs</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Card>

      <div className="mt-4 text-sm text-gray-600">
        <p>
          Showing {progress.data_points.length} workout sessions for {progress.exercise_name}
        </p>
      </div>
    </div>
  );
}
