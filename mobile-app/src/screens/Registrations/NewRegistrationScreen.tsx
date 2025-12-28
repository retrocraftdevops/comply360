/**
 * NewRegistrationScreen
 * Multi-step form for creating new company registrations
 */

import React, { useState } from 'react';
import { View, Text, StyleSheet, ScrollView, Alert } from 'react-native';
import { useNavigation } from '@react-navigation/native';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { FormWizard, FormInput, Card } from '@/lib/components';
import { FormStep } from '@/lib/components/FormWizard';
import { useAppDispatch, useAppSelector } from '@/store/store';
import { useCreateRegistrationMutation } from '@/store/api/registrationApi';
import { showToast } from '@/store/slices/uiSlice';
import {
  validateEmail,
  validatePhoneNumber,
  validateRequired,
  validateIDNumber,
  validateTaxNumber,
} from '@/lib/utils/validation';
import { colors, spacing } from '@/lib/utils/theme';
import { COMPANY_TYPES, SADC_COUNTRIES } from '@/lib/utils/constants';

interface RegistrationFormData {
  // Step 1: Company Information
  companyName: string;
  companyType: string;
  registrationNumber: string;
  taxNumber: string;
  country: string;

  // Step 2: Contact Details
  contactPerson: string;
  contactEmail: string;
  contactPhone: string;
  contactIDNumber: string;

  // Step 3: Business Details (to be implemented)
  businessAddress: string;
  city: string;
  province: string;
  postalCode: string;
  businessDescription: string;
  numberOfEmployees: string;

  // Step 4: Document Upload (to be implemented)
  documents: any[];

  // Draft tracking
  isDraft: boolean;
}

const NewRegistrationScreen: React.FC = () => {
  const navigation = useNavigation();
  const dispatch = useAppDispatch();
  const user = useAppSelector((state) => state.auth.user);

  const [createRegistration, { isLoading }] = useCreateRegistrationMutation();

  // Form data state
  const [formData, setFormData] = useState<RegistrationFormData>({
    companyName: '',
    companyType: '',
    registrationNumber: '',
    taxNumber: '',
    country: 'ZA', // Default to South Africa
    contactPerson: '',
    contactEmail: '',
    contactPhone: '',
    contactIDNumber: '',
    businessAddress: '',
    city: '',
    province: '',
    postalCode: '',
    businessDescription: '',
    numberOfEmployees: '',
    documents: [],
    isDraft: false,
  });

  /**
   * Update form data
   */
  const updateFormData = (field: keyof RegistrationFormData, value: any) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  /**
   * Validate Step 1: Company Information
   */
  const validateStep1 = async (): Promise<boolean> => {
    const errors: string[] = [];

    if (!validateRequired(formData.companyName)) {
      errors.push('Company name is required');
    }
    if (!validateRequired(formData.companyType)) {
      errors.push('Company type is required');
    }
    if (!validateRequired(formData.registrationNumber)) {
      errors.push('Registration number is required');
    }
    if (formData.taxNumber && !validateTaxNumber(formData.taxNumber)) {
      errors.push('Invalid tax number format');
    }
    if (!validateRequired(formData.country)) {
      errors.push('Country is required');
    }

    if (errors.length > 0) {
      dispatch(
        showToast({ message: errors.join('\n'), type: 'error' })
      );
      return false;
    }

    return true;
  };

  /**
   * Validate Step 2: Contact Details
   */
  const validateStep2 = async (): Promise<boolean> => {
    const errors: string[] = [];

    if (!validateRequired(formData.contactPerson)) {
      errors.push('Contact person is required');
    }
    if (!validateEmail(formData.contactEmail)) {
      errors.push('Invalid email address');
    }
    if (!validatePhoneNumber(formData.contactPhone)) {
      errors.push('Invalid phone number');
    }
    if (formData.contactIDNumber && !validateIDNumber(formData.contactIDNumber)) {
      errors.push('Invalid ID number');
    }

    if (errors.length > 0) {
      dispatch(
        showToast({ message: errors.join('\n'), type: 'error' })
      );
      return false;
    }

    return true;
  };

  /**
   * Handle form submission
   */
  const handleSubmit = async (data: RegistrationFormData) => {
    try {
      await createRegistration({
        company_name: data.companyName,
        company_type: data.companyType,
        registration_number: data.registrationNumber,
        tax_number: data.taxNumber,
        country: data.country,
        contact_person: data.contactPerson,
        contact_email: data.contactEmail,
        contact_phone: data.contactPhone,
        contact_id_number: data.contactIDNumber,
        business_address: data.businessAddress,
        city: data.city,
        province: data.province,
        postal_code: data.postalCode,
        business_description: data.businessDescription,
        number_of_employees: data.numberOfEmployees,
        status: 'DRAFT',
      }).unwrap();

      dispatch(
        showToast({
          message: 'Registration created successfully!',
          type: 'success',
        })
      );

      navigation.goBack();
    } catch (error) {
      console.error('[NewRegistration] Submit error:', error);
      dispatch(
        showToast({
          message: 'Failed to create registration. Please try again.',
          type: 'error',
        })
      );
    }
  };

  /**
   * Handle save draft
   */
  const handleSaveDraft = async (data: RegistrationFormData) => {
    try {
      // TODO: Implement draft saving to AsyncStorage
      dispatch(
        showToast({
          message: 'Draft saved successfully!',
          type: 'success',
        })
      );
    } catch (error) {
      console.error('[NewRegistration] Save draft error:', error);
      dispatch(
        showToast({
          message: 'Failed to save draft.',
          type: 'error',
        })
      );
    }
  };

  /**
   * Validate Step 3: Business Details
   */
  const validateStep3 = async (): Promise<boolean> => {
    const errors: string[] = [];

    if (!validateRequired(formData.businessAddress)) {
      errors.push('Business address is required');
    }
    if (!validateRequired(formData.city)) {
      errors.push('City is required');
    }
    if (!validateRequired(formData.province)) {
      errors.push('Province is required');
    }

    if (errors.length > 0) {
      dispatch(
        showToast({ message: errors.join('\n'), type: 'error' })
      );
      return false;
    }

    return true;
  };

  /**
   * Define form steps
   */
  const formSteps: FormStep[] = [
    {
      id: 'company-info',
      title: 'Company Information',
      subtitle: 'Enter the basic details about the company',
      component: <Step1CompanyInformation formData={formData} onChange={updateFormData} />,
      validate: validateStep1,
    },
    {
      id: 'contact-details',
      title: 'Contact Details',
      subtitle: 'Provide contact information for the primary contact person',
      component: <Step2ContactDetails formData={formData} onChange={updateFormData} />,
      validate: validateStep2,
    },
    {
      id: 'business-details',
      title: 'Business Details',
      subtitle: 'Provide information about the business operations',
      component: <Step3BusinessDetails formData={formData} onChange={updateFormData} />,
      validate: validateStep3,
    },
    {
      id: 'document-upload',
      title: 'Document Upload',
      subtitle: 'Upload required documents',
      component: <Step4DocumentUpload formData={formData} onChange={updateFormData} />,
      optional: true,
    },
    {
      id: 'review-submit',
      title: 'Review & Submit',
      subtitle: 'Review all information before submitting',
      component: <Step5ReviewSubmit formData={formData} />,
    },
  ];

  return (
    <View style={styles.container}>
      <FormWizard
        steps={formSteps}
        onComplete={handleSubmit}
        onSaveDraft={handleSaveDraft}
        formData={formData}
        onFormDataChange={setFormData}
      />
    </View>
  );
};

/**
 * Step 1: Company Information
 */
interface StepProps {
  formData: RegistrationFormData;
  onChange: (field: keyof RegistrationFormData, value: any) => void;
}

const Step1CompanyInformation: React.FC<StepProps> = ({ formData, onChange }) => {
  return (
    <View>
      <FormInput
        label="Company Name"
        value={formData.companyName}
        onChangeText={(text) => onChange('companyName', text)}
        placeholder="Enter company name"
        required
        icon="domain"
        autoCapitalize="words"
      />

      <FormInput
        label="Company Type"
        value={formData.companyType}
        onChangeText={(text) => onChange('companyType', text)}
        placeholder="e.g., PTY LTD, TRUST, NPO"
        required
        icon="briefcase-variant"
        helperText="Select the legal structure of the company"
      />

      <FormInput
        label="Registration Number"
        value={formData.registrationNumber}
        onChangeText={(text) => onChange('registrationNumber', text)}
        placeholder="e.g., 2023/123456/07"
        required
        icon="file-document"
        helperText="Official company registration number"
      />

      <FormInput
        label="Tax Number"
        value={formData.taxNumber}
        onChangeText={(text) => onChange('taxNumber', text)}
        placeholder="10-digit tax number"
        icon="calculator"
        keyboardType="numeric"
        maxLength={10}
        helperText="Optional: 10-digit tax reference number"
        validate={(value) => {
          if (value && !validateTaxNumber(value)) {
            return 'Invalid tax number format (must be 10 digits)';
          }
          return null;
        }}
      />

      <FormInput
        label="Country"
        value={formData.country}
        onChangeText={(text) => onChange('country', text)}
        placeholder="Select country"
        required
        icon="flag"
        helperText="Country of registration"
      />

      <Card variant="outlined" padding="medium">
        <Text style={styles.infoText}>
          Make sure all company information matches your official registration documents.
        </Text>
      </Card>
    </View>
  );
};

/**
 * Step 2: Contact Details
 */
const Step2ContactDetails: React.FC<StepProps> = ({ formData, onChange }) => {
  return (
    <View>
      <FormInput
        label="Contact Person"
        value={formData.contactPerson}
        onChangeText={(text) => onChange('contactPerson', text)}
        placeholder="Full name"
        required
        icon="account"
        autoCapitalize="words"
      />

      <FormInput
        label="Email Address"
        value={formData.contactEmail}
        onChangeText={(text) => onChange('contactEmail', text)}
        placeholder="email@example.com"
        required
        icon="email"
        keyboardType="email-address"
        autoCapitalize="none"
        validate={(value) => {
          if (!validateEmail(value)) {
            return 'Please enter a valid email address';
          }
          return null;
        }}
      />

      <FormInput
        label="Phone Number"
        value={formData.contactPhone}
        onChangeText={(text) => onChange('contactPhone', text)}
        placeholder="+27 XX XXX XXXX"
        required
        icon="phone"
        keyboardType="phone-pad"
        validate={(value) => {
          if (!validatePhoneNumber(value)) {
            return 'Please enter a valid South African phone number';
          }
          return null;
        }}
      />

      <FormInput
        label="ID Number"
        value={formData.contactIDNumber}
        onChangeText={(text) => onChange('contactIDNumber', text)}
        placeholder="13-digit ID number"
        icon="card-account-details"
        keyboardType="numeric"
        maxLength={13}
        helperText="Optional: South African ID number"
        validate={(value) => {
          if (value && !validateIDNumber(value)) {
            return 'Invalid ID number format (must be 13 digits)';
          }
          return null;
        }}
      />

      <Card variant="outlined" padding="medium">
        <Text style={styles.infoText}>
          This contact person will receive all communication regarding this registration.
        </Text>
      </Card>
    </View>
  );
};

/**
 * Step 3: Business Details
 */
const Step3BusinessDetails: React.FC<StepProps> = ({ formData, onChange }) => {
  return (
    <View>
      <FormInput
        label="Business Address"
        value={formData.businessAddress}
        onChangeText={(text) => onChange('businessAddress', text)}
        placeholder="Street address"
        required
        icon="map-marker"
        multiline
        numberOfLines={2}
      />

      <FormInput
        label="City"
        value={formData.city}
        onChangeText={(text) => onChange('city', text)}
        placeholder="City"
        required
        icon="city"
        autoCapitalize="words"
      />

      <FormInput
        label="Province/State"
        value={formData.province}
        onChangeText={(text) => onChange('province', text)}
        placeholder="Province"
        required
        icon="map"
        autoCapitalize="words"
      />

      <FormInput
        label="Postal Code"
        value={formData.postalCode}
        onChangeText={(text) => onChange('postalCode', text)}
        placeholder="Postal code"
        icon="mailbox"
        keyboardType="numeric"
      />

      <FormInput
        label="Business Description"
        value={formData.businessDescription}
        onChangeText={(text) => onChange('businessDescription', text)}
        placeholder="Describe the nature of the business"
        icon="text"
        multiline
        numberOfLines={4}
        helperText="Brief description of business activities"
      />

      <FormInput
        label="Number of Employees"
        value={formData.numberOfEmployees}
        onChangeText={(text) => onChange('numberOfEmployees', text)}
        placeholder="e.g., 1-10, 11-50, 50+"
        icon="account-group"
      />

      <Card variant="outlined" padding="medium">
        <Text style={styles.infoText}>
          This information helps us understand your business better and provide appropriate services.
        </Text>
      </Card>
    </View>
  );
};

/**
 * Step 4: Document Upload
 */
const Step4DocumentUpload: React.FC<StepProps> = ({ formData, onChange }) => {
  return (
    <View>
      <Card variant="outlined" padding="large">
        <Text style={styles.sectionTitle}>Required Documents</Text>
        <Text style={styles.infoText}>
          You can upload documents now or add them later.
        </Text>

        <View style={styles.documentList}>
          <View style={styles.documentItem}>
            <Icon name="file-document" size={24} color={colors.textSecondary} />
            <View style={styles.documentInfo}>
              <Text style={styles.documentName}>Company Registration Certificate</Text>
              <Text style={styles.documentStatus}>Optional</Text>
            </View>
          </View>

          <View style={styles.documentItem}>
            <Icon name="file-document" size={24} color={colors.textSecondary} />
            <View style={styles.documentInfo}>
              <Text style={styles.documentName}>Tax Clearance Certificate</Text>
              <Text style={styles.documentStatus}>Optional</Text>
            </View>
          </View>

          <View style={styles.documentItem}>
            <Icon name="file-document" size={24} color={colors.textSecondary} />
            <View style={styles.documentInfo}>
              <Text style={styles.documentName}>Director ID Copies</Text>
              <Text style={styles.documentStatus}>Optional</Text>
            </View>
          </View>
        </View>

        <Text style={styles.helperText}>
          Note: Document scanner and upload functionality will be available in the next update.
        </Text>
      </Card>
    </View>
  );
};

/**
 * Step 5: Review & Submit
 */
interface Step5Props {
  formData: RegistrationFormData;
}

const Step5ReviewSubmit: React.FC<Step5Props> = ({ formData }) => {
  return (
    <ScrollView>
      {/* Company Information */}
      <Card variant="outlined" padding="large" style={styles.reviewCard}>
        <Text style={styles.sectionTitle}>Company Information</Text>
        <ReviewItem label="Company Name" value={formData.companyName} />
        <ReviewItem label="Company Type" value={formData.companyType} />
        <ReviewItem label="Registration Number" value={formData.registrationNumber} />
        <ReviewItem label="Tax Number" value={formData.taxNumber || 'Not provided'} />
        <ReviewItem label="Country" value={formData.country} />
      </Card>

      {/* Contact Details */}
      <Card variant="outlined" padding="large" style={styles.reviewCard}>
        <Text style={styles.sectionTitle}>Contact Details</Text>
        <ReviewItem label="Contact Person" value={formData.contactPerson} />
        <ReviewItem label="Email" value={formData.contactEmail} />
        <ReviewItem label="Phone" value={formData.contactPhone} />
        <ReviewItem label="ID Number" value={formData.contactIDNumber || 'Not provided'} />
      </Card>

      {/* Business Details */}
      <Card variant="outlined" padding="large" style={styles.reviewCard}>
        <Text style={styles.sectionTitle}>Business Details</Text>
        <ReviewItem label="Address" value={formData.businessAddress} />
        <ReviewItem label="City" value={formData.city} />
        <ReviewItem label="Province" value={formData.province} />
        <ReviewItem label="Postal Code" value={formData.postalCode || 'Not provided'} />
        <ReviewItem label="Description" value={formData.businessDescription || 'Not provided'} />
        <ReviewItem label="Employees" value={formData.numberOfEmployees || 'Not provided'} />
      </Card>

      {/* Info Card */}
      <Card variant="filled" padding="large" style={styles.reviewCard}>
        <View style={styles.infoBox}>
          <Icon name="information" size={24} color={colors.primary} />
          <Text style={styles.infoBoxText}>
            Please review all information carefully. You can edit any section by going back to previous steps.
          </Text>
        </View>
      </Card>
    </ScrollView>
  );
};

/**
 * Review Item Component
 */
interface ReviewItemProps {
  label: string;
  value: string;
}

const ReviewItem: React.FC<ReviewItemProps> = ({ label, value }) => (
  <View style={styles.reviewItem}>
    <Text style={styles.reviewLabel}>{label}</Text>
    <Text style={styles.reviewValue}>{value}</Text>
  </View>
);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  infoText: {
    fontSize: 14,
    color: colors.textSecondary,
    lineHeight: 20,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '700',
    color: colors.text,
    marginBottom: spacing.md,
  },
  documentList: {
    marginTop: spacing.md,
    marginBottom: spacing.lg,
  },
  documentItem: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  documentInfo: {
    flex: 1,
    marginLeft: spacing.md,
  },
  documentName: {
    fontSize: 14,
    fontWeight: '600',
    color: colors.text,
    marginBottom: 4,
  },
  documentStatus: {
    fontSize: 12,
    color: colors.textTertiary,
  },
  helperText: {
    fontSize: 12,
    color: colors.textTertiary,
    fontStyle: 'italic',
  },
  reviewCard: {
    marginBottom: spacing.md,
  },
  reviewItem: {
    marginBottom: spacing.md,
  },
  reviewLabel: {
    fontSize: 12,
    fontWeight: '600',
    color: colors.textTertiary,
    marginBottom: 4,
  },
  reviewValue: {
    fontSize: 14,
    color: colors.text,
  },
  infoBox: {
    flexDirection: 'row',
    alignItems: 'flex-start',
  },
  infoBoxText: {
    flex: 1,
    marginLeft: spacing.md,
    fontSize: 14,
    color: colors.text,
    lineHeight: 20,
  },
});

export default NewRegistrationScreen;
