import { render, screen, fireEvent } from '@testing-library/react';
import Button from './Button';

describe('Button Component', () => {
  test('renders button with label', () => {
    render(<Button label="Click me" />);
    const buttonElement = screen.getByText('Click me');
    expect(buttonElement).toBeInTheDocument();
  });

  test('calls onClick handler when clicked', () => {
    const handleClick = jest.fn();
    render(<Button label="Click me" onClick={handleClick} />);
    const buttonElement = screen.getByText('Click me');
    
    fireEvent.click(buttonElement);
    
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  test('does not call onClick when disabled', () => {
    const handleClick = jest.fn();
    render(<Button label="Click me" onClick={handleClick} disabled />);
    const buttonElement = screen.getByText('Click me');
    
    fireEvent.click(buttonElement);
    
    expect(handleClick).not.toHaveBeenCalled();
  });

  test('renders with correct type attribute', () => {
    render(<Button label="Submit" type="submit" />);
    const buttonElement = screen.getByText('Submit');
    expect(buttonElement).toHaveAttribute('type', 'submit');
  });

  test('defaults to button type when type not specified', () => {
    render(<Button label="Click me" />);
    const buttonElement = screen.getByText('Click me');
    expect(buttonElement).toHaveAttribute('type', 'button');
  });

  test('applies custom className', () => {
    const customClass = 'custom-button-class';
    render(<Button label="Click me" className={customClass} />);
    const buttonElement = screen.getByText('Click me');
    expect(buttonElement).toHaveClass(customClass);
  });

  test('is disabled when disabled prop is true', () => {
    render(<Button label="Click me" disabled />);
    const buttonElement = screen.getByText('Click me');
    expect(buttonElement).toBeDisabled();
  });

  test('is not disabled by default', () => {
    render(<Button label="Click me" />);
    const buttonElement = screen.getByText('Click me');
    expect(buttonElement).not.toBeDisabled();
  });
});
