import React, { useState, useEffect } from 'react';
import { Container, Typography, List, ListItem, ListItemText, ListItemSecondaryAction, IconButton, Box, Button } from '@mui/material';
import { Delete as DeleteIcon, Edit as EditIcon } from '@mui/icons-material';

interface Task {
  id: number;
  user_id: number;
  title: string;
  description: string;
  priority: number;
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
    const newPriority = parseInt(prompt('Enter new priority') || '0', 10);
    if (newTitle && newDescription) {
      const response = await fetch('/tasks/update', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ id, title: newTitle, description: newDescription, priority: newPriority }),
      });
      if (response.ok) {
        setTasks(tasks.map(task => task.id === id ? { ...task, title: newTitle, description: newDescription, priority: newPriority } : task));
      } else {
        alert('Task update failed');
      }
    }
  };

  return (
    <Container>
      <Box sx={{ marginTop: 8, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <Typography component="h1" variant="h5">
          Task List
        </Typography>
        <List>
          {tasks.map(task => (
            <ListItem key={task.id}>
              <ListItemText
                primary={task.title}
                secondary={`Priority: ${task.priority} - ${task.description}`}
              />
              <ListItemSecondaryAction>
                <IconButton edge="end" aria-label="edit" onClick={() => handleUpdate(task.id)}>
                  <EditIcon />
                </IconButton>
                <IconButton edge="end" aria-label="delete" onClick={() => handleDelete(task.id)}>
                  <DeleteIcon />
                </IconButton>
              </ListItemSecondaryAction>
            </ListItem>
          ))}
        </List>
      </Box>
    </Container>
  );
};

export default TaskList;
