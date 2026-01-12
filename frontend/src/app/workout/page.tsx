import { useState } from 'react';
import { useNavigate } from 'react-router';
import { useAuth } from '../../context/AuthContext';
import { useWorkout } from '../../context/WorkoutContext';
import { Button } from '../../components/ui/Button';
import { Card } from '../../components/ui/Card';
import { Spinner } from '../../components/ui/Spinner';
import { ExerciseForm } from '../../components/workout/ExerciseForm';
import { ExerciseCard } from '../../components/workout/ExerciseCard';
import { Input } from '../../components/ui/Input';

export default function WorkoutPage() {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const { activeWorkout, loading, startWorkout, addExercise, deleteExercise, completeWorkout } = useWorkout();
  const [workoutName, setWorkoutName] = useState('');
  const [showExerciseForm, setShowExerciseForm] = useState(false);
  const [completing, setCompleting] = useState(false);

  if (!isAuthenticated) {
    navigate('/login');
    return null;
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[60vh]">
        <Spinner size="lg" />
      </div>
    );
  }

  const handleStartWorkout = async () => {
    await startWorkout(workoutName || undefined);
    setWorkoutName('');
  };

  const handleAddExercise = async (name: string, sets: import('../../lib/types').SetInput[]) => {
    await addExercise({ name, sets });
    setShowExerciseForm(false);
  };

  const handleDeleteExercise = async (exerciseId: number) => {
    if (confirm('Are you sure you want to delete this exercise?')) {
      await deleteExercise(exerciseId);
    }
  };

  const handleCompleteWorkout = async () => {
    if (!activeWorkout?.exercises || activeWorkout.exercises.length === 0) {
      alert('Please add at least one exercise before completing the workout');
      return;
    }

    if (confirm('Complete this workout?')) {
      setCompleting(true);
      try {
        await completeWorkout();
        navigate('/history');
      } finally {
        setCompleting(false);
      }
    }
  };

  if (!activeWorkout) {
    return (
      <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <Card>
          <h1 className="text-3xl font-bold text-center mb-6">Start a Workout</h1>
          <div className="space-y-4">
            <Input
              label="Workout Name (Optional)"
              placeholder="e.g., Push Day, Leg Day"
              value={workoutName}
              onChange={(e) => setWorkoutName(e.target.value)}
            />
            <Button onClick={handleStartWorkout} className="w-full" size="lg">
              Start Workout
            </Button>
          </div>
        </Card>

        <div className="mt-8 text-center text-gray-600">
          <p>No active workout. Start one to begin logging exercises!</p>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-6">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold">{activeWorkout.name || 'Current Workout'}</h1>
            <p className="text-gray-600 mt-1">
              Started {new Date(activeWorkout.started_at).toLocaleTimeString()}
            </p>
          </div>
          <Button
            variant="primary"
            onClick={handleCompleteWorkout}
            disabled={completing}
          >
            {completing ? 'Completing...' : 'Finish Workout'}
          </Button>
        </div>
      </div>

      <div className="space-y-6">
        {activeWorkout.exercises && activeWorkout.exercises.length > 0 ? (
          <>
            <div>
              <h2 className="text-xl font-semibold mb-4">Exercises</h2>
              <div className="space-y-4">
                {activeWorkout.exercises.map((exercise) => (
                  <ExerciseCard
                    key={exercise.id}
                    exercise={exercise}
                    onDelete={() => handleDeleteExercise(exercise.id)}
                  />
                ))}
              </div>
            </div>
          </>
        ) : (
          <Card>
            <p className="text-center text-gray-600">
              No exercises yet. Add your first exercise below!
            </p>
          </Card>
        )}

        {showExerciseForm ? (
          <ExerciseForm
            onSubmit={handleAddExercise}
            onCancel={() => setShowExerciseForm(false)}
          />
        ) : (
          <Button
            onClick={() => setShowExerciseForm(true)}
            variant="secondary"
            size="lg"
            className="w-full"
          >
            + Add Exercise
          </Button>
        )}
      </div>
    </div>
  );
}
