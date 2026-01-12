import { useState, type FormEvent } from 'react';
import { Button } from '../ui/Button';
import { Input } from '../ui/Input';
import { Card } from '../ui/Card';
import type { SetInput } from '../../lib/types';

interface ExerciseFormProps {
  onSubmit: (name: string, sets: SetInput[]) => Promise<void>;
  onCancel?: () => void;
}

export function ExerciseForm({ onSubmit, onCancel }: ExerciseFormProps) {
  const [name, setName] = useState('');
  const [sets, setSets] = useState<SetInput[]>([
    { reps: 10, weight: undefined, is_bodyweight: false },
  ]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError('');

    if (!name.trim()) {
      setError('Exercise name is required');
      return;
    }

    if (sets.length === 0) {
      setError('At least one set is required');
      return;
    }

    setLoading(true);
    try {
      await onSubmit(name.trim(), sets);
      setName('');
      setSets([{ reps: 10, weight: undefined, is_bodyweight: false }]);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add exercise');
    } finally {
      setLoading(false);
    }
  };

  const addSet = () => {
    const lastSet = sets[sets.length - 1];
    setSets([...sets, { ...lastSet }]);
  };

  const removeSet = (index: number) => {
    if (sets.length > 1) {
      setSets(sets.filter((_, i) => i !== index));
    }
  };

  const updateSet = (index: number, field: keyof SetInput, value: number | boolean) => {
    const newSets = [...sets];
    if (field === 'is_bodyweight' && typeof value === 'boolean') {
      newSets[index] = { ...newSets[index], is_bodyweight: value, weight: value ? undefined : newSets[index].weight };
    } else if (field === 'reps' && typeof value === 'number') {
      newSets[index] = { ...newSets[index], reps: value };
    } else if (field === 'weight' && typeof value === 'number') {
      newSets[index] = { ...newSets[index], weight: value };
    }
    setSets(newSets);
  };

  return (
    <Card>
      <h3 className="text-lg font-semibold mb-4">Add Exercise</h3>
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Exercise Name"
          placeholder="e.g., Bench Press"
          value={name}
          onChange={(e) => setName(e.target.value)}
          disabled={loading}
          required
        />

        <div className="space-y-2">
          <label className="block text-sm font-medium text-gray-700">Sets</label>
          {sets.map((set, index) => (
            <div key={index} className="flex gap-2 items-start">
              <span className="text-sm text-gray-500 pt-2 w-8">#{index + 1}</span>
              <Input
                type="number"
                placeholder="Reps"
                value={set.reps}
                onChange={(e) => updateSet(index, 'reps', parseInt(e.target.value) || 0)}
                disabled={loading}
                required
                min="1"
                className="flex-1"
              />
              {!set.is_bodyweight && (
                <Input
                  type="number"
                  placeholder="Weight (lbs)"
                  value={set.weight || ''}
                  onChange={(e) => updateSet(index, 'weight', parseFloat(e.target.value) || 0)}
                  disabled={loading}
                  required
                  min="0"
                  step="0.5"
                  className="flex-1"
                />
              )}
              <label className="flex items-center gap-2 pt-2 whitespace-nowrap">
                <input
                  type="checkbox"
                  checked={set.is_bodyweight}
                  onChange={(e) => updateSet(index, 'is_bodyweight', e.target.checked)}
                  disabled={loading}
                  className="rounded"
                />
                <span className="text-sm">Bodyweight</span>
              </label>
              {sets.length > 1 && (
                <Button
                  type="button"
                  variant="danger"
                  size="sm"
                  onClick={() => removeSet(index)}
                  disabled={loading}
                >
                  Remove
                </Button>
              )}
            </div>
          ))}
          <Button
            type="button"
            variant="secondary"
            size="sm"
            onClick={addSet}
            disabled={loading}
          >
            + Add Set
          </Button>
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded text-sm">
            {error}
          </div>
        )}

        <div className="flex gap-2">
          <Button type="submit" disabled={loading} className="flex-1">
            {loading ? 'Adding...' : 'Add Exercise'}
          </Button>
          {onCancel && (
            <Button type="button" variant="secondary" onClick={onCancel} disabled={loading}>
              Cancel
            </Button>
          )}
        </div>
      </form>
    </Card>
  );
}
