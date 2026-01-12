import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import { api } from '../lib/api';
import type { WorkoutSession, CreateExerciseRequest, UpdateExerciseRequest } from '../lib/types';

interface WorkoutContextType {
  activeWorkout: WorkoutSession | null;
  loading: boolean;
  startWorkout: (name?: string) => Promise<void>;
  refreshWorkout: () => Promise<void>;
  addExercise: (data: CreateExerciseRequest) => Promise<void>;
  updateExercise: (exerciseId: number, data: UpdateExerciseRequest) => Promise<void>;
  deleteExercise: (exerciseId: number) => Promise<void>;
  completeWorkout: () => Promise<void>;
}

const WorkoutContext = createContext<WorkoutContextType | undefined>(undefined);

export function WorkoutProvider({ children }: { children: ReactNode }) {
  const [activeWorkout, setActiveWorkout] = useState<WorkoutSession | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadActiveWorkout();
  }, []);

  const loadActiveWorkout = async () => {
    try {
      const response = await api.get<{ workout: WorkoutSession | null }>('/workouts/active');
      setActiveWorkout(response.workout);
    } catch (error) {
      console.error('Failed to load active workout:', error);
      setActiveWorkout(null);
    } finally {
      setLoading(false);
    }
  };

  const startWorkout = async (name?: string) => {
    setLoading(true);
    try {
      const response = await api.post<{ workout: WorkoutSession }>('/workouts', { name });
      setActiveWorkout(response.workout);
    } finally {
      setLoading(false);
    }
  };

  const refreshWorkout = async () => {
    if (!activeWorkout) return;

    try {
      const response = await api.get<{ workout: WorkoutSession }>(`/workouts/${activeWorkout.id}`);
      setActiveWorkout(response.workout);
    } catch (error) {
      console.error('Failed to refresh workout:', error);
    }
  };

  const addExercise = async (data: CreateExerciseRequest) => {
    if (!activeWorkout) throw new Error('No active workout');

    await api.post<{ exercise: unknown }>(`/workouts/${activeWorkout.id}/exercises`, data);
    await refreshWorkout();
  };

  const updateExercise = async (exerciseId: number, data: UpdateExerciseRequest) => {
    if (!activeWorkout) throw new Error('No active workout');

    await api.put(`/workouts/${activeWorkout.id}/exercises/${exerciseId}`, data);
    await refreshWorkout();
  };

  const deleteExercise = async (exerciseId: number) => {
    if (!activeWorkout) throw new Error('No active workout');

    await api.delete(`/workouts/${activeWorkout.id}/exercises/${exerciseId}`);
    await refreshWorkout();
  };

  const completeWorkout = async () => {
    if (!activeWorkout) throw new Error('No active workout');

    await api.post(`/workouts/${activeWorkout.id}/complete`);
    setActiveWorkout(null);
  };

  return (
    <WorkoutContext.Provider
      value={{
        activeWorkout,
        loading,
        startWorkout,
        refreshWorkout,
        addExercise,
        updateExercise,
        deleteExercise,
        completeWorkout,
      }}
    >
      {children}
    </WorkoutContext.Provider>
  );
}

export function useWorkout() {
  const context = useContext(WorkoutContext);
  if (context === undefined) {
    throw new Error('useWorkout must be used within a WorkoutProvider');
  }
  return context;
}
