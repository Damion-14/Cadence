import { Card } from '../ui/Card';
import { Button } from '../ui/Button';
import type { Exercise } from '../../lib/types';

interface ExerciseCardProps {
  exercise: Exercise;
  onDelete?: () => void;
  showActions?: boolean;
}

export function ExerciseCard({ exercise, onDelete, showActions = true }: ExerciseCardProps) {
  return (
    <Card>
      <div className="flex justify-between items-start mb-3">
        <h3 className="text-lg font-semibold">{exercise.name}</h3>
        {showActions && onDelete && (
          <Button variant="danger" size="sm" onClick={onDelete}>
            Delete
          </Button>
        )}
      </div>

      <div className="space-y-2">
        {exercise.sets && exercise.sets.length > 0 ? (
          <table className="w-full text-sm">
            <thead className="text-gray-600">
              <tr>
                <th className="text-left py-1">Set</th>
                <th className="text-left py-1">Reps</th>
                <th className="text-left py-1">Weight</th>
              </tr>
            </thead>
            <tbody>
              {exercise.sets.map((set, index) => (
                <tr key={set.id} className="border-t">
                  <td className="py-2">{index + 1}</td>
                  <td className="py-2">{set.reps}</td>
                  <td className="py-2">
                    {set.is_bodyweight ? (
                      <span className="text-gray-600">Bodyweight</span>
                    ) : (
                      `${set.weight} lbs`
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <p className="text-gray-500 text-sm">No sets recorded</p>
        )}
      </div>

      {exercise.sets && exercise.sets.length > 0 && (
        <div className="mt-3 pt-3 border-t text-sm text-gray-600">
          <div className="flex justify-between">
            <span>Total Sets:</span>
            <span className="font-medium">{exercise.sets.length}</span>
          </div>
          <div className="flex justify-between">
            <span>Total Reps:</span>
            <span className="font-medium">
              {exercise.sets.reduce((sum, set) => sum + set.reps, 0)}
            </span>
          </div>
          {!exercise.sets.every((s) => s.is_bodyweight) && (
            <div className="flex justify-between">
              <span>Total Volume:</span>
              <span className="font-medium">
                {exercise.sets
                  .reduce((sum, set) => sum + (set.weight || 0) * set.reps, 0)
                  .toFixed(1)}{' '}
                lbs
              </span>
            </div>
          )}
        </div>
      )}
    </Card>
  );
}
