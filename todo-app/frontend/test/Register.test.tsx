import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import Register from './Register';

test('renders Register component', () => {
  render(<Register />);
  const linkElement = screen.getByText(/Register/i);
  expect(linkElement).toBeInTheDocument();
});

test('allows user to register', async () => {
  render(<Register />);
  
  fireEvent.change(screen.getByLabelText(/Username/i), { target: { value: 'testuser' } });
  fireEvent.change(screen.getByLabelText(/Password/i), { target: { value: 'password' } });
  fireEvent.click(screen.getByRole('button', { name: /Register/i }));

  const alertMock = jest.spyOn(window, 'alert').mockImplementation();
  expect(alertMock).toHaveBeenCalledWith('Registration successful');
});

