export interface User {
  id: number;
  email: string;
  username: string;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}

export interface WorkoutSession {
  id: number;
  user_id: number;
  name: string;
  status: 'active' | 'completed';
  started_at: string;
  completed_at?: string;
  created_at: string;
  updated_at: string;
  exercises?: Exercise[];
}

export interface Exercise {
  id: number;
  workout_session_id: number;
  name: string;
  order_index: number;
  created_at: string;
  updated_at: string;
  sets?: Set[];
}

export interface Set {
  id: number;
  exercise_id: number;
  set_number: number;
  reps: number;
  weight?: number;
  is_bodyweight: boolean;
  created_at: string;
  updated_at: string;
}

export interface SetInput {
  reps: number;
  weight?: number;
  is_bodyweight: boolean;
}

export interface CreateWorkoutRequest {
  name?: string;
}

export interface CreateExerciseRequest {
  name: string;
  sets: SetInput[];
}

export interface UpdateExerciseRequest {
  name?: string;
  sets?: SetInput[];
}

export interface PersonalRecord {
  exercise_name: string;
  max_weight?: number;
  max_reps: number;
  max_volume?: number;
  achieved_at: string;
}

export interface WorkoutSummary {
  id: number;
  name: string;
  completed_at: string;
  exercise_count: number;
  total_sets: number;
  total_volume: number;
}

export interface WorkoutInWeek {
  id: number;
  name: string;
  completed_at: string;
  day_of_week: number;
}

export interface WeeklySummary {
  week: string;
  total_workouts: number;
  total_exercises: number;
  total_volume: number;
  workouts: WorkoutInWeek[];
}

export interface ProgressDataPoint {
  date: string;
  max_weight?: number;
  max_reps: number;
  volume: number;
}

export interface ProgressResponse {
  exercise_name: string;
  data_points: ProgressDataPoint[];
}

export interface HistoryResponse {
  workouts: WorkoutSummary[];
  total: number;
}

export interface PRsResponse {
  prs: PersonalRecord[];
}

export interface ApiError {
  error: {
    code: string;
    message: string;
    details?: unknown;
    request_id: string;
  };
}
