import React, { useState, FormEvent } from 'react';

const AddTagToTask: React.FC = () => {
  const [taskID, setTaskID] = useState<number>(0);
  const [tagID, setTagID] = useState<number>(0);

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    const response = await fetch('/tasks/add-tag', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ task_id: taskID, tag_id: tagID }),
    });
    if (response.ok) {
      alert('Tag added to task successfully');
    } else {
      alert('Failed to add tag to task');
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>
          Task ID:
          <input
            type="number"
            value={taskID}
            onChange={(e) => setTaskID(parseInt(e.target.value))}
          />
        </label>
      </div>
      <div>
        <label>
          Tag ID:
          <input
            type="number"
            value={tagID}
            onChange={(e) => setTagID(parseInt(e.target.value))}
          />
        </label>
      </div>
      <button type="submit">Add Tag to Task</button>
    </form>
  );
};

export default AddTagToTask;

