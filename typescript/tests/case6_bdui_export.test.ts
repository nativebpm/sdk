import test from 'node:test';
import assert from 'node:assert';
import { z, WorkflowBuilder } from '../dist/index.js';

test('Кейс 6: Экспорт под Server-Driven UI (BDUI)', async (t) => {
  // 1. Описание богатых схем пользовательского интерфейса (BDUI) на Zod
  const OnboardingFormSchema = z.object({
    full_name: z.string().min(2, 'Имя должно содержать минимум 2 символа'),
    email: z.string(),
    country: z.enum(['KZ', 'UZ', 'KG', 'AE', 'US']),
    terms_accepted: z.boolean(),
  });

  const KYCVerificationFormSchema = z.object({
    document_type: z.enum(['passport', 'id_card', 'driving_license']),
    document_number: z.string().min(5),
    monthly_income: z.number().positive(),
  });

  const FeedbackFormSchema = z.object({
    rating: z.number().int(),
    comment: z.string().optional(),
  });

  // 2. Сборка многошагового процесса взаимодействия с пользователем
  const bduiWorkflow = new WorkflowBuilder('customer_journey_bdui', 'Customer Journey Process')
    .start('start')
    .userTask('step_onboarding', 'Первичная регистрация', {
      form: OnboardingFormSchema,
      formId: 'onboarding_form',
      candidateGroups: 'registration_team',
    })
    .userTask('step_kyc', 'Верификация личности', {
      form: KYCVerificationFormSchema,
      formId: 'kyc_form',
      candidateGroups: 'compliance',
    })
    .userTask('step_feedback', 'Опрос удовлетворенности', {
      form: FeedbackFormSchema,
      formId: 'feedback_form',
    })
    .end('end', 'Процесс завершен');

  // 3. Вызов метода workflow.extractForms()
  const forms = bduiWorkflow.extractForms();

  // 4. Проверка полноты словаря форм для фронтенда (React, Vue, Capacitor)
  assert.ok(forms.onboarding_form, 'Форма onboarding_form должна присутствовать в словаре');
  assert.ok(forms.kyc_form, 'Форма kyc_form должна присутствовать в словаре');
  assert.ok(forms.feedback_form, 'Форма feedback_form должна присутствовать в словаре');

  // 5. Проверка соответствия Draft 2020-12 и корректности метаданных для Server-Driven UI генератора
  const onboardingSchema = forms.onboarding_form;
  assert.strictEqual(onboardingSchema.$schema, 'https://json-schema.org/draft/2020-12/schema');
  assert.strictEqual(onboardingSchema.type, 'object');
  assert.strictEqual(onboardingSchema.properties.full_name.type, 'string');
  assert.strictEqual(onboardingSchema.properties.full_name.minLength, 2);
  assert.deepStrictEqual(onboardingSchema.properties.country.enum, ['KZ', 'UZ', 'KG', 'AE', 'US']);
  assert.strictEqual(onboardingSchema.properties.terms_accepted.type, 'boolean');
  assert.ok(onboardingSchema.required.includes('full_name'));
  assert.ok(onboardingSchema.required.includes('email'));
  assert.ok(onboardingSchema.required.includes('country'));
  assert.ok(onboardingSchema.required.includes('terms_accepted'));

  const kycSchema = forms.kyc_form;
  assert.strictEqual(kycSchema.$schema, 'https://json-schema.org/draft/2020-12/schema');
  assert.deepStrictEqual(kycSchema.properties.document_type.enum, ['passport', 'id_card', 'driving_license']);
  assert.strictEqual(kycSchema.properties.monthly_income.type, 'number');

  // 6. Имитация рендеринга динамических виджетов фронтенда (BDUI Component Mapper)
  function simulateFrontendWidgetRender(formDefinition: any): string[] {
    const renderedWidgets: string[] = [];
    for (const [field, prop] of Object.entries<any>(formDefinition.properties)) {
      if (prop.enum) {
        renderedWidgets.push(`UISelectInput(${field}, [${prop.enum.join(', ')}])`);
      } else if (prop.type === 'boolean') {
        renderedWidgets.push(`UISwitchToggle(${field})`);
      } else if (prop.type === 'number') {
        renderedWidgets.push(`UINumericInput(${field})`);
      } else {
        renderedWidgets.push(`UITextInput(${field})`);
      }
    }
    return renderedWidgets;
  }

  const renderedOnboarding = simulateFrontendWidgetRender(onboardingSchema);
  assert.deepStrictEqual(renderedOnboarding, [
    'UITextInput(full_name)',
    'UITextInput(email)',
    'UISelectInput(country, [KZ, UZ, KG, AE, US])',
    'UISwitchToggle(terms_accepted)',
  ]);

  const renderedKyc = simulateFrontendWidgetRender(kycSchema);
  assert.deepStrictEqual(renderedKyc, [
    'UISelectInput(document_type, [passport, id_card, driving_license])',
    'UITextInput(document_number)',
    'UINumericInput(monthly_income)',
  ]);
});
