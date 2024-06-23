import React from 'react';
import { render, screen } from '@testing-library/react';
import TaskList from '../src/components/TaskList';

test('renders TaskList component', () => {
  render(<TaskList />);
  const linkElement = screen.getByText(/Task List/i);
  expect(linkElement).toBeInTheDocument();
});

