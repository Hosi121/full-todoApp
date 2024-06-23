import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import CreateTask from './CreateTask';

test('renders CreateTask component', () => {
  render(<CreateTask />);
  const linkElement = screen.getByText(/Create Task/i);
  expect(linkElement).toBeInTheDocument();
});

test('allows user to create task', async () => {
  render(<CreateTask />);
  
  fireEvent.change(screen.getByLabelText(/Title/i), { target: { value: 'Test Task' } });
  fireEvent.change(screen.getByLabelText(/Description/i), { target: { value: 'Test Description' } });
  fireEvent.change(screen.getByLabelText(/Priority/i), { target: { value: '1' } });
  fireEvent.click(screen.getByRole('button', { name: /Create Task/i }));

  const alertMock = jest.spyOn(window, 'alert').mockImplementation();
  expect(alertMock).toHaveBeenCalledWith('Task created successfully');
});

