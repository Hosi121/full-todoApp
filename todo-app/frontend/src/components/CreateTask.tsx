import React, { useState, FormEvent } from 'react';

const CreateTask: React.FC = () => {
  const [title, setTitle] = useState<string>('');
  const [description, setDescription] = useState<string>('');
  const [priority, setPriority] = useState<number>(0);
  const [userID, setUserID] = useState<number>(1); // テスト用に固定ユーザーIDを使用

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    const response = await fetch('/tasks', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ user_id: userID, title, description, priority }),
    });
    if (response.ok) {
      alert('Task created successfully');
    } else {
      alert('Task creation failed');
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>
          Title:
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
          />
        </label>
      </div>
      <div>
        <label>
          Description:
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
        </label>
      </div>
      <div>
        <label>
          Priority:
          <input
            type="number"
            value={priority}
            onChange={(e) => setPriority(parseInt(e.target.value))}
          />
        </label>
      </div>
      <button type="submit">Create Task</button>
    </form>
  );
};

export default CreateTask;

