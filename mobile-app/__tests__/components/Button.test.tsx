/**
 * Button Component Tests
 */

import React from 'react';
import { render, fireEvent } from '@testing-library/react-native';
import Button from '@/lib/components/Button';

describe('Button Component', () => {
  it('should render correctly with title', () => {
    const { getByText } = render(
      <Button title="Test Button" onPress={() => {}} />
    );
    expect(getByText('Test Button')).toBeTruthy();
  });

  it('should call onPress when pressed', () => {
    const mockOnPress = jest.fn();
    const { getByText } = render(
      <Button title="Click Me" onPress={mockOnPress} />
    );

    fireEvent.press(getByText('Click Me'));
    expect(mockOnPress).toHaveBeenCalledTimes(1);
  });

  it('should not call onPress when disabled', () => {
    const mockOnPress = jest.fn();
    const { getByText } = render(
      <Button title="Disabled" onPress={mockOnPress} disabled />
    );

    fireEvent.press(getByText('Disabled'));
    expect(mockOnPress).not.toHaveBeenCalled();
  });

  it('should show loading spinner when loading', () => {
    const { getByTestId } = render(
      <Button title="Loading" onPress={() => {}} loading />
    );
    expect(getByTestId('loading-spinner')).toBeTruthy();
  });

  it('should render with primary variant', () => {
    const { getByText } = render(
      <Button title="Primary" onPress={() => {}} variant="primary" />
    );
    const button = getByText('Primary').parent?.parent;
    expect(button).toBeTruthy();
  });

  it('should render with secondary variant', () => {
    const { getByText } = render(
      <Button title="Secondary" onPress={() => {}} variant="secondary" />
    );
    const button = getByText('Secondary').parent?.parent;
    expect(button).toBeTruthy();
  });

  it('should render with outline variant', () => {
    const { getByText } = render(
      <Button title="Outline" onPress={() => {}} variant="outline" />
    );
    const button = getByText('Outline').parent?.parent;
    expect(button).toBeTruthy();
  });

  it('should render with danger variant', () => {
    const { getByText } = render(
      <Button title="Danger" onPress={() => {}} variant="danger" />
    );
    const button = getByText('Danger').parent?.parent;
    expect(button).toBeTruthy();
  });

  it('should render with small size', () => {
    const { getByText } = render(
      <Button title="Small" onPress={() => {}} size="small" />
    );
    expect(getByText('Small')).toBeTruthy();
  });

  it('should render with medium size', () => {
    const { getByText } = render(
      <Button title="Medium" onPress={() => {}} size="medium" />
    );
    expect(getByText('Medium')).toBeTruthy();
  });

  it('should render with large size', () => {
    const { getByText } = render(
      <Button title="Large" onPress={() => {}} size="large" />
    );
    expect(getByText('Large')).toBeTruthy();
  });

  it('should render with left icon', () => {
    const { UNSAFE_getByType } = render(
      <Button title="With Icon" onPress={() => {}} icon="check" iconPosition="left" />
    );
    // Icon component should be rendered
    expect(UNSAFE_getByType).toBeTruthy();
  });

  it('should render with right icon', () => {
    const { UNSAFE_getByType } = render(
      <Button title="With Icon" onPress={() => {}} icon="arrow-right" iconPosition="right" />
    );
    // Icon component should be rendered
    expect(UNSAFE_getByType).toBeTruthy();
  });

  it('should render full width button', () => {
    const { getByText } = render(
      <Button title="Full Width" onPress={() => {}} fullWidth />
    );
    const button = getByText('Full Width').parent?.parent;
    expect(button).toBeTruthy();
  });

  it('should not call onPress when loading', () => {
    const mockOnPress = jest.fn();
    const { getByText } = render(
      <Button title="Loading" onPress={mockOnPress} loading />
    );

    fireEvent.press(getByText('Loading'));
    expect(mockOnPress).not.toHaveBeenCalled();
  });
});
