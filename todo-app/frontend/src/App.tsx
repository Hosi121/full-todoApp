import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Register from './components/Register';
import CreateTask from './components/CreateTask';
import TaskList from './components/TaskList';
import CreateTag from './components/CreateTag';
import AddTagToTask from './components/AddTagToTask';
import FilterTasks from './components/FilterTasks';

const App: React.FC = () => {
  return (
    <Router>
      <Routes>
        <Route path="/register" element={<Register />} />
        <Route path="/create-task" element={<CreateTask />} />
        <Route path="/tasks" element={<TaskList />} />
        <Route path="/create-tag" element={<CreateTag />} />
        <Route path="/add-tag-to-task" element={<AddTagToTask />} />
        <Route path="/filter-tasks" element={<FilterTasks />} />
      </Routes>
    </Router>
  );
}

export default App;

