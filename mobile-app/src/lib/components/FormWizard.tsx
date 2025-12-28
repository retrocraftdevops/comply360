/**
 * FormWizard Component
 * Multi-step form container with navigation and validation
 */

import React, { ReactNode, useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  KeyboardAvoidingView,
  Platform,
} from 'react-native';
import Button from './Button';
import StepIndicator from './StepIndicator';
import { colors, spacing, fonts } from '@/lib/utils/theme';

export interface FormStep {
  id: string;
  title: string;
  subtitle?: string;
  component: ReactNode;
  validate?: () => Promise<boolean>;
  optional?: boolean;
}

export interface FormWizardProps {
  steps: FormStep[];
  onComplete: (data: any) => Promise<void>;
  onSaveDraft?: (data: any) => Promise<void>;
  initialStep?: number;
  showSaveDraft?: boolean;
  formData?: any;
  onFormDataChange?: (data: any) => void;
}

const FormWizard: React.FC<FormWizardProps> = ({
  steps,
  onComplete,
  onSaveDraft,
  initialStep = 0,
  showSaveDraft = true,
  formData = {},
  onFormDataChange,
}) => {
  const [currentStep, setCurrentStep] = useState(initialStep);
  const [isValidating, setIsValidating] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSavingDraft, setIsSavingDraft] = useState(false);

  const totalSteps = steps.length;
  const isFirstStep = currentStep === 0;
  const isLastStep = currentStep === totalSteps - 1;
  const currentStepData = steps[currentStep];

  /**
   * Handle next step navigation
   */
  const handleNext = async () => {
    // Validate current step if validator exists
    if (currentStepData.validate) {
      setIsValidating(true);
      try {
        const isValid = await currentStepData.validate();
        if (!isValid) {
          setIsValidating(false);
          return;
        }
      } catch (error) {
        console.error('[FormWizard] Validation error:', error);
        setIsValidating(false);
        return;
      }
      setIsValidating(false);
    }

    // Move to next step or submit
    if (isLastStep) {
      await handleSubmit();
    } else {
      setCurrentStep(currentStep + 1);
    }
  };

  /**
   * Handle previous step navigation
   */
  const handlePrevious = () => {
    if (!isFirstStep) {
      setCurrentStep(currentStep - 1);
    }
  };

  /**
   * Handle form submission
   */
  const handleSubmit = async () => {
    setIsSubmitting(true);
    try {
      await onComplete(formData);
    } catch (error) {
      console.error('[FormWizard] Submit error:', error);
      // Error handling should be done by parent component
    } finally {
      setIsSubmitting(false);
    }
  };

  /**
   * Handle save draft
   */
  const handleSaveDraft = async () => {
    if (!onSaveDraft) return;

    setIsSavingDraft(true);
    try {
      await onSaveDraft(formData);
    } catch (error) {
      console.error('[FormWizard] Save draft error:', error);
    } finally {
      setIsSavingDraft(false);
    }
  };

  /**
   * Jump to specific step
   */
  const handleStepPress = (stepIndex: number) => {
    // Only allow going to previous steps
    if (stepIndex < currentStep) {
      setCurrentStep(stepIndex);
    }
  };

  return (
    <KeyboardAvoidingView
      style={styles.container}
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      keyboardVerticalOffset={Platform.OS === 'ios' ? 90 : 0}
    >
      {/* Step Indicator */}
      <View style={styles.stepIndicatorContainer}>
        <StepIndicator
          steps={steps.map((step) => ({
            id: step.id,
            title: step.title,
            optional: step.optional,
          }))}
          currentStep={currentStep}
          onStepPress={handleStepPress}
        />
      </View>

      {/* Step Content */}
      <ScrollView
        style={styles.content}
        contentContainerStyle={styles.contentContainer}
        showsVerticalScrollIndicator={false}
        keyboardShouldPersistTaps="handled"
      >
        {/* Step Title */}
        <View style={styles.header}>
          <Text style={styles.title}>{currentStepData.title}</Text>
          {currentStepData.subtitle && (
            <Text style={styles.subtitle}>{currentStepData.subtitle}</Text>
          )}
          {currentStepData.optional && (
            <Text style={styles.optional}>(Optional)</Text>
          )}
        </View>

        {/* Step Component */}
        <View style={styles.stepContent}>{currentStepData.component}</View>
      </ScrollView>

      {/* Navigation Footer */}
      <View style={styles.footer}>
        {/* Save Draft Button */}
        {showSaveDraft && onSaveDraft && (
          <View style={styles.draftButtonContainer}>
            <Button
              title="Save Draft"
              onPress={handleSaveDraft}
              variant="ghost"
              size="small"
              loading={isSavingDraft}
              disabled={isValidating || isSubmitting}
              icon="content-save-outline"
            />
          </View>
        )}

        {/* Navigation Buttons */}
        <View style={styles.navigationButtons}>
          {/* Previous Button */}
          {!isFirstStep && (
            <View style={styles.navButton}>
              <Button
                title="Previous"
                onPress={handlePrevious}
                variant="outline"
                disabled={isValidating || isSubmitting || isSavingDraft}
                icon="chevron-left"
                iconPosition="left"
              />
            </View>
          )}

          {/* Next/Submit Button */}
          <View style={[styles.navButton, isFirstStep && styles.navButtonFull]}>
            <Button
              title={isLastStep ? 'Submit' : 'Next'}
              onPress={handleNext}
              variant="primary"
              loading={isValidating || isSubmitting}
              disabled={isSavingDraft}
              icon={isLastStep ? 'check' : 'chevron-right'}
              iconPosition="right"
              fullWidth
            />
          </View>
        </View>

        {/* Step Counter */}
        <View style={styles.stepCounter}>
          <Text style={styles.stepCounterText}>
            Step {currentStep + 1} of {totalSteps}
          </Text>
        </View>
      </View>
    </KeyboardAvoidingView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  stepIndicatorContainer: {
    backgroundColor: '#FFFFFF',
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.lg,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  content: {
    flex: 1,
  },
  contentContainer: {
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.xl,
    paddingBottom: spacing['2xl'],
  },
  header: {
    marginBottom: spacing.xl,
  },
  title: {
    fontSize: fonts['2xl'],
    fontWeight: '700',
    color: colors.text,
    marginBottom: spacing.xs,
  },
  subtitle: {
    fontSize: fonts.base,
    color: colors.textSecondary,
    lineHeight: 24,
  },
  optional: {
    fontSize: fonts.sm,
    color: colors.textTertiary,
    marginTop: spacing.xs,
    fontStyle: 'italic',
  },
  stepContent: {
    flex: 1,
  },
  footer: {
    backgroundColor: '#FFFFFF',
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
  draftButtonContainer: {
    alignItems: 'center',
    marginBottom: spacing.sm,
  },
  navigationButtons: {
    flexDirection: 'row',
    gap: spacing.md,
  },
  navButton: {
    flex: 1,
  },
  navButtonFull: {
    flex: 1,
  },
  stepCounter: {
    alignItems: 'center',
    marginTop: spacing.md,
  },
  stepCounterText: {
    fontSize: fonts.sm,
    color: colors.textTertiary,
  },
});

export default FormWizard;
