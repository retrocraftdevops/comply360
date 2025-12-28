/**
 * StepIndicator Component
 * Visual progress indicator for multi-step forms
 */

import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { colors, spacing, fonts } from '@/lib/utils/theme';

export interface Step {
  id: string;
  title: string;
  optional?: boolean;
}

export interface StepIndicatorProps {
  steps: Step[];
  currentStep: number;
  onStepPress?: (stepIndex: number) => void;
}

const StepIndicator: React.FC<StepIndicatorProps> = ({
  steps,
  currentStep,
  onStepPress,
}) => {
  const totalSteps = steps.length;

  /**
   * Get step status
   */
  const getStepStatus = (stepIndex: number): 'completed' | 'current' | 'upcoming' => {
    if (stepIndex < currentStep) return 'completed';
    if (stepIndex === currentStep) return 'current';
    return 'upcoming';
  };

  /**
   * Render step circle
   */
  const renderStepCircle = (stepIndex: number, status: string) => {
    const isClickable = onStepPress && stepIndex < currentStep;

    const circleContent =
      status === 'completed' ? (
        <Icon name="check" size={16} color="#FFFFFF" />
      ) : (
        <Text
          style={[
            styles.stepNumber,
            status === 'current' && styles.stepNumberCurrent,
          ]}
        >
          {stepIndex + 1}
        </Text>
      );

    return (
      <TouchableOpacity
        style={[
          styles.stepCircle,
          status === 'completed' && styles.stepCircleCompleted,
          status === 'current' && styles.stepCircleCurrent,
          status === 'upcoming' && styles.stepCircleUpcoming,
        ]}
        onPress={isClickable ? () => onStepPress(stepIndex) : undefined}
        disabled={!isClickable}
        activeOpacity={isClickable ? 0.7 : 1}
      >
        {circleContent}
      </TouchableOpacity>
    );
  };

  /**
   * Render connector line
   */
  const renderConnector = (stepIndex: number) => {
    if (stepIndex === totalSteps - 1) return null;

    const status = getStepStatus(stepIndex);

    return (
      <View
        style={[
          styles.connector,
          status === 'completed' && styles.connectorCompleted,
        ]}
      />
    );
  };

  return (
    <View style={styles.container}>
      {steps.map((step, index) => {
        const status = getStepStatus(index);

        return (
          <View key={step.id} style={styles.stepContainer}>
            {/* Step Circle and Connector */}
            <View style={styles.stepIndicator}>
              {renderStepCircle(index, status)}
              {renderConnector(index)}
            </View>

            {/* Step Label */}
            <View style={styles.stepLabel}>
              <Text
                style={[
                  styles.stepTitle,
                  status === 'current' && styles.stepTitleCurrent,
                  status === 'upcoming' && styles.stepTitleUpcoming,
                ]}
                numberOfLines={2}
              >
                {step.title}
              </Text>
              {step.optional && status === 'upcoming' && (
                <Text style={styles.optionalLabel}>Optional</Text>
              )}
            </View>
          </View>
        );
      })}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    alignItems: 'flex-start',
  },
  stepContainer: {
    flex: 1,
    alignItems: 'center',
  },
  stepIndicator: {
    flexDirection: 'row',
    alignItems: 'center',
    width: '100%',
    marginBottom: spacing.sm,
  },
  stepCircle: {
    width: 32,
    height: 32,
    borderRadius: 16,
    alignItems: 'center',
    justifyContent: 'center',
    borderWidth: 2,
  },
  stepCircleCompleted: {
    backgroundColor: colors.primary,
    borderColor: colors.primary,
  },
  stepCircleCurrent: {
    backgroundColor: '#FFFFFF',
    borderColor: colors.primary,
  },
  stepCircleUpcoming: {
    backgroundColor: '#FFFFFF',
    borderColor: colors.border,
  },
  stepNumber: {
    fontSize: fonts.sm,
    fontWeight: '600',
    color: colors.textTertiary,
  },
  stepNumberCurrent: {
    color: colors.primary,
  },
  connector: {
    flex: 1,
    height: 2,
    backgroundColor: colors.border,
    marginLeft: spacing.xs,
  },
  connectorCompleted: {
    backgroundColor: colors.primary,
  },
  stepLabel: {
    alignItems: 'center',
    paddingHorizontal: spacing.xs,
  },
  stepTitle: {
    fontSize: fonts.xs,
    fontWeight: '600',
    color: colors.text,
    textAlign: 'center',
  },
  stepTitleCurrent: {
    color: colors.primary,
  },
  stepTitleUpcoming: {
    color: colors.textTertiary,
  },
  optionalLabel: {
    fontSize: 10,
    color: colors.textTertiary,
    marginTop: 2,
    fontStyle: 'italic',
  },
});

export default StepIndicator;
