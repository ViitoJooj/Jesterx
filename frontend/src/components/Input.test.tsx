import { render, screen, fireEvent } from '@testing-library/react';
import Input from './Input';

describe('Input Component', () => {
  test('renders input element', () => {
    render(<Input />);
    const inputElement = screen.getByRole('textbox');
    expect(inputElement).toBeInTheDocument();
  });

  test('displays value prop', () => {
    const testValue = 'test value';
    const handleChange = jest.fn();
    render(<Input value={testValue} onChange={handleChange} />);
    const inputElement = screen.getByRole('textbox') as HTMLInputElement;
    expect(inputElement.value).toBe(testValue);
  });

  test('calls onChange handler when value changes', () => {
    const handleChange = jest.fn();
    render(<Input onChange={handleChange} />);
    const inputElement = screen.getByRole('textbox');
    
    fireEvent.change(inputElement, { target: { value: 'new value' } });
    
    expect(handleChange).toHaveBeenCalledTimes(1);
  });

  test('applies custom className', () => {
    const customClass = 'custom-input-class';
    render(<Input className={customClass} />);
    const inputElement = screen.getByRole('textbox');
    expect(inputElement).toHaveClass(customClass);
  });

  test('passes through additional props', () => {
    render(<Input placeholder="Enter text" type="email" />);
    const inputElement = screen.getByPlaceholderText('Enter text');
    expect(inputElement).toHaveAttribute('type', 'email');
  });
});
