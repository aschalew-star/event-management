import { defineNuxtPlugin } from '#app'
import { defineRule, configure } from 'vee-validate'
import { required, email, min, max, confirmed } from '@vee-validate/rules'
import { Form, Field, ErrorMessage } from 'vee-validate';

export default defineNuxtPlugin((nuxtApp) => {
 
  nuxtApp.vueApp.component('VeeForm', Form);
  nuxtApp.vueApp.component('VeeField', Field);
  nuxtApp.vueApp.component('VeeErrorMessage', ErrorMessage);

  // Register rules
  defineRule('required', required)
  defineRule('email', email)
  defineRule('min', min)
  defineRule('max', max)
  defineRule('confirmed', confirmed)

  // Configure validation messages
  configure({
    generateMessage: (ctx) => {
      const messages: Record<string, string> = {
        required: `The ${ctx.field} field is required`,
        email: `The ${ctx.field} must be a valid email`,
        min: `The ${ctx.field} must be at least ${ctx.rule?.params?.[0] ?? 0} characters`,
        max: `The ${ctx.field} must be at most ${ctx.rule?.params?.[0] ?? 0} characters`,
        confirmed: `The ${ctx.field} confirmation does not match`
      }

      return messages[ctx.rule?.name as string] || `The ${ctx.field} field is invalid`
    }
  })
})