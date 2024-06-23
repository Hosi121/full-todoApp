import React, { useState, useEffect } from 'react';

interface Task {
  id: number;
  user_id: number;
  title: string;
  description: string;
  created_at: string;
  updated_at: string;
}

const TaskList: React.FC = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [userID, setUserID] = useState<number>(1); // テスト用に固定ユーザーIDを使用

  useEffect(() => {
    const fetchTasks = async () => {
      const response = await fetch(`/tasks/list?user_id=${userID}`);
      const data = await response.json();
      setTasks(data);
    };
    fetchTasks();
  }, [userID]);

  const handleDelete = async (id: number) => {
    const response = await fetch('/tasks/delete', {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ id }),
    });
    if (response.ok) {
      setTasks(tasks.filter(task => task.id !== id));
    } else {
      alert('Task deletion failed');
    }
  };

  const handleUpdate = async (id: number) => {
    const newTitle = prompt('Enter new title');
    const newDescription = prompt('Enter new description');
    if (newTitle && newDescription) {
      const response = await fetch('/tasks/update', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ id, title: newTitle, description: newDescription }),
      });
      if (response.ok) {
        setTasks(tasks.map(task => task.id === id ? { ...task, title: newTitle, description: newDescription } : task));
      } else {
        alert('Task update failed');
      }
    }
  };

  return (
    <div>
      <h2>Task List</h2>
      <ul>
        {tasks.map(task => (
          <li key={task.id}>
            <h3>{task.title}</h3>
            <p>{task.description}</p>
            <button onClick={() => handleUpdate(task.id)}>Edit</button>
            <button onClick={() => handleDelete(task.id)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default TaskList;

