import { useEffect, useState } from 'react';
import { useNavigate, useParams, Link } from 'react-router';
import { useAuth } from '../../../context/AuthContext';
import { api } from '../../../lib/api';
import { Button } from '../../../components/ui/Button';
import { Spinner } from '../../../components/ui/Spinner';
import { ExerciseCard } from '../../../components/workout/ExerciseCard';
import type { WorkoutSession } from '../../../lib/types';

export default function WorkoutDetailPage() {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const { isAuthenticated } = useAuth();
  const [workout, setWorkout] = useState<WorkoutSession | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }

    if (id) {
      loadWorkout(parseInt(id));
    }
  }, [id, isAuthenticated, navigate]);

  const loadWorkout = async (workoutId: number) => {
    try {
      const data = await api.get<{ workout: WorkoutSession }>(`/workouts/${workoutId}`);
      setWorkout(data.workout);
    } catch (error) {
      console.error('Failed to load workout:', error);
      alert('Failed to load workout');
      navigate('/history');
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

  if (!workout) {
    return (
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <p>Workout not found</p>
      </div>
    );
  }

  const totalVolume = workout.exercises?.reduce(
    (sum, ex) =>
      sum +
      (ex.sets?.reduce((setSum, set) => setSum + (set.weight || 0) * set.reps, 0) || 0),
    0
  ) || 0;

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-6">
        <Link to="/history">
          <Button variant="secondary" size="sm">
            ← Back to History
          </Button>
        </Link>
      </div>

      <div className="mb-8">
        <h1 className="text-3xl font-bold">{workout.name || 'Workout'}</h1>
        <div className="mt-2 text-gray-600 space-y-1">
          <p>
            {new Date(workout.completed_at || workout.started_at).toLocaleDateString('en-US', {
              weekday: 'long',
              year: 'numeric',
              month: 'long',
              day: 'numeric',
              hour: 'numeric',
              minute: '2-digit',
            })}
          </p>
          <div className="flex gap-4 text-sm">
            <span>{workout.exercises?.length || 0} exercises</span>
            <span>
              {workout.exercises?.reduce(
                (sum, ex) => sum + (ex.sets?.length || 0),
                0
              ) || 0}{' '}
              sets
            </span>
            {totalVolume > 0 && (
              <span className="font-medium text-blue-600">
                {totalVolume.toFixed(0)} lbs total volume
              </span>
            )}
          </div>
        </div>
      </div>

      <div className="space-y-4">
        {workout.exercises && workout.exercises.length > 0 ? (
          workout.exercises.map((exercise) => (
            <ExerciseCard key={exercise.id} exercise={exercise} showActions={false} />
          ))
        ) : (
          <p className="text-gray-600">No exercises recorded</p>
        )}
      </div>
    </div>
  );
}
